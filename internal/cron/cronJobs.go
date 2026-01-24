package cron

import (
	"context"
	"fmt"
	"log"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/dinesht04/go-micro/internal/worker"
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
}

func CreateNewCronJobStation() *CronJobStation {
	c := cron.New()
	return &CronJobStation{
		cron: c,
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

	c.Jobs[userEmailId] = record
	fmt.Println("cron job added successfully")
	return nil
}

func (c *CronJobStation) Unsubscribe(clientID string, userEmailId string) {
	Record, ok := c.Jobs[userEmailId]
	if !ok {
		log.Fatal("Record doesnt exist how to unsubscruibe?")
	}
	c.cron.Remove(Record.cronId)
	delete(c.Jobs, userEmailId)
	fmt.Println("cron job removed successfully")
}

func RegisterCronSendingEmailJob(c *CronJobStation, userEmailId string, frequency string, contentType string) (cron.EntryID, error) {
	fmt.Println("Registering for the job")

	cronId, err := c.cron.AddFunc("@hourly", func() {
		//send mail, get content here from rdb and pass it onto sendMail()

		//this stuff goes to logs
		fmt.Println("Sending a mail")
		content, err := c.rdb.HGet(c.context, "subscriptionContentMap", contentType).Result()
		if err != nil {
			if err == redis.Nil {
				fmt.Println("This type of content doesnt exist")
			} else {
				fmt.Println("Err accessing content type")
				fmt.Println(err)
			}
		}

		email := &data.Email{
			Recipient: userEmailId,
			Subject:   "Automated Mail",
			Content:   content,
		}

		success, err := worker.SendEmail(email)
		if !success {
			fmt.Println("Err sending email")
			fmt.Println(err)
		}

	})

	return cronId, err
}
