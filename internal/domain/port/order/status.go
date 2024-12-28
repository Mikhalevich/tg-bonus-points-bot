package order

import (
	"fmt"
)

type Status string

const (
	StatusCreated    Status = "created"
	StatusInProgress Status = "in_progress"
	StatusReady      Status = "ready"
	StatusCompleted  Status = "completed"
	StatusCanceled   Status = "canceled"
)

func (s Status) String() string {
	return string(s)
}

func (s Status) HumanReadable() string {
	switch s {
	case StatusCreated:
		return "Created"
	case StatusInProgress:
		return "In Progress"
	case StatusReady:
		return "Ready"
	case StatusCompleted:
		return "Completed"
	case StatusCanceled:
		return "Canceled"
	}

	return ""
}

func StatusFromString(s string) (Status, error) {
	status := Status(s)
	switch status {
	case StatusCreated, StatusInProgress, StatusReady, StatusCompleted, StatusCanceled:
		return status, nil
	default:
		return Status("invalid"), fmt.Errorf("invalid status: %s", s)
	}
}
