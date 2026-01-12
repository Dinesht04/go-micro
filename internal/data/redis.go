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

	encodedTask, err := json.Marshal(&task)
	if err != nil {
		log.Fatal(err)
	}

	err = rdb.Set(ctx, "sampleTask", encodedTask, 0).Err()
	if err != nil {
		log.Fatal(err)
	}

	savedTask := rdb.Get(ctx, "sampleTask")
	if err == redis.Nil {
		fmt.Println("sampleTask dont exist")
	} else if err != nil {
		log.Fatal(err)
	} else {
		var decodedTask Task

		byteTask, err := savedTask.Bytes()
		if err != nil {
			log.Panic(err)
		}
		err = json.Unmarshal(byteTask, &decodedTask)
		if err != nil {
			fmt.Println(string(byteTask))
			log.Panic("unamrshal", err)
		}
		fmt.Println("sample task:", decodedTask)
	}

	return rdb

}
