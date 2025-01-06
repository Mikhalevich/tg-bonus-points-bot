package model

import (
	"fmt"
	"sort"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type Order struct {
	ID               int    `db:"id"`
	ChatID           int64  `db:"chat_id"`
	Status           string `db:"status"`
	VerificationCode string `db:"verification_code"`
}

type OrderTimeline struct {
	ID        int       `db:"order_id"`
	Status    string    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ToPortOrder(dbOrder *Order, dbTimeline []OrderTimeline) (*order.Order, error) {
	orderStatus, err := order.StatusFromString(dbOrder.Status)
	if err != nil {
		return nil, fmt.Errorf("status from string: %w", err)
	}

	sort.Slice(dbTimeline, func(i, j int) bool {
		return dbTimeline[i].UpdatedAt.Sub(dbTimeline[j].UpdatedAt) < 0
	})

	portTimeline := make([]order.StatusTime, 0, len(dbTimeline))

	for _, t := range dbTimeline {
		status, err := order.StatusFromString(t.Status)
		if err != nil {
			return nil, fmt.Errorf("timeline status from string: %w", err)
		}

		portTimeline = append(portTimeline, order.StatusTime{
			Status: status,
			Time:   t.UpdatedAt,
		})
	}

	return &order.Order{
		ID:               order.IDFromInt(dbOrder.ID),
		ChatID:           msginfo.ChatIDFromInt(dbOrder.ChatID),
		Status:           orderStatus,
		VerificationCode: dbOrder.VerificationCode,
		Timeline:         portTimeline,
	}, nil
}
