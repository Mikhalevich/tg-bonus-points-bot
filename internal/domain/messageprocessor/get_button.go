package messageprocessor

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/perror"
)

func (m *MessageProcessor) GetButton(ctx context.Context, id button.ID) (*button.Button, error) {
	btn, err := m.buttonRepository.GetButton(ctx, id)
	if err != nil {
		if m.buttonRepository.IsNotFoundError(err) {
			return nil, perror.NotFound("button not found")
		}

		return nil, fmt.Errorf("get button from repository: %w", err)
	}

	return btn, nil
}
