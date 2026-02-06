package data

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

type Email struct {
	Content   string
	Subject   string
	Recipient string
}

type VerifyOtpParams struct {
	UserEmail string `json:"userEmail" binding:"required,email"`
	Otp       string `json:"otp" binding:"required,min=4,max=8"`
}

type UpdateContent struct {
	ContentType string `json:"content_type" binding:"required"`
	Content     string `json:"content" binding:"required"`
	Subject     string `json:"subject" binding:"required"`
}
