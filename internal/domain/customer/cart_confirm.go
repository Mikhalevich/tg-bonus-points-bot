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

//nolint:cyclop
func (c *Customer) CartConfirm(
	ctx context.Context,
	info msginfo.Info,
	cartID cart.ID,
	currencyID currency.ID,
) error {
	storeInfo, err := c.storeInfoByID(ctx, c.storeID)
	if err != nil {
		return fmt.Errorf("check for active: %w", err)
	}

	if !storeInfo.IsActive {
		c.sender.SendText(ctx, info.ChatID, storeInfo.ClosedStoreMessage)
		return nil
	}

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

	orderedProducts, productsInfo, err := c.makeOrderedProducts(ctx, cartItems, currencyID)
	if err != nil {
		return fmt.Errorf("products info: %w", err)
	}

	createdOrder, err := c.repository.CreateOrder(ctx, makeCreateOrderInput(info.ChatID, orderedProducts, currencyID))
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

	if err := c.sendOrderInvoice(ctx, info.ChatID, currencyID, createdOrder, productsInfo); err != nil {
		return fmt.Errorf("send order invoice: %w", err)
	}

	c.sender.DeleteMessage(ctx, info.ChatID, info.MessageID)

	return nil
}

func (c *Customer) sendOrderInvoice(
	ctx context.Context,
	chatID msginfo.ChatID,
	currencyID currency.ID,
	createdOrder *order.Order,
	productsInfo map[product.ProductID]product.Product,
) error {
	curr, err := c.repository.GetCurrencyByID(ctx, currencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	buttons, err := c.makeInvoiceButtons(ctx, chatID, createdOrder, curr)

	if err != nil {
		return fmt.Errorf("cancel order button: %w", err)
	}

	if err := c.sender.SendOrderInvoice(ctx, chatID, message.OrderInvoice(),
		makeOrderDescription(createdOrder.Products, productsInfo),
		createdOrder,
		productsInfo,
		curr.Code,
		buttons...,
	); err != nil {
		return fmt.Errorf("send order invoice: %w", err)
	}

	return nil
}

func makeCreateOrderInput(
	chatID msginfo.ChatID,
	orderedProducts []order.OrderedProduct,
	currencyID currency.ID,
) port.CreateOrderInput {
	return port.CreateOrderInput{
		ChatID:              chatID,
		Status:              order.StatusWaitingPayment,
		StatusOperationTime: time.Now(),
		VerificationCode:    generateVerificationCode(),
		Products:            orderedProducts,
		CurrencyID:          currencyID,
	}
}

func (c *Customer) makeInvoiceButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	ord *order.Order,
	curr *currency.Currency,
) ([]button.InlineKeyboardButtonRow, error) {
	payBtn := button.Pay(fmt.Sprintf("%s, %s", message.Pay(), curr.FormatPrice(ord.TotalPrice())))

	cancelBtn, err := button.CancelOrder(chatID, message.Cancel(), ord.ID, false)
	if err != nil {
		return nil, fmt.Errorf("cancel order button: %w", err)
	}

	inlineCancelBtn, err := c.buttonRepository.SetButton(ctx, cancelBtn)
	if err != nil {
		return nil, fmt.Errorf("cancel order button: %w", err)
	}

	return []button.InlineKeyboardButtonRow{
		button.InlineRow(payBtn),
		button.InlineRow(inlineCancelBtn),
	}, nil
}

func makeOrderDescription(
	orderedProducts []order.OrderedProduct,
	productsInfo map[product.ProductID]product.Product,
) string {
	positions := make([]string, 0, len(orderedProducts))

	for _, v := range orderedProducts {
		positions = append(positions, fmt.Sprintf("%s x%d", productsInfo[v.ProductID].Title, v.Count))
	}

	return strings.Join(positions, ", ")
}

func generateVerificationCode() string {
	//nolint:gosec
	return fmt.Sprintf("%03d", rand.Intn(1000))
}

func (c *Customer) makeOrderedProducts(
	ctx context.Context,
	cartProducts []cart.CartProduct,
	currencyID currency.ID,
) ([]order.OrderedProduct, map[product.ProductID]product.Product, error) {
	if len(cartProducts) == 0 {
		return nil, nil, nil
	}

	ids := make([]product.ProductID, 0, len(cartProducts))

	for _, v := range cartProducts {
		ids = append(ids, v.ProductID)
	}

	productMap, err := c.repository.GetProductsByIDs(ctx, ids, currencyID)
	if err != nil {
		return nil, nil, fmt.Errorf("get products by ids: %w", err)
	}

	output := make([]order.OrderedProduct, 0, len(cartProducts))

	for _, v := range cartProducts {
		productInfo, ok := productMap[v.ProductID]
		if !ok {
			return nil, nil, fmt.Errorf("missing product id: %d", v.ProductID.Int())
		}

		output = append(output, order.OrderedProduct{
			ProductID:  v.ProductID,
			CategoryID: v.CategoryID,
			Count:      v.Count,
			Price:      productInfo.Price,
		})
	}

	return output, productMap, nil
}
