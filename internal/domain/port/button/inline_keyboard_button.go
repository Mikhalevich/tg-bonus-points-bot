package button

type InlineKeyboardButton struct {
	ID      ID
	Caption string
	Pay     bool
}

type InlineKeyboardButtonRow []InlineKeyboardButton

func InlineRow(buttons ...InlineKeyboardButton) InlineKeyboardButtonRow {
	return buttons
}

func Pay(caption string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Caption: caption,
		Pay:     true,
	}
}
