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
	StatusRejected   Status = "rejected"
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
	case StatusRejected:
		return "Rejected"
	}

	return ""
}

func StatusFromString(s string) (Status, error) {
	status := Status(s)
	switch status {
	case StatusCreated, StatusInProgress, StatusReady, StatusCompleted, StatusCanceled, StatusRejected:
		return status, nil
	default:
		return Status("invalid"), fmt.Errorf("invalid status: %s", s)
	}
}
