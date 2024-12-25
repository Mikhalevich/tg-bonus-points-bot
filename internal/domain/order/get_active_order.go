package order

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (o *Order) GetActiveOrder(ctx context.Context, info msginfo.Info) error {
	activeOrder, err := o.repository.GetOrderByChatIDAndStatus(
		ctx,
		info.ChatID,
		order.StatusCreated,
		order.StatusInProgress,
		order.StatusReady,
	)

	if err != nil {
		if o.repository.IsNotFoundError(err) {
			o.sender.ReplyText(ctx, info.ChatID, info.MessageID, "no active orders")
			return nil
		}

		return fmt.Errorf("get order by chat_id: %w", err)
	}

	o.sender.ReplyTextMarkdown(ctx, info.ChatID, info.MessageID, formatOrder(activeOrder, o.sender.EscapeMarkdown))

	return nil
}

func formatOrder(o *order.Order, escaper func(string) string) string {
	format := []string{
		fmt.Sprintf("order id: *%s*", escaper(o.ID.String())),
		fmt.Sprintf("status: *%s*", formatStatus(o.Status)),
		fmt.Sprintf("verification code: *%s*", escaper(o.VerificationCode)),
	}

	for _, t := range o.Timeline {
		format = append(format, fmt.Sprintf(
			"%s Time: *%s*",
			formatStatus(t.Status),
			escaper(t.Time.Format(time.RFC3339))),
		)
	}

	return strings.Join(format, "\n")
}

func formatStatus(s order.Status) string {
	switch s {
	case order.StatusCreated:
		return "Pending"
	case order.StatusInProgress:
		return "In Progress"
	case order.StatusReady:
		return "Ready"
	case order.StatusCompleted:
		return "Completed"
	case order.StatusCanceled:
		return "Canceled"
	}

	return ""
}
