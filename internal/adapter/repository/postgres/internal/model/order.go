package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type Order struct {
	ID               int          `db:"id"`
	ChatID           int64        `db:"chat_id"`
	Status           string       `db:"status"`
	VerificationCode string       `db:"verification_code"`
	CreatedAt        time.Time    `db:"created_at"`
	InProgressAt     sql.NullTime `db:"in_progress_at"`
	ReadyAt          sql.NullTime `db:"ready_at"`
	CompletedAt      sql.NullTime `db:"completed_at"`
	CanceledAt       sql.NullTime `db:"canceled_at"`
}

func ToPortOrder(o *Order) (*order.Order, error) {
	status, err := order.StatusFromString(o.Status)
	if err != nil {
		return nil, fmt.Errorf("create status from string: %w", err)
	}

	timeline := []order.StatusTime{
		{
			Status: order.StatusCreated,
			Time:   o.CreatedAt,
		},
	}

	timeline = appendTimelineStatus(timeline, o.InProgressAt, order.StatusInProgress)
	timeline = appendTimelineStatus(timeline, o.ReadyAt, order.StatusReady)
	timeline = appendTimelineStatus(timeline, o.CompletedAt, order.StatusCompleted)
	timeline = appendTimelineStatus(timeline, o.CanceledAt, order.StatusCanceled)

	return &order.Order{
		ID:               order.IDFromInt(o.ID),
		ChatID:           msginfo.ChatIDFromInt(o.ChatID),
		Status:           status,
		VerificationCode: o.VerificationCode,
		Timeline:         timeline,
	}, nil
}

func appendTimelineStatus(
	timeline []order.StatusTime,
	sqlTime sql.NullTime,
	status order.Status,
) []order.StatusTime {
	if !sqlTime.Valid {
		return timeline
	}

	return append(timeline, order.StatusTime{
		Status: status,
		Time:   sqlTime.Time,
	})
}
