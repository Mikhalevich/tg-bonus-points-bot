package currency

import (
	"fmt"
	"math"
)

type ID int

func (id ID) Int() int {
	return int(id)
}

func IDFromInt(id int) ID {
	return ID(id)
}

type Currency struct {
	ID         ID
	Code       string
	Exp        int
	DecimalSep string
	MinAmount  int
	MaxAmount  int
	IsEnabled  bool
}

func (c Currency) FormatPrice(price int) string {
	if c.Exp == 0 {
		return fmt.Sprintf("%d %s", price, c.Code)
	}

	var (
		div = int(math.Pow10(c.Exp))
		rub = price / div
		kop = price % div
	)

	return fmt.Sprintf("%d%s%0*d %s", rub, c.DecimalSep, c.Exp, kop, c.Code)
}
