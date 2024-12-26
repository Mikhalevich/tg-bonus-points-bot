package response

import (
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type Order struct {
	ID               string       `json:"id"`
	Status           string       `json:"status"`
	VerificationCode string       `json:"verification_code"`
	Timeline         []StatusTime `json:"timeline"`
}

type StatusTime struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func ToOrder(o *order.Order) *Order {
	timeline := make([]StatusTime, 0, len(o.Timeline))
	for _, tl := range o.Timeline {
		timeline = append(timeline, StatusTime{
			Status: tl.Status.String(),
			Time:   tl.Time,
		})
	}

	return &Order{
		ID:               o.ID.String(),
		Status:           o.ID.String(),
		VerificationCode: o.VerificationCode,
		Timeline:         timeline,
	}
}
