package botconsumer

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/botconsumer/tghandler"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
)

func makeRoutes(tbot *tgbot.TGBot, handler *tghandler.TGHandler) {
	tbot.AddTextCommand("/start", handler.Start)

	tbot.AddMenuCommand("/order", "order food", handler.Order)
	tbot.AddMenuCommand("/order_info", "information about active order", handler.GetActiveOrder)
	tbot.AddMenuCommand("/queue", "current orders queue size", handler.OrderQueueSize)
	tbot.AddMenuCommand("/history", "view history orders", handler.OrderHistory)

	tbot.AddDefaultHandler(handler.DefaultHandler)
	tbot.AddDefaultCallbackQueryHander(handler.DefaultCallbackQuery)
}
