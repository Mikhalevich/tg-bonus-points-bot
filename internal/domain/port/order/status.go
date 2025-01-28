package order

import (
	"fmt"
)

type Status string

const (
	StatusWaitingPayment    Status = "waiting_payment"
	StatusPaymentInProgress Status = "payment_in_progress"
	StatusConfirmed         Status = "confirmed"
	StatusInProgress        Status = "in_progress"
	StatusReady             Status = "ready"
	StatusCompleted         Status = "completed"
	StatusCanceled          Status = "canceled"
	StatusRejected          Status = "rejected"
)

func (s Status) String() string {
	return string(s)
}

func (s Status) HumanReadable() string {
	switch s {
	case StatusWaitingPayment:
		return "Waiting Payment"
	case StatusPaymentInProgress:
		return "Payment In Progress"
	case StatusConfirmed:
		return "Confirmed"
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
	case StatusWaitingPayment,
		StatusPaymentInProgress,
		StatusConfirmed,
		StatusInProgress,
		StatusReady,
		StatusCompleted,
		StatusCanceled,
		StatusRejected:
		return status, nil
	default:
		return Status("invalid"), fmt.Errorf("invalid status: %s", s)
	}
}
