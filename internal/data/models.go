package data

type Task struct {
	Id      int    `json:"id" binding:"required"`
	Task    string `json:"task" binding:"required"`
	Retries int    `json:"retries" binding:"required"`
}
