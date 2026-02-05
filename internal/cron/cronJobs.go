package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

type userRecord struct {
	ClientEmail string
	cronId      cron.EntryID
}

type CronJobStation struct {
	context context.Context
	rdb     *redis.Client
	cron    *cron.Cron
	Jobs    map[string]userRecord
	logger  *slog.Logger
}

func CreateNewCronJobStation(ctx context.Context, rdb *redis.Client, logger *slog.Logger) *CronJobStation {
	c := cron.New()
	c.Start()
	return &CronJobStation{
		cron:    c,
		Jobs:    make(map[string]userRecord),
		context: ctx,
		rdb:     rdb,
		logger:  logger,
	}
}

func (c *CronJobStation) Subscribe(userEmailId string, frequency string, contentType string) error {

	cronId, err := RegisterCronSendingEmailJob(c, userEmailId, frequency, contentType)
	if err != nil {
		return err
	}

	record := userRecord{
		cronId: cronId,
	}

	c.Jobs[userEmailId+contentType] = record
	c.logger.Info("Cron Job Added Successfully!",
		"cronid", cronId,
		"userEmailID", userEmailId,
		"frequency", frequency,
		"contentType", contentType)
	return nil
}

func (c *CronJobStation) Unsubscribe(userEmailId string, contentType string) error {
	Record, ok := c.Jobs[userEmailId+contentType]
	if !ok {
		return fmt.Errorf("Record doesnt exist how to unsubscruibe?")
	}
	c.cron.Remove(Record.cronId)
	delete(c.Jobs, userEmailId+contentType)
	c.logger.Info("Cron Job Removed Successfully!",
		"cronid", Record.cronId,
		"userEmailID", userEmailId,
		"contentType", contentType)

	return nil
}

func RegisterCronSendingEmailJob(c *CronJobStation, userEmailId string, frequency string, contentType string) (cron.EntryID, error) {

	cronId, err := c.cron.AddFunc(frequency, func() {
		content, err := c.rdb.HGetAll(c.context, "subscriptionContentMap"+contentType).Result()
		if err != nil {
			if err == redis.Nil {
				fmt.Println("This type of content doesnt exist")
				c.logger.Info("The type of content doesn't exist in db")
				return
			} else {
				c.logger.Info("Error accessing content type from db")
				fmt.Println(err)
				return
			}
		}

		taskName := "Automated Email to: " + userEmailId

		messageTask := data.Task{
			Id:   uuid.NewString(),
			Task: taskName,
			Type: "message",
			Payload: data.Payload{
				UserID:  userEmailId,
				Subject: content["subject"],
				Content: content["content"],
			},
			Retries: 3,
		}

		encodedTask, err := json.Marshal(&messageTask)
		if err != nil {
			c.logger.Info("Error decoding task")
			fmt.Println(err)
			return
		}

		err = c.rdb.RPush(c.context, "taskQueue", encodedTask).Err()
		if err != nil {
			fmt.Println("Error Pushing to task Queue")
			c.logger.Info("Error pushing task to queue", "taskId", messageTask.Id)
			fmt.Println(err)
			return
		} else {
			c.logger.Info("Cron job succefful, task added to Queue", "taskId", messageTask.Id)
		}

	})

	return cronId, err
}
