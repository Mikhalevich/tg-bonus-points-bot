package order

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
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

func IDFromString(id string) (ID, error) {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse int: %w", err)
	}

	return ID(intID), nil
}

type Status string

func (s Status) String() string {
	return string(s)
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

const (
	StatusCreated    Status = "created"
	StatusInProgress Status = "in_progress"
	StatusReady      Status = "ready"
	StatusCompleted  Status = "completed"
	StatusCanceled   Status = "canceled"
)

type Order struct {
	ID               ID
	ChatID           msginfo.ChatID
	Status           Status
	VerificationCode string
	Timeline         []StatusTime
}

type StatusTime struct {
	Status Status
	Time   time.Time
}
