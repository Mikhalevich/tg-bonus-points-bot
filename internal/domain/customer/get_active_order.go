package customer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) GetActiveOrder(ctx context.Context, info msginfo.Info) error {
	activeOrder, err := c.repository.GetOrderByChatIDAndStatus(
		ctx,
		info.ChatID,
		order.StatusWaitingPayment,
		order.StatusPaymentInProgress,
		order.StatusConfirmed,
		order.StatusInProgress,
		order.StatusReady,
	)

	if err != nil {
		if c.repository.IsNotFoundError(err) {
			c.sender.ReplyText(ctx, info.ChatID, info.MessageID, "no active orders")
			return nil
		}

		return fmt.Errorf("get order by chat_id: %w", err)
	}

	productsInfo, err := c.repository.GetProductsByIDs(ctx, activeOrder.ProductIDs(), activeOrder.CurrencyID)
	if err != nil {
		return fmt.Errorf("get products by ids: %w", err)
	}

	formattedOrder := formatOrder(activeOrder, productsInfo, c.sender.EscapeMarkdown)

	if activeOrder.CanCancel() {
		cancelBtn, err := c.buttonRepository.SetButton(ctx, button.CancelOrder(info.ChatID, message.Cancel(), activeOrder.ID))
		if err != nil {
			return fmt.Errorf("make cancel order button: %w", err)
		}

		c.sender.ReplyTextMarkdown(ctx, info.ChatID, info.MessageID, formattedOrder, button.InlineRow(cancelBtn))
	} else {
		c.sender.ReplyTextMarkdown(ctx, info.ChatID, info.MessageID, formattedOrder)
	}

	return nil
}

func formatOrder(
	ord *order.Order,
	productsInfo map[product.ProductID]product.Product,
	escaper func(string) string,
) string {
	format := []string{
		fmt.Sprintf("order id: *%s*", escaper(ord.ID.String())),
		fmt.Sprintf("status: *%s*", ord.Status.HumanReadable()),
		fmt.Sprintf("verification code: *%s*", escaper(ord.VerificationCode)),
	}

	for _, t := range ord.Timeline {
		format = append(format, fmt.Sprintf(
			"%s Time: *%s*",
			t.Status.HumanReadable(),
			escaper(t.Time.Format(time.RFC3339))),
		)
	}

	for _, v := range ord.Products {
		format = append(format, fmt.Sprintf("%s x%d %d",
			escaper(productsInfo[v.ProductID].Title), v.Count, v.Price))
	}

	return strings.Join(format, "\n")
}
