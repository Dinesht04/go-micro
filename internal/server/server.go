package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	Port string
	rdb  *redis.Client
}

func NewServer(rdb *redis.Client) *Server {
	server := &Server{
		Port: ":8080",
		rdb:  rdb,
	}
	return server
}

func (s *Server) StartServer() {
	//start server and pass params into redis
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/task", func(ctx *gin.Context) {
		var task data.Task
		err := ctx.ShouldBind(&task)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": "INVALID FORMAT",
			})
		}
		fmt.Println(task)

		encodedTask, err := json.Marshal(&task)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": "Error while marhsalling task",
			})
			return
		}

		err = s.rdb.RPush(ctx, "taskQueue", encodedTask).Err()
		if err != nil {
			log.Fatal(err)
		}

		//mantain a map in memory

		//how to implement the retries mechanic?
		//how will the queue insertion work? - draw on excalidraw

	})

	r.GET("/task", func(ctx *gin.Context) {

		//log tasks here?

		val := s.rdb.RPop(ctx, "taskQueue")

		if err := val.Err(); err == redis.Nil {
			ctx.JSON(http.StatusOK, gin.H{
				"Status": "Queue is empty",
			})
			return
		}

		var Task data.Task

		encodedTask, err := val.Bytes()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": "Error while decoding redis string",
			})
			fmt.Println(val.Err())
			return
		}

		err = json.Unmarshal(encodedTask, &Task)

		fmt.Println("redis string: ", val.String())

		ctx.JSON(http.StatusOK, gin.H{
			"task": Task,
		})

	})

	//8080
	r.Run(s.Port)
}
