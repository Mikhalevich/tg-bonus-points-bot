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
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) CartConfirm(
	ctx context.Context,
	info msginfo.Info,
	cartID cart.ID,
	currencyID currency.ID,
) error {
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

	cartProducts, err := c.orderedProducts(ctx, cartItems, currencyID)
	if err != nil {
		return fmt.Errorf("products info: %w", err)
	}

	input := makeCreateOrderInput(info.ChatID, cartProducts, currencyID)

	createdOrder, err := c.repository.CreateOrder(ctx, input)
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

	buttons, err := c.makeInvoiceButtons(ctx, info.ChatID, createdOrder)

	if err != nil {
		return fmt.Errorf("cancel order button: %w", err)
	}

	if err := c.sender.SendOrderInvoice(ctx, info.ChatID, message.OrderInvoice(),
		makeOrderDescription(cartProducts),
		createdOrder,
		buttons...,
	); err != nil {
		return fmt.Errorf("send order invoice: %w", err)
	}

	c.sender.DeleteMessage(ctx, info.ChatID, info.MessageID)

	return nil
}

func makeCreateOrderInput(
	chatID msginfo.ChatID,
	cartProducts []order.OrderedProduct,
	currencyID currency.ID,
) port.CreateOrderInput {
	return port.CreateOrderInput{
		ChatID:              chatID,
		Status:              order.StatusWaitingPayment,
		StatusOperationTime: time.Now(),
		VerificationCode:    generateVerificationCode(),
		Products:            cartProducts,
		CurrencyID:          currencyID,
	}
}

func (c *Customer) makeInvoiceButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	ord *order.Order,
) ([]button.InlineKeyboardButtonRow, error) {
	payBtn := button.Pay(fmt.Sprintf("%s, %s", message.Pay(), ord.TotalPriceHumanReadable()))

	cancelBtn, err := c.buttonRepository.SetButton(ctx, button.CancelOrder(chatID, message.Cancel(), ord.ID))
	if err != nil {
		return nil, fmt.Errorf("cancel order button: %w", err)
	}

	return []button.InlineKeyboardButtonRow{
		button.InlineRow(payBtn),
		button.InlineRow(cancelBtn),
	}, nil
}

func makeOrderDescription(products []order.OrderedProduct) string {
	positions := make([]string, 0, len(products))

	for _, v := range products {
		positions = append(positions, fmt.Sprintf("%s x%d", v.Product.Title, v.Count))
	}

	return strings.Join(positions, ", ")
}

func generateVerificationCode() string {
	//nolint:gosec
	return fmt.Sprintf("%03d", rand.Intn(1000))
}

func (c *Customer) orderedProducts(
	ctx context.Context,
	cartProducts []cart.CartProduct,
	currencyID currency.ID,
) ([]order.OrderedProduct, error) {
	if len(cartProducts) == 0 {
		return nil, nil
	}

	ids := make([]product.ProductID, 0, len(cartProducts))

	for _, v := range cartProducts {
		ids = append(ids, v.ProductID)
	}

	productMap, err := c.repository.GetProductsByIDs(ctx, ids, currencyID)
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
