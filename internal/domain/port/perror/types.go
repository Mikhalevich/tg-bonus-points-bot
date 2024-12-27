package perror

type Type int

const (
	TypeNotFound Type = iota + 1
	TypeAlreadyExists
)
