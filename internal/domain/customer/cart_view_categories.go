package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/currency"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) CartViewCategories(
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

		return fmt.Errorf("cart items: %w", err)
	}

	orderedProducts, _, err := c.makeOrderedProducts(ctx, cartItems, currencyID)
	if err != nil {
		return fmt.Errorf("cart products: %w", err)
	}

	categories, err := c.repository.GetCategories(ctx)
	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	curr, err := c.repository.GetCurrencyByID(ctx, currencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	buttons, err := c.makeCartCategoriesButtons(
		ctx,
		info.ChatID,
		cartID,
		categories,
		orderedProducts,
		curr,
	)
	if err != nil {
		return fmt.Errorf("make order buttons: %w", err)
	}

	c.sender.EditTextMessage(ctx, info.ChatID, info.MessageID, message.OrderCategoryPage(), buttons...)

	return nil
}

func (c *Customer) makeCartCategoriesButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	cartID cart.ID,
	categories []product.Category,
	orderedProducts []order.OrderedProduct,
	curr *currency.Currency,
) ([]button.InlineKeyboardButtonRow, error) {
	buttons := make([]button.ButtonRow, 0, len(categories)+1)

	for _, v := range categories {
		title := makeViewCategoryButtonTitle(v, orderedProducts, curr)

		b, err := button.CartViewCategoryProducts(chatID, title, cartID, v.ID, curr.ID)
		if err != nil {
			return nil, fmt.Errorf("create cart view button: %w", err)
		}

		buttons = append(buttons, button.Row(b))
	}

	cancelCartBtn, err := button.CartCancel(chatID, message.Cancel(), cartID)
	if err != nil {
		return nil, fmt.Errorf("cancel cart button: %w", err)
	}

	confirmCartBtn, err := button.CartConfirm(
		chatID,
		makePriceButtonTitle(orderedProducts, curr),
		cartID,
		curr.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("confirm cart button: %w", err)
	}

	buttons = append(buttons, []button.Button{
		cancelCartBtn,
		confirmCartBtn,
	})

	inlineKeyboardButtonRows, err := c.buttonRepository.SetButtonRows(ctx, buttons...)
	if err != nil {
		return nil, fmt.Errorf("set button rows: %w", err)
	}

	return inlineKeyboardButtonRows, nil
}

func makeViewCategoryButtonTitle(
	category product.Category,
	orderedProducts []order.OrderedProduct,
	curr *currency.Currency,
) string {
	var (
		count int
		price int
	)

	for _, v := range orderedProducts {
		if category.ID == v.CategoryID {
			price += v.Price * v.Count
			count += v.Count
		}
	}

	if count > 0 {
		return fmt.Sprintf("%s [x%d %s]", category.Title, count, curr.FormatPrice(price))
	}

	return category.Title
}

func makePriceButtonTitle(
	orderedProducts []order.OrderedProduct,
	curr *currency.Currency,
) string {
	price := 0
	for _, v := range orderedProducts {
		price += v.Price * v.Count
	}

	if price > 0 {
		return fmt.Sprintf("%s [%s]", message.Confirm(), curr.FormatPrice(price))
	}

	return message.Confirm()
}
