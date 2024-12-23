package model

import (
	"database/sql"
	"time"
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
