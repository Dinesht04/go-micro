package data

import "encoding/json"

type Task struct {
	Task    string
	Retries int
}

func (t Task) MarshalBinary() (data []byte, err error) {
	return json.Marshal(t)
}

func (t *Task) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}
