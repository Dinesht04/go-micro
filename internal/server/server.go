package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Port string
}

func NewServer() *Server {
	server := &Server{
		Port: "8080",
	}
	return server
}

func StartServer() {
	//start server and pass params into redis
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//8080
	r.Run()
}
