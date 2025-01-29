package currency

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
