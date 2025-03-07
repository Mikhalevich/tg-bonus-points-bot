package cartprocessing

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/internal/message"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/store"
)

func (c *CartProcessing) Create(ctx context.Context, info msginfo.Info) error {
	storeInfo, err := c.storeInfoByID(ctx, c.storeID)
	if err != nil {
		return fmt.Errorf("check for active: %w", err)
	}

	if !storeInfo.IsActive {
		c.sender.SendText(ctx, info.ChatID, storeInfo.ClosedStoreMessage)
		return nil
	}

	categories, err := c.repository.GetCategories(ctx)
	if err != nil {
		return fmt.Errorf("get products: %w", err)
	}

	cartID, err := c.cart.StartNewCart(ctx, info.ChatID)
	if err != nil {
		return fmt.Errorf("start new cart: %w", err)
	}

	curr, err := c.repository.GetCurrencyByID(ctx, storeInfo.Store.DefaultCurrencyID)
	if err != nil {
		return fmt.Errorf("get currency by id: %w", err)
	}

	buttons, err := c.makeCartCategoriesButtons(
		ctx,
		info.ChatID,
		cartID,
		categories,
		nil,
		curr,
	)
	if err != nil {
		return fmt.Errorf("make order buttons: %w", err)
	}

	c.sender.ReplyText(ctx, info.ChatID, info.MessageID, message.OrderCategoryPage(), buttons...)

	return nil
}

type storeInfo struct {
	Store              *store.Store
	IsActive           bool
	ClosedStoreMessage string
}

func (c *CartProcessing) storeInfoByID(ctx context.Context, storeID store.ID) (*storeInfo, error) {
	s, err := c.storeInfo.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("get store by id: %w", err)
	}

	currentTime := c.timeProvider.Now()

	nextWorkingTime, isActive := s.Schedule.NextWorkingTime(currentTime)
	if !isActive {
		return &storeInfo{
			Store:              s,
			IsActive:           false,
			ClosedStoreMessage: message.StoreClosed(currentTime, nextWorkingTime),
		}, nil
	}

	return &storeInfo{
		Store:    s,
		IsActive: true,
	}, nil
}
