package worker

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type WorkStation struct {
	Rdb     *redis.Client
	Workers int
}

func NewWorkStation(rdb *redis.Client, num int) *WorkStation {
	return &WorkStation{
		Rdb:     rdb,
		Workers: num,
	}
}

func (w *WorkStation) StartWorkers(ctx context.Context) {

	for range w.Workers {
		go Worker(w.Rdb, ctx)
	}
}

func Worker(rdb *redis.Client, ctx context.Context) {

	for {

		results, err := rdb.BLPop(ctx, time.Minute, "taskQueue").Result()
		if err == redis.Nil {
			fmt.Println("NOthign found within timeout, waiting for 1 min again")
			continue
		}

		result := results[1]
		fmt.Println("task popped is:", result)

		var task data.Task

		err = json.Unmarshal([]byte(result), &task)
		if err != nil {
			log.Fatal("inside worker, first task", err)
		}

		task.Retries = task.Retries - 1

		task.Id = uuid.NewString()

		taskType := task.Type
		// status := sendEmail()

		switch taskType {
		case "generateOtp":
			GenerateOtp(task, rdb, ctx)
		case "verifyOtp":
			VerifyOtp()
		case "message":
			Sendmessage()
		case "subscribe":
			Subscribe()
		case "unsubscribe":
			Unsubscribe()
		default:
			fmt.Println("Random shi bruh")
		}

		// if !status {
		// 	fmt.Println("Performing Task: ", task.Id, " Failed!, Adding back to queue")
		// 	fmt.Println("Retries left: ", task.Retries)

		// 	if task.Retries <= 0 {
		// 		fmt.Println("Task: ", task.Task, " Retries ended, returning...")
		// 		continue
		// 	}

		// 	encodedTask, err := json.Marshal(&task)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	fmt.Println("Inserting again....")
		// 	status, err := rdb.RPush(ctx, "taskQueue", encodedTask).Result()
		// 	fmt.Println("Inserting again Result: ", status, " | err: ", err)

		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// } else {
		// 	fmt.Println("Performed Task: ", task.Task, " Successfully!")
		// }

	}

}

func executeTask() bool {
	max := big.NewInt(30)
	failure := big.NewInt(10)
	rand, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal(err)
	}
	if rand.Int64() > failure.Int64() {
		return true
	} else {
		return false
	}
}
