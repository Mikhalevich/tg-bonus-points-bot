package page

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/order"
)

type Page struct {
	Number int
	Total  int
}

func Current(
	history []order.HistoryOrder,
	totalOrders int,
	direction ScrollDirection,
	pageSize int,
	ordering Ordering,
) Page {
	var (
		cp = currentPage(history, direction, pageSize)
		tp = totalPages(totalOrders, pageSize)
	)

	if ordering == DESC {
		cp = tp - cp + 1
	}

	return Page{
		Number: cp,
		Total:  tp,
	}
}

func First(totalOrders int, pageSize int, ordering Ordering) Page {
	if ordering == ASC {
		return firstPage(totalOrders, pageSize)
	}

	return lastPage(totalOrders, pageSize)
}

func Last(totalOrders int, pageSize int, ordering Ordering) Page {
	if ordering == ASC {
		return lastPage(totalOrders, pageSize)
	}

	return firstPage(totalOrders, pageSize)
}

func firstPage(totalOrders int, pageSize int) Page {
	return Page{
		Number: 1,
		Total:  totalPages(totalOrders, pageSize),
	}
}

func lastPage(totalOrders int, pageSize int) Page {
	total := totalPages(totalOrders, pageSize)

	return Page{
		Number: total,
		Total:  total,
	}
}

func totalPages(totalOrders int, pageSize int) int {
	return totalOrders/pageSize + 1
}

func currentPage(history []order.HistoryOrder, direction ScrollDirection, pageSize int) int {
	if direction == BottomToTop {
		return (history[len(history)-1].SerialNumber-1)/pageSize + 1
	}

	return (history[0].SerialNumber-1)/pageSize + 1
}

func CalculateOrderIDForNextPage(twoPageOrders []order.HistoryOrder, pageSize int) order.ID {
	if len(twoPageOrders) > pageSize {
		return twoPageOrders[pageSize-1].ID
	}

	return 0
}

func TruncateOrdersToPageSize(orders []order.HistoryOrder, pageSize int) []order.HistoryOrder {
	if len(orders) > pageSize {
		return orders[:pageSize]
	}

	return orders
}
