package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/dinesht04/go-micro/internal/email"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	Port   string
	rdb    *redis.Client
	logger *slog.Logger
}

func NewServer(rdb *redis.Client, logger *slog.Logger) *Server {
	server := &Server{
		Port:   ":8080",
		rdb:    rdb,
		logger: logger,
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
			log.Fatal(err)
		}
		// fmt.Println(task)
		task.Id = uuid.NewString()

		encodedTask, err := json.Marshal(&task)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ERROR": "Error while marhsalling task",
			})
			return
		}

		err = s.rdb.RPush(ctx, "taskQueue", encodedTask).Err()
		if err != nil {
			s.logger.Info("Error pushing task to the Queue", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			ctx.Abort()
			return
		} else {
			s.logger.Info("Pushed Task to the Queue",
				"taskId", task.Id,
				"taskName", task.Task)

			ctx.JSON(http.StatusOK, gin.H{
				"status":   "success",
				"message":  "Task Added to the Queue Successfully",
				"taskId":   task.Id,
				"taskName": task.Task,
				"taskType": task.Type,
			})
		}

	})

	r.POST("/verify", func(ctx *gin.Context) {
		var req data.VerifyOtpParams
		err := ctx.ShouldBind(&req)
		if err != nil {
			s.logger.Info("Invalid Request Format", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Request Format",
			})
			ctx.Abort()
			return
		}

		verified, err := email.VerifyOtp(req, s.rdb, ctx)
		if err != nil {
			s.logger.Info("Error while verifying OTP", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"type":     "otp verification",
			"verified": verified,
		})

	})

	r.POST("/updateSubscriptionContent", func(ctx *gin.Context) {

		var subReq data.UpdateContent

		err := ctx.ShouldBind(&subReq)
		if err != nil {
			s.logger.Info("Invalid Request Format", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Request Format",
			})
			ctx.Abort()
			return
		}

		fields := []string{
			"subject", subReq.Subject,
			"content", subReq.Content,
		}

		err = s.rdb.HSet(ctx, "subscriptionContentMap"+subReq.ContentType, fields).Err()
		if err != nil {
			s.logger.Info("Error while updating subscription content map", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			ctx.Abort()
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"success": "Content map updated succesfully",
			})
			return
		}

	})

	r.GET("/metrics", func(ctx *gin.Context) {
		//us redis to store and access total jobs, successful jobs, etv
		totalTasksExecuted, err := s.rdb.Get(ctx, "totalTasksExecuted").Result()
		if err != nil {
			if err == redis.Nil {
				totalTasksExecuted = "0"
			} else {
				s.logger.Info("Error while accessing total tasks executed", "error", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
				ctx.Abort()
				return
			}
		}
		totalTasksFailed, err := s.rdb.Get(ctx, "totalTasksFailed").Result()
		if err != nil {
			if err == redis.Nil {
				totalTasksFailed = "0"
			} else {
				s.logger.Info("Error while accessing total tasks failed", "error", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
				ctx.Abort()
				return
			}

		}
		totalTasksSuccessful, err := s.rdb.Get(ctx, "totalTasksSuccessful").Result()
		if err != nil {
			if err == redis.Nil {
				totalTasksFailed = "0"
			} else {
				s.logger.Info("Error while accessing total successful tasks", "error", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
				ctx.Abort()
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Total Jobs Executed": totalTasksExecuted,
			"Jobs Successful":     totalTasksSuccessful,
			"Jobs Failed":         totalTasksFailed,
		})

	})

	//8080
	r.Run(s.Port)
}
