package buttonprovider

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/perror"
)

func (b *ButtonProvider) GetButton(ctx context.Context, id button.ID) (*button.Button, error) {
	btn, err := b.repository.GetButton(ctx, id)
	if err != nil {
		if b.repository.IsNotFoundError(err) {
			return nil, perror.NotFound("button not found")
		}

		return nil, fmt.Errorf("get button from repository: %w", err)
	}

	return btn, nil
}
