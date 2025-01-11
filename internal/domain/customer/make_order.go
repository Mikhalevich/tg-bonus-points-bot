package customer

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/flag"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) MakeOrder(ctx context.Context, info msginfo.Info) error {
	categories, err := c.repository.GetCategoryProducts(ctx, product.Filter{
		Products: flag.Enabled,
		Category: flag.Enabled,
	})

	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	buttons, err := c.makeOrderButtons(ctx, info.ChatID, categories)
	if err != nil {
		return fmt.Errorf("make order buttons: %w", err)
	}

	c.sender.ReplyText(ctx, info.ChatID, info.MessageID, message.OrderCategoryPage(), buttons...)

	return nil
}

func generateVerificationCode() string {
	//nolint:gosec
	return fmt.Sprintf("%03d", rand.Intn(1000))
}

func (c *Customer) makeOrderButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	categories []product.Category,
) ([]button.InlineKeyboardButtonRow, error) {
	buttons := make([]button.InlineKeyboardButtonRow, 0, len(categories)+1)

	for _, v := range categories {
		b, err := c.makeInlineKeyboardButton(ctx, button.ViewCategoryProducts(chatID, v.ID), v.Title)
		if err != nil {
			return nil, fmt.Errorf("category order button: %w", err)
		}

		buttons = append(buttons, button.Row(b))
	}

	cancelBtn, err := c.makeInlineKeyboardButton(ctx, button.CancelCart(chatID), message.Cancel())
	if err != nil {
		return nil, fmt.Errorf("cancel order button: %w", err)
	}

	confirmBtn, err := c.makeInlineKeyboardButton(ctx, button.ConfirmCart(chatID), message.Confirm())
	if err != nil {
		return nil, fmt.Errorf("confirm order button: %w", err)
	}

	buttons = append(buttons, button.Row(cancelBtn, confirmBtn))

	return buttons, nil
}
