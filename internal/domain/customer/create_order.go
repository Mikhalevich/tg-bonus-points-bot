package customer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) CreateOrder(ctx context.Context, info msginfo.Info, cartID cart.ID) error {
	cartProducts, err := c.cart.GetProducts(ctx, cartID)
	if err != nil {
		if c.cart.IsNotFoundError(err) {
			c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.CartOrderUnavailable())
			return nil
		}

		return fmt.Errorf("get cart products: %w", err)
	}

	if len(cartProducts) == 0 {
		c.sender.SendText(ctx, info.ChatID, message.NoProductsForOrder())
		return nil
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
			c.sender.SendText(ctx, info.ChatID, message.AlreadyHasActiveOrder())
			return nil
		}

		return fmt.Errorf("repository create order: %w", err)
	}

	if err := c.cart.Clear(ctx, info.ChatID, cartID); err != nil {
		return fmt.Errorf("clear cart: %w", err)
	}

	if err := c.sendOrderQRImage(ctx, info, order.Order{
		ID:               id,
		ChatID:           info.ChatID,
		Status:           input.Status,
		VerificationCode: input.VerificationCode,
		Products:         input.Products,
	}); err != nil {
		return fmt.Errorf("send order qr image: %w", err)
	}

	c.sender.DeleteMessage(ctx, info.ChatID, info.MessageID)

	return nil
}

func generateVerificationCode() string {
	//nolint:gosec
	return fmt.Sprintf("%03d", rand.Intn(1000))
}

func (c *Customer) orderProducts(ctx context.Context, cartProducts []port.CartItem) ([]product.ProductCount, error) {
	ids := make([]product.ID, 0, len(cartProducts))

	for _, v := range cartProducts {
		ids = append(ids, v.ProductID)
	}

	productMap, err := c.repository.GetProductsByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get products by ids: %w", err)
	}

	output := make([]product.ProductCount, 0, len(cartProducts))

	for _, v := range cartProducts {
		productInfo, ok := productMap[v.ProductID]
		if !ok {
			return nil, fmt.Errorf("missing product id: %d", v.ProductID.Int())
		}

		output = append(output, product.ProductCount{
			Product: product.Product{
				ID:        v.ProductID,
				Title:     productInfo.Title,
				Price:     productInfo.Price,
				IsEnabled: productInfo.IsEnabled,
				CreatedAt: productInfo.CreatedAt,
				UpdatedAt: productInfo.UpdatedAt,
			},
			Count: v.Count,
		})
	}

	return output, nil
}

func (c *Customer) sendOrderQRImage(ctx context.Context, info msginfo.Info, ord order.Order) error {
	cancelBtn, err := c.buttonRepository.SetButton(ctx, button.CancelOrder(info.ChatID, message.Cancel(), ord.ID))
	if err != nil {
		return fmt.Errorf("cancel order button: %w", err)
	}

	png, err := c.qrCode.GeneratePNG(ord.ID.String())
	if err != nil {
		return fmt.Errorf("qrcode generate png: %w", err)
	}

	if err := c.sender.SendPNGMarkdown(
		ctx,
		info.ChatID,
		formatOrder(&ord, c.sender.EscapeMarkdown),
		png,
		button.InlineRow(cancelBtn),
	); err != nil {
		return fmt.Errorf("send png: %w", err)
	}

	return nil
}
