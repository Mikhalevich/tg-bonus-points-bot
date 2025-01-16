package customer

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/button"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

func (c *Customer) StartNewCart(ctx context.Context, info msginfo.Info) error {
	categories, err := c.repository.GetCategories(ctx)

	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	cartID, err := c.cart.StartNewCart(ctx, info.ChatID)
	if err != nil {
		return fmt.Errorf("start new cart: %w", err)
	}

	buttons, err := c.makeCartCategoriesButtons(ctx, info.ChatID, cartID, categories, nil)
	if err != nil {
		return fmt.Errorf("make order buttons: %w", err)
	}

	c.sender.ReplyText(ctx, info.ChatID, info.MessageID, message.OrderCategoryPage(), buttons...)

	return nil
}

func (c *Customer) makeCartCategoriesButtons(
	ctx context.Context,
	chatID msginfo.ChatID,
	cartID cart.ID,
	categories []product.Category,
	cartProducts []cart.CartProduct,
) ([]button.InlineKeyboardButtonRow, error) {
	buttons := make([]button.ButtonRow, 0, len(categories)+1)

	for _, v := range categories {
		b, err := button.CartViewCategoryProducts(chatID, v.Title, cartID, v.ID)
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
		makePriceButtonTitle(message.Confirm(), cartProducts),
		cartID,
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

func makePriceButtonTitle(
	caption string,
	cartProducts []cart.CartProduct,
) string {
	price := 0
	for _, v := range cartProducts {
		price += v.Product.Price * v.Count
	}

	if price > 0 {
		return fmt.Sprintf("%s [%d]", caption, price)
	}

	return caption
}
