package orderaction

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/messageprocessor/button"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/product"
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/infra/logger"
)

func (o *OrderAction) GetActiveOrder(ctx context.Context, info msginfo.Info) error {
	activeOrder, err := o.repository.GetOrderByChatIDAndStatus(
		ctx,
		info.ChatID,
		order.StatusWaitingPayment,
		order.StatusPaymentInProgress,
		order.StatusConfirmed,
		order.StatusInProgress,
		order.StatusReady,
	)

	if err != nil {
		if o.repository.IsNotFoundError(err) {
			o.replyPlainText(ctx, info.ChatID, info.MessageID, "no active orders")

			return nil
		}

		return fmt.Errorf("get order by chat_id: %w", err)
	}

	productsInfo, err := o.repository.GetProductsByIDs(ctx, activeOrder.ProductIDs(), activeOrder.CurrencyID)
	if err != nil {
		return fmt.Errorf("get products by ids: %w", err)
	}

	curr, err := o.repository.GetCurrencyByID(ctx, activeOrder.CurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	position := o.orderQueuePosition(ctx, activeOrder)

	if err := o.replyCancelOrderMessage(
		ctx,
		info.ChatID,
		info.MessageID,
		activeOrder,
		curr,
		productsInfo,
		position,
	); err != nil {
		return fmt.Errorf("cancel order reply: %w", err)
	}

	return nil
}

func (o *OrderAction) orderQueuePosition(ctx context.Context, activeOrder *order.Order) int {
	if !activeOrder.InQueue() {
		return 0
	}

	pos, err := o.repository.GetOrderPositionByStatus(
		ctx,
		activeOrder.ID,
		order.StatusConfirmed,
		order.StatusInProgress,
	)

	if err != nil {
		if o.repository.IsNotFoundError(err) {
			return 0
		}

		logger.FromContext(ctx).WithError(err).Error("failed to get order position")

		return 0
	}

	return pos
}

func (o *OrderAction) replyCancelOrderMessage(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	activeOrder *order.Order,
	curr *currency.Currency,
	productsInfo map[product.ProductID]product.Product,
	queuePosition int,
) error {
	formattedOrder := formatOrder(activeOrder, curr, productsInfo, queuePosition, o.sender.EscapeMarkdown)

	if !activeOrder.CanCancel() {
		o.replyMarkdown(ctx, chatID, messageID, formattedOrder)

		return nil
	}

	cancelBtn, err := button.CancelOrder(chatID, message.Cancel(), activeOrder.ID, true)
	if err != nil {
		return fmt.Errorf("cancel order button: %w", err)
	}

	o.replyMarkdown(ctx, chatID, messageID, formattedOrder, button.Row(cancelBtn))

	return nil
}

func formatOrder(
	ord *order.Order,
	curr *currency.Currency,
	productsInfo map[product.ProductID]product.Product,
	queuePosition int,
	escaper func(string) string,
) string {
	format := []string{
		fmt.Sprintf("order id: *%s*", escaper(ord.ID.String())),
		fmt.Sprintf("status: *%s*", ord.Status.HumanReadable()),
		fmt.Sprintf("verification code: *%s*", escaper(ord.VerificationCode)),
		fmt.Sprintf("daily position: *%d*", ord.DailyPosition),
		fmt.Sprintf("total price: *%s*", curr.FormatPrice(ord.TotalPrice)),
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
		format = append(format, fmt.Sprintf("%s x%d %s",
			escaper(productsInfo[v.ProductID].Title), v.Count, curr.FormatPrice(v.Price)))
	}

	if queuePosition > 0 {
		format = append(format, fmt.Sprintf("position in queue: *%d*", queuePosition))
	}

	return strings.Join(format, "\n")
}
