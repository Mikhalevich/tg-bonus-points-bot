package buttonprovider

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

type ButtonProvider struct {
	repository port.ButtonRepositoryReader
}

func New(repository port.ButtonRepositoryReader) *ButtonProvider {
	return &ButtonProvider{
		repository: repository,
	}
}
