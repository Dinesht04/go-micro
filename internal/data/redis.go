package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// type RedisClient struct {
// 	Client *redis.NewClient
// }

type Task struct {
	Task    string
	Retries int
}

func (t Task) MarshalBinary() (data []byte, err error) {
	return json.Marshal(t)
}

func (t *Task) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

func NewRedisClient(ctx context.Context) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	task := Task{
		Task:    "sample task",
		Retries: 3,
	}

	err := rdb.Set(ctx, "sampleTask", task, 0).Err()
	if err != nil {
		log.Fatal(err)
	}

	savedTask := &Task{}

	err = rdb.Get(ctx, "sampleTask").Scan(savedTask)
	if err == redis.Nil {
		fmt.Println("sampleTask dont exist")
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("sample Task: ", savedTask.Retries)
	}

	return rdb

}
