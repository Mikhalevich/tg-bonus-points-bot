package model

import (
	"fmt"
	"sort"
	"time"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/cart"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/product"
)

type Order struct {
	ID               int    `db:"id"`
	ChatID           int64  `db:"chat_id"`
	Status           string `db:"status"`
	VerificationCode string `db:"verification_code"`
}

type OrderTimeline struct {
	ID        int       `db:"order_id"`
	Status    string    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
}

type OrderProduct struct {
	OrderID   int `db:"order_id"`
	ProductID int `db:"product_id"`
	Count     int `db:"count"`
	Price     int `db:"price"`
}

type OrderProductFull struct {
	ProductID int       `db:"product_id"`
	Count     int       `db:"count"`
	Price     int       `db:"price"`
	Title     string    `db:"title"`
	IsEnabled bool      `db:"is_enabled"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func PortToOrderProducts(id order.ID, portProducts []cart.CartProduct) []OrderProduct {
	dbProducts := make([]OrderProduct, 0, len(portProducts))

	for _, v := range portProducts {
		dbProducts = append(dbProducts, OrderProduct{
			OrderID:   id.Int(),
			ProductID: v.Product.ID.Int(),
			Count:     v.Count,
			Price:     v.Product.Price,
		})
	}

	return dbProducts
}

func toPortCartProducts(dbProducts []OrderProductFull) []cart.CartProduct {
	portProducts := make([]cart.CartProduct, 0, len(dbProducts))

	for _, v := range dbProducts {
		portProducts = append(portProducts, cart.CartProduct{
			Product: product.Product{
				ID:        product.ProductIDFromInt(v.ProductID),
				Title:     v.Title,
				Price:     v.Price,
				IsEnabled: v.IsEnabled,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			},
			Count: v.Count,
		})
	}

	return portProducts
}

func ToPortOrder(
	dbOrder *Order,
	dbOrderProducts []OrderProductFull,
	dbTimeline []OrderTimeline,
) (*order.Order, error) {
	orderStatus, err := order.StatusFromString(dbOrder.Status)
	if err != nil {
		return nil, fmt.Errorf("status from string: %w", err)
	}

	portTimeline, err := toPortTimeline(dbTimeline)
	if err != nil {
		return nil, fmt.Errorf("timeline: %w", err)
	}

	return &order.Order{
		ID:               order.IDFromInt(dbOrder.ID),
		ChatID:           msginfo.ChatIDFromInt(dbOrder.ChatID),
		Status:           orderStatus,
		VerificationCode: dbOrder.VerificationCode,
		Timeline:         portTimeline,
		Products:         toPortCartProducts(dbOrderProducts),
	}, nil
}

func toPortTimeline(dbTimeline []OrderTimeline) ([]order.StatusTime, error) {
	sort.Slice(dbTimeline, func(i, j int) bool {
		return dbTimeline[i].UpdatedAt.Sub(dbTimeline[j].UpdatedAt) < 0
	})

	portTimeline := make([]order.StatusTime, 0, len(dbTimeline))

	for _, t := range dbTimeline {
		status, err := order.StatusFromString(t.Status)
		if err != nil {
			return nil, fmt.Errorf("timeline status from string: %w", err)
		}

		portTimeline = append(portTimeline, order.StatusTime{
			Status: status,
			Time:   t.UpdatedAt,
		})
	}

	return portTimeline, nil
}
