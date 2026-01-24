package data

type Task struct {
	Id      string  `json:"id"`
	Task    string  `json:"task" binding:"required"`
	Type    string  `json:"type" binding:"required,oneof= generateOtp verifyOtp message subscribe unsubscribe"`
	Payload Payload `json:"payload" binding:"required"`
	Retries int     `json:"retries" binding:"required"`
}

type Payload struct {
	UserID      string `json:"userId" binding:"required,email"`
	ContentType string `json:"content_type" binding:"required_if=Type subscribe"`
	Length      int    `json:"length" binding:"required_if=Type generateOtp,lte=8"`
	Frequency   string `json:"frequency" binding:"required_if=Type subscribe"`
	Content     string `json:"content" binding:"required_if=Type message"`
	Subject     string `json:"subject" binding:"requred_if=Type message"`
}

type Email struct {
	Content   string
	Subject   string
	Recipient string
}

type VerifyOtpParams struct {
	UserID string `json:"userId" binding:"required,email"`
	OTP    string `json:"otp" binding:"required,min=4,max=8"`
}
