package cart

import (
	"errors"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/redis/go-redis/v9"
)

var _ port.Cart = (*Cart)(nil)

type Cart struct {
	client *redis.Client
	ttl    time.Duration
}

func New(client *redis.Client, ttl time.Duration) *Cart {
	return &Cart{
		client: client,
		ttl:    ttl,
	}
}

func (c *Cart) IsNotFoundError(err error) bool {
	return errors.Is(err, redis.Nil)
}
