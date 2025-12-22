package request

import (
	"github.com/Mikhalevich/tg-coffee-shop-bot/internal/domain/port/flag"
)

type Filter string

const (
	FilterEnabled  Filter = "enabled"
	FilterDisabled Filter = "disabled"
	FilterAll      Filter = "all"
)

func (f Filter) ToPortState() flag.State {
	switch f {
	case FilterEnabled:
		return flag.Enabled
	case FilterDisabled:
		return flag.Disabled
	case FilterAll:
		return flag.All
	}

	return flag.All
}
