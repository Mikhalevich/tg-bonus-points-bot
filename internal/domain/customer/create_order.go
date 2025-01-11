package customer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) CreateOrder(ctx context.Context, info msginfo.Info) error {
	cartProducts, err := c.cart.GetProducts(ctx, info.ChatID)
	if err != nil {
		if c.cart.IsNotFoundError(err) {
			c.sender.SendText(ctx, info.ChatID, message.NoProductsForOrder())
			return nil
		}

		return fmt.Errorf("get cart products: %w", err)
	}

	orderProducts, err := c.orderProducts(ctx, cartProducts)
	if err != nil {
		return fmt.Errorf("products info: %w", err)
	}

	input := port.CreateOrderInput{
		ChatID:              info.ChatID,
		Status:              order.StatusConfirmed,
		StatusOperationTime: time.Now(),
		VerificationCode:    generateVerificationCode(),
		Products:            orderProducts,
	}

	id, err := c.repository.CreateOrder(ctx, input)

	if err != nil {
		if c.repository.IsAlreadyExistsError(err) {
			c.sender.ReplyText(ctx, info.ChatID, info.MessageID, message.AlreadyHasActiveOrder())
			return nil
		}

		return fmt.Errorf("repository create order: %w", err)
	}

	if err := c.cart.Clear(ctx, info.ChatID); err != nil {
		return fmt.Errorf("clear cart: %w", err)
	}

	cancelBtn, err := c.makeInlineKeyboardButton(ctx, button.CancelOrder(info.ChatID, id), message.Cancel())
	if err != nil {
		return fmt.Errorf("cancel order button: %w", err)
	}

	png, err := c.qrCode.GeneratePNG(id.String())
	if err != nil {
		return fmt.Errorf("qrcode generate png: %w", err)
	}

	if err := c.sender.SendPNGMarkdown(
		ctx,
		info.ChatID,
		formatOrder(&order.Order{
			ID:               id,
			ChatID:           info.ChatID,
			Status:           input.Status,
			VerificationCode: input.VerificationCode,
		}, c.sender.EscapeMarkdown),
		png,
		button.Row(cancelBtn),
	); err != nil {
		return fmt.Errorf("send png: %w", err)
	}

	c.sender.DeleteMessage(ctx, info.ChatID, info.MessageID)

	return nil
}

func generateVerificationCode() string {
	//nolint:gosec
	return fmt.Sprintf("%03d", rand.Intn(1000))
}

func (c *Customer) orderProducts(ctx context.Context, cartProducts []port.CartItem) ([]product.ProductCount, error) {
	products := make([]product.ProductCount, 0, len(cartProducts))

	for _, v := range cartProducts {
		products = append(products, product.ProductCount{
			Product: product.Product{
				ID:    v.ProductID,
				Title: "product test title",
				Price: 100,
			},
			Count: 1,
		})
	}

	return products, nil
}
