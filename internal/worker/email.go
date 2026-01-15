package worker

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/smtp"
	"os"

	"github.com/dinesht04/go-micro/internal/data"
	"github.com/redis/go-redis/v9"
)

func sendEmail(email *data.Email) bool {
	fmt.Println("sending mail?")

	auth := smtp.PlainAuth("", os.Getenv("smtp_user"), os.Getenv("smtp_pass"), os.Getenv("smtp_server"))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{email.Recipient}

	msgTo := fmt.Sprintf("To: %s\r\n", email.Recipient)

	msgSubject := fmt.Sprintf("%s\r\n", email.Subject)

	msgContent := fmt.Sprintf("%s\r\n", email.Content)

	msg := []byte(msgTo +
		msgSubject +
		"\r\n" +
		msgContent)
	err := smtp.SendMail(fmt.Sprintf("%s:%s", os.Getenv("smtp_server"), os.Getenv("smtp_port")), auth, os.Getenv("smtp_user"), to, msg)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func GenerateOtp(task data.Task, rdb *redis.Client, ctx context.Context) {

	var otp string

	for range task.Payload.Length {
		otp += string(GenerateRandomNumber())
	}

	fmt.Println("OTP generated bruh: ", otp)

	err := rdb.HSetEXWithArgs(ctx, "otp_hashmap", &redis.HSetEXOptions{
		Condition:      redis.HSetEXFNX,
		ExpirationType: redis.HSetEXExpirationEX,
		ExpirationVal:  120,
	}, task.Payload.UserID, otp).Err()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("OTP stored in HMAP: successfully")

	email := &data.Email{
		Content:   fmt.Sprintf("Your OTP is: %s", otp),
		Recipient: task.Payload.UserID,
		Subject:   "OTP Requested",
	}

	sendEmail(email)

}

func VerifyOtp() {}

func Sendmessage() {}

func Subscribe() {}

func Unsubscribe() {}

func GenerateRandomNumber() string {
	maxInt := big.NewInt(9)
	randNum, err := rand.Int(rand.Reader, maxInt)
	if err != nil {
		log.Fatal(err)
	}

	return randNum.String()
}
