package customer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) GetActiveOrder(ctx context.Context, chatID msginfo.ChatID, messageID msginfo.MessageID) error {
	activeOrder, err := c.repository.GetOrderByChatIDAndStatus(
		ctx,
		chatID,
		order.StatusCreated,
		order.StatusInProgress,
		order.StatusReady,
	)

	if err != nil {
		if c.repository.IsNotFoundError(err) {
			c.sender.ReplyText(ctx, chatID, messageID, "no active orders")
			return nil
		}

		return fmt.Errorf("get order by chat_id: %w", err)
	}

	formattedOrder := formatOrder(activeOrder, c.sender.EscapeMarkdown)

	if isOrderCancelable(activeOrder.Status) {
		cancelButton := button.CancelOrder(chatID, activeOrder.ID)
		if err := c.buttonRepository.StoreButton(ctx, &cancelButton); err != nil {
			return fmt.Errorf("store cancel button: %w", err)
		}

		c.sender.ReplyTextMarkdown(ctx, chatID, messageID, formattedOrder, cancelOrderButton(cancelButton.ID))
	} else {
		c.sender.ReplyTextMarkdown(ctx, chatID, messageID, formattedOrder)
	}

	return nil
}

func cancelOrderButton(id button.ID) port.Button {
	return port.Button{
		Text: "Cancel",
		ID:   id,
	}
}

func formatOrder(o *order.Order, escaper func(string) string) string {
	format := []string{
		fmt.Sprintf("order id: *%s*", escaper(o.ID.String())),
		fmt.Sprintf("status: *%s*", o.Status.HumanReadable()),
		fmt.Sprintf("verification code: *%s*", escaper(o.VerificationCode)),
	}

	for _, t := range o.Timeline {
		format = append(format, fmt.Sprintf(
			"%s Time: *%s*",
			t.Status.HumanReadable(),
			escaper(t.Time.Format(time.RFC3339))),
		)
	}

	return strings.Join(format, "\n")
}
