package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(401, gin.H{
				"error": "unauthorized",
				"msg":   "auth header missing",
			})
			ctx.Abort()
			return
		}

		verified, err := VerifyJWT(authHeader)
		if err != nil {
			ctx.JSON(401, gin.H{
				"authorized": verified,
				"error":      err,
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

type SignUpRequest struct {
	ID string `json:"id" binding:"required,email"`
}

func HandleSignup(rdb *redis.Client) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		//check redis hashmap if an api key has been assigned to the specific userID or not. if not then continue..
		var req SignUpRequest
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"error":  err,
				"msg":    "error while binding",
			})
			ctx.Abort()
			return
		}

		res, err := rdb.HExists(ctx, "UserList", req.ID).Result()
		if res {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": false,
				"error":  "Api key already registered in this username. If you wish to get a new api key, Deregister.",
			})
			ctx.Abort()
			return
		}

		token, err := CreateJWT(req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"error":  err,
				"msg":    "error while craeting jwt",
			})
			ctx.Abort()
			return
		}

		err = rdb.HSet(ctx, "UserList", req.ID, token).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"error":  err,
				"msg":    "error while setting userlist",
			})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"apiKey": token,
		})
	}
}

type DeRegister struct {
	ID string `json:"id" binding:"required,email"`
}

func HandleDeRegister(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//in case user has lost, remove the userID from the hashmap so that another one can be assigned to them,
		var req DeRegister
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"error":  err,
				"msg":    "error while binding",
			})
			ctx.Abort()
			return
		}

		res, err := rdb.HExists(ctx, "UserList", req.ID).Result()
		fmt.Println(res, err)
		if !res {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": false,
				"error":  "Api key for user doesn't exist.",
			})
			ctx.Abort()
			return
		}

		err = rdb.HDel(ctx, "UserList", req.ID).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": false,
				"error":  err,
				"msg":    "error deleting user from list",
			})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": true,
			"msg":    "Deregistered successfully! You may now register for a new api key",
		})

	}
}
