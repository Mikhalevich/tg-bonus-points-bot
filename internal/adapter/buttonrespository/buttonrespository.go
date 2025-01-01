package buttonrespository

import (
	"errors"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/redis/go-redis/v9"
)

var _ port.ButtonRepository = (*ButtonRepository)(nil)

type ButtonRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func New(client *redis.Client, ttl time.Duration) *ButtonRepository {
	return &ButtonRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *ButtonRepository) IsNotFoundError(err error) bool {
	return errors.Is(err, redis.Nil)
}
