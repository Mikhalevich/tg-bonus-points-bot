package customer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/flag"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) MakeOrder(ctx context.Context, info msginfo.Info) error {
	input := port.CreateOrderInput{
		ChatID:              info.ChatID,
		Status:              order.StatusAssembling,
		StatusOperationTime: time.Now(),
		VerificationCode:    generateVerificationCode(),
	}

	id, err := c.repository.CreateOrder(ctx, input)

	if err != nil {
		if c.repository.IsAlreadyExistsError(err) {
			c.sender.ReplyText(ctx, info.ChatID, info.MessageID, message.AlreadyHasActiveOrder())
			return nil
		}

		return fmt.Errorf("repository create order: %w", err)
	}

	categories, err := c.repository.GetCategoryProducts(ctx, product.Filter{
		Products: flag.Enabled,
		Category: flag.Enabled,
	})
	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	buttons, err := c.makeOrderButtons(ctx, info.ChatID, id, categories)
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
	orderID order.ID,
	categories []product.Category,
) ([]button.InlineKeyboardButtonRow, error) {
	buttons := make([]button.InlineKeyboardButtonRow, 0, len(categories)+1)

	for _, v := range categories {
		b, err := c.makeInlineKeyboardButton(ctx, button.ViewCategory(chatID, orderID, v.ID), v.Title)
		if err != nil {
			return nil, fmt.Errorf("category order button: %w", err)
		}

		buttons = append(buttons, button.Row(b))
	}

	cancelBtn, err := c.makeInlineKeyboardButton(ctx, button.CancelOrderEditMsg(chatID, orderID), message.Cancel())
	if err != nil {
		return nil, fmt.Errorf("cancel order button: %w", err)
	}

	confirmBtn, err := c.makeInlineKeyboardButton(ctx, button.ConfirmOrder(chatID, orderID), message.Confirm())
	if err != nil {
		return nil, fmt.Errorf("confirm order button: %w", err)
	}

	buttons = append(buttons, button.Row(cancelBtn, confirmBtn))

	return buttons, nil
}
