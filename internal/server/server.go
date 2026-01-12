package server

import (
	"fmt"
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
		Port: "8080",
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

		//how to implement the retries mechanic?
		//how will the queue insertion work? - draw on excalidraw
	})

	//8080
	r.Run(s.Port)
}
