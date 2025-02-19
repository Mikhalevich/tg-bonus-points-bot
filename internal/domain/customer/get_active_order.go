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
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/infra/logger"
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

	position := c.orderQueuePosition(ctx, activeOrder)

	if err := c.replyCancelOrderMessage(
		ctx,
		info.ChatID,
		info.MessageID,
		activeOrder,
		productsInfo,
		position,
	); err != nil {
		return fmt.Errorf("cancel order reply: %w", err)
	}

	return nil
}

func (c *Customer) orderQueuePosition(ctx context.Context, activeOrder *order.Order) int {
	if !activeOrder.InQueue() {
		return 0
	}

	pos, err := c.repository.GetOrderPositionByStatus(
		ctx,
		activeOrder.ID,
		order.StatusConfirmed,
		order.StatusInProgress,
	)

	if err != nil {
		if c.repository.IsNotFoundError(err) {
			return 0
		}

		logger.FromContext(ctx).WithError(err).Error("failed to get order position")

		return 0
	}

	return pos
}

func (c *Customer) replyCancelOrderMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	activeOrder *order.Order,
	productsInfo map[product.ProductID]product.Product,
	queuePosition int,
) error {
	formattedOrder := formatOrder(activeOrder, productsInfo, queuePosition, c.sender.EscapeMarkdown)

	if !activeOrder.CanCancel() {
		c.sender.ReplyTextMarkdown(ctx, chatID, messageID, formattedOrder)
		return nil
	}

	cancelBtn, err := button.CancelOrder(chatID, message.Cancel(), activeOrder.ID, true)
	if err != nil {
		return fmt.Errorf("cancel order button: %w", err)
	}

	inlineCancelBtn, err := c.buttonRepository.SetButton(ctx, cancelBtn)
	if err != nil {
		return fmt.Errorf("make cancel order button: %w", err)
	}

	c.sender.ReplyTextMarkdown(ctx, chatID, messageID, formattedOrder, button.InlineRow(inlineCancelBtn))

	return nil
}

func formatOrder(
	ord *order.Order,
	productsInfo map[product.ProductID]product.Product,
	queuePosition int,
	escaper func(string) string,
) string {
	format := []string{
		fmt.Sprintf("order id: *%s*", escaper(ord.ID.String())),
		fmt.Sprintf("status: *%s*", ord.Status.HumanReadable()),
		fmt.Sprintf("verification code: *%s*", escaper(ord.VerificationCode)),
		fmt.Sprintf("daily position: *%d*", ord.DailyPosition),
		fmt.Sprintf("created\\_at: *%s*", escaper(ord.CreatedAt.Format(time.RFC3339))),
		fmt.Sprintf("updated\\_at: *%s*", escaper(ord.UpdatedAt.Format(time.RFC3339))),
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

	if queuePosition > 0 {
		format = append(format, fmt.Sprintf("position in queue: *%d*", queuePosition))
	}

	return strings.Join(format, "\n")
}
