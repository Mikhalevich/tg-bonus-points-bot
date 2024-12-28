package perror

type Type int

const (
	TypeInvalid Type = iota
	TypeNotFound
	TypeAlreadyExists
	TypeInvalidParam
)
