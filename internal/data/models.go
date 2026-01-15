package data

type Task struct {
	Id      string  `json:"id"`
	Task    string  `json:"task" binding:"required"`
	Type    string  `json:"type" binding:"required,oneof= generateOtp verifyOtp message subscribe unsubscribe"`
	Payload Payload `json:"payload" binding:"required"`
	Retries int     `json:"retries" binding:"required"`
}

type Payload struct {
	UserID    string `json:"userId" binding:"required,email"`
	Length    int    `json:"length" binding:"required_if=Type generateOtp,lte=8"`
	Frequency int    `json:"frequency" binding:"required,lte=8"`
}

type Email struct {
	Content   string
	Subject   string
	Recipient string
}
