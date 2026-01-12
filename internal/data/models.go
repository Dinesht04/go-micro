package data

type Task struct {
	Task    string `json:"task" binding:"required"`
	Retries int    `json:"retries" binding:"required"`
}
