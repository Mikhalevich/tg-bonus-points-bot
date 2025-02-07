package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

func (c *Customer) OrderPaymentInProgress(
	ctx context.Context,
	paymentID string,
	orderID order.ID,
	currency string,
	totalAmount int,
) error {
	storeInfo, err := c.storeInfoByID(ctx, stubForStoreID)
	if err != nil {
		return fmt.Errorf("check for active: %w", err)
	}

	if !storeInfo.IsActive {
		if err := c.sender.AnswerOrderPayment(ctx, paymentID, false, storeInfo.ClosedStoreMessage); err != nil {
			return fmt.Errorf("answer payment for store is closed: %w", err)
		}

		return nil
	}

	res, err := c.setOrderInProgress(ctx, orderID, totalAmount)
	if err != nil {
		return fmt.Errorf("set order in progress: %w", err)
	}

	if err := c.sender.AnswerOrderPayment(ctx, paymentID, res.OK, res.ErrorMsg); err != nil {
		return fmt.Errorf("answer order payment: %w", err)
	}

	return nil
}

type answerOrderPaymentResult struct {
	OK       bool
	ErrorMsg string
}

func (c *Customer) setOrderInProgress(
	ctx context.Context,
	orderID order.ID,
	totalAmount int,
) (*answerOrderPaymentResult, error) {
	ord, err := c.repository.GetOrderByID(ctx, orderID)
	if err != nil {
		if c.repository.IsNotFoundError(err) {
			return &answerOrderPaymentResult{
				OK:       false,
				ErrorMsg: message.OrderNotExists(),
			}, nil
		}

		return nil, fmt.Errorf("get order by id: %w", err)
	}

	if ord.Status != order.StatusWaitingPayment {
		return &answerOrderPaymentResult{
			OK:       false,
			ErrorMsg: message.OrderStatus(ord.Status),
		}, nil
	}

	if ord.TotalPrice() != totalAmount {
		return &answerOrderPaymentResult{
			OK:       false,
			ErrorMsg: message.OrderTotalPriceIncorrect(),
		}, nil
	}

	if _, err := c.repository.UpdateOrderStatus(
		ctx,
		orderID,
		time.Now(),
		order.StatusPaymentInProgress,
		order.StatusWaitingPayment,
	); err != nil {
		return nil, fmt.Errorf("update order status: %w", err)
	}

	return &answerOrderPaymentResult{
		OK: true,
	}, nil
}
