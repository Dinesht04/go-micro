okay so this is an email microservice written in go. it supports otp generation, verification, sending a message and subscribing and unsubscribing to emails.
uses :-
github.com/gin-gonic/gin v1.11.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.17.2
	github.com/robfig/cron/v3 v3.0.1

Uses smtp to send emails.

env vars required:
smtp_user=
smtp_pass=
smtp_server=
smtp_port=

includes 5 endpoints
1)/ping GET
    health checkpoint
2)/task POST
    send a task to add to the queue
    payload:-

    type Task struct {
	Id      string  `json:"id"`
	Task    string  `json:"task" binding:"required"`
	Type    string  `json:"type" binding:"required,oneof= generateOtp message subscribe unsubscribe"`
	Payload Payload `json:"payload" binding:"required"`
	Retries int     `json:"retries" binding:"required"`
    }

    type Payload struct {
	UserID      string `json:"userId" binding:"required,email"`
	ContentType string `json:"content_type" binding:"required_if=Type subscribe,required_if=Type unsubscribe"`
	Length      int    `json:"length" binding:"required_if=Type generateOtp,lte=8"`
	Frequency   string `json:"frequency" binding:"omitempty,required_if=Type subscribe,oneof= @monthly @weekly @daily @hourly"`
	Content     string `json:"content" binding:"required_if=Type message,required_if=Type subscribe"`
	Subject     string `json:"subject" binding:"required_if=Type message,required_if=Type subscribe"`
    }

    4 Types of Taks:-
    1) Message - sends a message
    2)GenerateOTP - generates otp and stores it in a redis hashmap along with the user's email
    3)Subscribe - frequency can be hourly, daily, weekly, monthly and yearly
    4)Unsubscribe -  unsubscribes

3)/verify POST
    verify the otp
    Payload:-
    type VerifyOtpParams struct {
        UserID string `json:"userId" binding:"required,email"`
        OTP    string `json:"otp" binding:"required,min=4,max=8"`
    }

4)/updateSubscriptionContent POST
    update the subscription content
    payload:-
    type UpdateContent struct {
	ContentType string `json:"content_type" binding:"required"`
	Content     string `json:"content" binding:"required"`
	Subject     string `json:"subject" binding:"required"`
    }

5)/metrics GET
    gives the metrics