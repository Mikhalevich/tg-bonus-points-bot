package botconsumer

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/botconsumer/tghandler"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/app/internal/tgbot"
)

func makeRoutes(tbot *tgbot.TGBot, h *tghandler.TGHandler) {
	tbot.AddTextCommand("/start", h.Start)

	tbot.AddMenuCommand("/order", "make order", h.Order)
	tbot.AddMenuCommand("/get_active_order", "retrieve active order", h.GetActiveOrder)
	tbot.AddMenuCommand("/test_products", "test products", h.TestProducts)

	tbot.AddDefaultTextHandler(h.DefaultText)
	tbot.AddDefaultCallbackQueryHander(h.DefaultCallbackQuery)
}
