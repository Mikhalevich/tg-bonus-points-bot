package botconsumer

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/botconsumer/tghandler"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
)

func makeRoutes(tbot *tgbot.TGBot, h *tghandler.TGHandler) {
	tbot.AddTextCommand("/start", h.Start)

	tbot.AddMenuCommand("/order", "order food", h.Order)
	tbot.AddMenuCommand("/order_info", "information about active order", h.GetActiveOrder)

	tbot.AddDefaultHandler(h.DefaultHandler)
	tbot.AddDefaultCallbackQueryHander(h.DefaultCallbackQuery)
}
