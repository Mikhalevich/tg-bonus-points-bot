package port

import (
	"context"
	"time"
)

type DailyPositionGenerator interface {
	Position(ctx context.Context, t time.Time) (int, error)
}
