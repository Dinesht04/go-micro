package main

import (
	"context"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/dinesht04/go-micro/internal/server"
	"github.com/dinesht04/go-micro/internal/worker"
)

func main() {

	//connect to redis
	//start the server

	ctx := context.Background()

	rdb := data.NewRedisClient(ctx)
	server := server.NewServer(rdb)

	Workstation := worker.NewWorkStation(rdb, 3)
	Workstation.StartWorkers(ctx)

	server.StartServer()

}
