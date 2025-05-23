package cartprocessing

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/perror"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *CartProcessing) Confirm(
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

	orderedProducts, productsInfo, err := c.orderedProductsFromCart(ctx, cartID, currencyID)
	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.CartOrderUnavailable())

			return nil
		}

		return fmt.Errorf("ordered products from cart: %w", err)
	}

	if len(orderedProducts) == 0 {
		c.sender.SendText(ctx, info.ChatID, message.NoProductsForOrder())

		return nil
	}

	createdOrder, err := c.repository.CreateOrder(ctx, c.makeCreateOrderInput(info.ChatID, orderedProducts, currencyID))
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

func (c *CartProcessing) sendOrderInvoice(
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

func (c *CartProcessing) makeCreateOrderInput(
	chatID msginfo.ChatID,
	orderedProducts []order.OrderedProduct,
	currencyID currency.ID,
) port.CreateOrderInput {
	totalPrice := 0
	for _, v := range orderedProducts {
		totalPrice += v.Count * v.Price
	}

	return port.CreateOrderInput{
		ChatID:              chatID,
		Status:              order.StatusWaitingPayment,
		StatusOperationTime: c.timeProvider.Now(),
		VerificationCode:    "",
		TotalPrice:          totalPrice,
		Products:            orderedProducts,
		CurrencyID:          currencyID,
	}
}

func (c *CartProcessing) makeInvoiceButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	ord *order.Order,
	curr *currency.Currency,
) ([]button.InlineKeyboardButtonRow, error) {
	payBtn := button.Pay(fmt.Sprintf("%s, %s", message.Pay(), curr.FormatPrice(ord.TotalPrice)))

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

func (c *CartProcessing) orderedProductsFromCart(
	ctx context.Context,
	cartID cart.ID,
	currencyID currency.ID,
) ([]order.OrderedProduct, map[product.ProductID]product.Product, error) {
	cartProducts, err := c.cart.GetProducts(ctx, cartID)
	if err != nil {
		if c.cart.IsNotFoundError(err) {
			return nil, nil, perror.NotFound("cart not found")
		}

		return nil, nil, fmt.Errorf("get cart products: %w", err)
	}

	if len(cartProducts) == 0 {
		return nil, nil, nil
	}

	productIDs := make([]product.ProductID, 0, len(cartProducts))
	for _, v := range cartProducts {
		productIDs = append(productIDs, v.ProductID)
	}

	productsInfo, err := c.repository.GetProductsByIDs(ctx, productIDs, currencyID)
	if err != nil {
		return nil, nil, fmt.Errorf("get products by ids: %w", err)
	}

	output := make([]order.OrderedProduct, 0, len(cartProducts))

	for _, prod := range cartProducts {
		productInfo, ok := productsInfo[prod.ProductID]
		if !ok {
			return nil, nil, fmt.Errorf("missing product id: %d", prod.ProductID.Int())
		}

		output = append(output, order.OrderedProduct{
			ProductID:  prod.ProductID,
			CategoryID: prod.CategoryID,
			Count:      prod.Count,
			Price:      productInfo.Price,
		})
	}

	return output, productsInfo, nil
}
