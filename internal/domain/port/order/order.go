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

type Order struct {
	ID               ID
	ChatID           msginfo.ChatID
	Status           Status
	VerificationCode string
	Timeline         []StatusTime
}

func (o Order) IsSameChat(chatID msginfo.ChatID) bool {
	return o.ChatID == chatID
}

func (o Order) CanCancel() bool {
	return o.Status == StatusConfirmed
}

type StatusTime struct {
	Status Status
	Time   time.Time
}
