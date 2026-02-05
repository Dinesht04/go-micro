package email

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/smtp"
	"os"

	"github.com/dinesht04/go-micro/internal/cron"
	"github.com/dinesht04/go-micro/internal/data"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SendEmail(email *data.Email) (bool, error) {
	auth := smtp.PlainAuth("", os.Getenv("smtp_user"), os.Getenv("smtp_pass"), os.Getenv("smtp_server"))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{email.Recipient}

	msgTo := fmt.Sprintf("To: %s\r\n", email.Recipient)

	msgSubject := fmt.Sprintf("Subject: %s\r\n", email.Subject)

	msgContent := fmt.Sprintf("%s\r\n", email.Content)

	msg := []byte(msgTo +
		msgSubject +
		"\r\n" +
		msgContent)
	err := smtp.SendMail(fmt.Sprintf("%s:%s", os.Getenv("smtp_server"), os.Getenv("smtp_port")), auth, os.Getenv("smtp_user"), to, msg)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GenerateOtp(task data.Task, rdb *redis.Client, ctx context.Context) (bool, string, error) {

	var otp string

	for range task.Payload.Length {
		otp += string(GenerateRandomNumber())
	}

	_, err := rdb.HSetEXWithArgs(ctx, "otp_hashmap", &redis.HSetEXOptions{
		ExpirationType: redis.HSetEXExpirationEX,
		ExpirationVal:  120,
	}, task.Payload.UserID, otp).Result()
	if err != nil {
		return false, "Task failed during Redis insertion", err
	}

	email := &data.Email{
		Content:   fmt.Sprintf("Your OTP is: %s", otp),
		Recipient: task.Payload.UserID,
		Subject:   "OTP Requested",
	}

	status, err := SendEmail(email)
	if err != nil {
		return status, "Task failed while sending email", err
	} else {
		return status, "Task executed Successfully!", nil
	}

}

func VerifyOtp(data data.VerifyOtpParams, rdb *redis.Client, ctx *gin.Context) (bool, error) {

	res, err := rdb.HGet(ctx, "otp_hashmap", data.UserID).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		} else {
			return false, err
		}
	}

	if res == data.OTP {
		return true, nil
	} else {
		return false, nil
	}

}

func Sendmessage(task data.Task, rdb *redis.Client) (bool, string, error) {

	email := &data.Email{
		Subject:   task.Payload.Subject,
		Content:   task.Payload.Content,
		Recipient: task.Payload.UserID,
	}

	status, err := SendEmail(email)
	if err != nil {
		return status, "Sending Message email failed", err
	}

	return status, "Sent Message Successfully!", err

}

func Subscribe(task data.Task, rdb *redis.Client, ctx context.Context, c *cron.CronJobStation) (bool, string, error) {

	fields := []string{
		"subject", task.Payload.Subject,
		"content", task.Payload.Content,
	}

	err := rdb.HSet(ctx, "subscriptionContentMap"+task.Payload.ContentType, fields).Err()
	if err != nil {
		return false, "subscription content insertiong error", err
	}

	err = c.Subscribe(task.Payload.UserID, task.Payload.Frequency, task.Payload.ContentType)

	return true, "subscribed successfully", nil

	// content => This can be changed through /updateSubscriptionContent
	//Content should be accessed dynamically in cron job since it is subject to change.
}

func Unsubscribe(task data.Task, rdb *redis.Client, c *cron.CronJobStation) (bool, string, error) {
	err := c.Unsubscribe(task.Payload.UserID, task.Payload.ContentType)
	if err != nil {
		return false, "error unsubscribing", err
	}
	return true, "unsubscribed successfully", nil
}

func GenerateRandomNumber() string {
	maxInt := big.NewInt(9)
	randNum, err := rand.Int(rand.Reader, maxInt)
	if err != nil {
		log.Fatal(err)
	}

	return randNum.String()
}
