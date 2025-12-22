package orderaction

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
)

func (o *OrderAction) GetOrderByID(ctx context.Context, chatID msginfo.ChatID, orderID order.ID) error {
	ord, err := o.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		if o.repository.IsNotFoundError(err) {
			o.sender.SendText(ctx, chatID, message.InvalidOrder())

			return nil
		}

		return fmt.Errorf("get order by chat_id: %w", err)
	}

	if !ord.IsSameChat(chatID) {
		o.sender.SendText(ctx, chatID, message.InvalidOrder())

		return nil
	}

	productsInfo, err := o.repository.GetProductsByIDs(ctx, ord.ProductIDs(), ord.CurrencyID)
	if err != nil {
		return fmt.Errorf("get products by ids: %w", err)
	}

	curr, err := o.repository.GetCurrencyByID(ctx, ord.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	o.sender.SendText(ctx, chatID, formatOrder(ord, curr, productsInfo, 0, o.sender.EscapeMarkdown))

	return nil
}
