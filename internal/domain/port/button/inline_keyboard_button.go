package button

type InlineKeyboardButton struct {
	ID      ID
	Caption string
}

type InlineKeyboardButtonRow []InlineKeyboardButton

func Row(buttons ...InlineKeyboardButton) InlineKeyboardButtonRow {
	return buttons
}
