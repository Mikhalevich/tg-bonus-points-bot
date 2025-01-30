package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

var (
	stubForCurrencyID = currency.IDFromInt(1)
)

func (c *Customer) CartCreate(ctx context.Context, info msginfo.Info) error {
	categories, err := c.repository.GetCategories(ctx)

	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	cartID, err := c.cart.StartNewCart(ctx, info.ChatID)
	if err != nil {
		return fmt.Errorf("start new cart: %w", err)
	}

	curr, err := c.repository.GetCurrencyByID(ctx, stubForCurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	buttons, err := c.makeCartCategoriesButtons(
		ctx,
		info.ChatID,
		cartID,
		categories,
		nil,
		curr,
	)
	if err != nil {
		return fmt.Errorf("make order buttons: %w", err)
	}

	c.sender.ReplyText(ctx, info.ChatID, info.MessageID, message.OrderCategoryPage(), buttons...)

	return nil
}
