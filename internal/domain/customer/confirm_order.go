package customer

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) ConfirmOrder(ctx context.Context, info msginfo.Info, cartID cart.ID) error {
	cartItems, err := c.cart.GetProducts(ctx, cartID)
	if err != nil {
		if c.cart.IsNotFoundError(err) {
			c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.CartOrderUnavailable())
			return nil
		}

		return fmt.Errorf("get cart products: %w", err)
	}

	if len(cartItems) == 0 {
		c.sender.SendText(ctx, info.ChatID, message.NoProductsForOrder())
		return nil
	}

	cartProducts, err := c.orderedProducts(ctx, cartItems)
	if err != nil {
		return fmt.Errorf("products info: %w", err)
	}

	input := port.CreateOrderInput{
		ChatID:              info.ChatID,
		Status:              order.StatusWaitingPayment,
		StatusOperationTime: time.Now(),
		VerificationCode:    generateVerificationCode(),
		Products:            cartProducts,
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

	if err := c.sender.SendOrderInvoice(ctx, info.ChatID, message.OrderInvoice(),
		makeOrderDescription(cartProducts),
		convertToOrder(id, info.ChatID, input),
	); err != nil {
		return fmt.Errorf("send order invoice: %w", err)
	}

	c.sender.DeleteMessage(ctx, info.ChatID, info.MessageID)

	return nil
}

func makeOrderDescription(products []order.OrderedProduct) string {
	positions := make([]string, 0, len(products))

	for _, v := range products {
		positions = append(positions, fmt.Sprintf("%s x%d %d", v.Product.Title, v.Count, v.Count*v.Product.Price))
	}

	return strings.Join(positions, "\n")
}

func convertToOrder(id order.ID, chatID msginfo.ChatID, input port.CreateOrderInput) order.Order {
	return order.Order{
		ID:               id,
		ChatID:           chatID,
		Status:           input.Status,
		VerificationCode: input.VerificationCode,
		Products:         input.Products,
	}
}

func generateVerificationCode() string {
	//nolint:gosec
	return fmt.Sprintf("%03d", rand.Intn(1000))
}

func (c *Customer) orderedProducts(
	ctx context.Context,
	cartProducts []cart.CartProduct,
) ([]order.OrderedProduct, error) {
	if len(cartProducts) == 0 {
		return nil, nil
	}

	ids := make([]product.ProductID, 0, len(cartProducts))

	for _, v := range cartProducts {
		ids = append(ids, v.ProductID)
	}

	productMap, err := c.repository.GetProductsByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get products by ids: %w", err)
	}

	output := make([]order.OrderedProduct, 0, len(cartProducts))

	for _, v := range cartProducts {
		productInfo, ok := productMap[v.ProductID]
		if !ok {
			return nil, fmt.Errorf("missing product id: %d", v.ProductID.Int())
		}

		output = append(output, order.OrderedProduct{
			Product:    productInfo,
			CategoryID: v.CategoryID,
			Count:      v.Count,
		})
	}

	return output, nil
}

//nolint:unused
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
