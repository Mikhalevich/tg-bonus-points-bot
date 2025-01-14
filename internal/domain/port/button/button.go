package button

import (
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port/msginfo"
)

type ID string

func (id ID) String() string {
	return string(id)
}

func IDFromString(s string) ID {
	return ID(s)
}

type Button struct {
	ChatID    msginfo.ChatID
	Caption   string
	Operation Operation
	Payload   []byte
}

type ButtonRow []Button

func Row(buttons ...Button) ButtonRow {
	return buttons
}
