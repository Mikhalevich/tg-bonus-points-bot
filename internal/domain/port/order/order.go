package order

import (
	"strconv"
)

type ID int

func (id ID) Int() int {
	return int(id)
}

func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func IDFromInt(id int) ID {
	return ID(id)
}

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	StatusCreated    Status = "created"
	StatusInProgress Status = "in_progress"
	StatusReady      Status = "ready"
	StatusCompleted  Status = "completed"
	StatusCanceled   Status = "canceled"
)
