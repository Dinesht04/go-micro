package main

import (
	"context"

	"github.com/dinesht04/go-micro/internal/data"
)

func main() {

	//connect to redis
	//start the server

	ctx := context.Background()

	_ = data.NewRedisClient(ctx)

}
