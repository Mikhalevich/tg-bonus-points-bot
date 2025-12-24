package buttonprovider

import (
	"context"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/button"
)

type ButtonRepositoryReader interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
	IsNotFoundError(err error) bool
}

type ButtonProvider struct {
	repository ButtonRepositoryReader
}

func New(repository ButtonRepositoryReader) *ButtonProvider {
	return &ButtonProvider{
		repository: repository,
	}
}
