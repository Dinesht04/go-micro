package server

import (
	"log/slog"
	"net/http"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/dinesht04/go-micro/internal/services"
	"github.com/gin-gonic/gin"
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
		s.logger.Info("Endpoint health check")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/task", func(ctx *gin.Context) {
		var task data.Task
		err := ctx.ShouldBind(&task)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "false",
				"error":  "INVALID FORMAT",
				"msg":    err,
			})
		}

		validated, msg, err := services.ValidateTask(task, s.rdb, ctx)
		if err != nil {
			s.logger.Info("Error Validating Task", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "false",
				"error":  "Internal Server Error",
			})
			ctx.Abort()
			return
		}

		if !validated {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": validated,
				"error":  msg,
			})
			return
		}

		status, msg, err := services.ProduceTask(task, s.rdb, ctx)
		if err != nil {
			s.logger.Info("Error Producing Task", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "false",
				"error":  "Internal Server Error",
			})
			ctx.Abort()
			return
		}

		if !status {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": status,
				"error":  msg,
			})
			ctx.Abort()
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"status": status,
				"msg":    msg,
			})
			return
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

		verified, err := services.VerifyOtp(req, s.rdb, ctx)
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

	r.POST("/subscriptionContent", func(ctx *gin.Context) {
		var subreq data.CreateContent

		err := ctx.ShouldBind(&subreq)
		if err != nil {
			s.logger.Info("Invalid Request Format", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Request Format",
			})
			ctx.Abort()
			return
		}

		msg, err := services.CreateContentType(subreq, s.rdb, ctx)
		if err != nil {
			s.logger.Info("Error while updating subscription content map", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
				"msg":   msg,
			})
			ctx.Abort()
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"status": true,
				"msg":    msg,
			})
			return
		}

	})

	r.PUT("/subscriptionContent", func(ctx *gin.Context) {
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

		status, msg, err := services.UpdateContentType(subReq, s.rdb, ctx)
		if err != nil {
			s.logger.Info(msg, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": status,
				"error":  "Internal Server Error",
				msg:      msg,
			})
			ctx.Abort()
			return

		}

		if !status {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": status,
				"msg":    msg,
			})
			ctx.Abort()
			return
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"status": status,
				"msg":    msg,
			})
			return

		}

	})

	r.GET("/metrics", func(ctx *gin.Context) {
		//us redis to store and access total jobs, successful jobs, etv

		executed, failed, successful, msg, err := services.GetMetrics(s.rdb, ctx)
		if err != nil {
			s.logger.Info(msg, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"msg":    msg,
			})
			ctx.Abort()
			return

		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":              true,
			"Total Jobs Executed": executed,
			"Jobs Successful":     successful,
			"Jobs Failed":         failed,
		})

	})

	//8080
	r.Run(s.Port)
}
