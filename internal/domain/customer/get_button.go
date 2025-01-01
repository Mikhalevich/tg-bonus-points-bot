package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/perror"
)

func (c *Customer) GetButton(ctx context.Context, id button.ID) (*button.Button, error) {
	btn, err := c.buttonRepository.GetButton(ctx, id)
	if err != nil {
		if c.buttonRepository.IsNotFoundError(err) {
			return nil, perror.NotFound("button not found")
		}

		return nil, fmt.Errorf("get button from repository: %w", err)
	}

	return btn, nil
}
