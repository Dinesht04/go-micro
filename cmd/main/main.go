package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/dinesht04/go-micro/internal/cron"
	"github.com/dinesht04/go-micro/internal/data"
	"github.com/dinesht04/go-micro/internal/server"
	"github.com/dinesht04/go-micro/internal/worker"
	"github.com/joho/godotenv"
)

func main() {

	//connect to redis
	//start the server

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	rdb, err := data.NewRedisClient(ctx, logger)
	if err != nil {
		logger.Info("Error Initiating redis client", "error", err)
	}

	server := server.NewServer(rdb, logger)
	CronJobStation := cron.CreateNewCronJobStation(ctx, rdb, logger)

	Workstation := worker.NewWorkStation(rdb, 3, CronJobStation)
	Workstation.StartWorkers(ctx, logger)

	server.StartServer()

}
