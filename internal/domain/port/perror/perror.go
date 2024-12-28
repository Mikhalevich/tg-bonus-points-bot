package perror

import (
	"errors"
)

type Error struct {
	Type    Type
	Message string
}

func New(t Type, msg string) Error {
	return Error{
		Type:    t,
		Message: msg,
	}
}

func ParseError(err error) Error {
	var perr *Error
	if errors.As(err, &perr) {
		return *perr
	}

	return New(TypeInvalid, "internal error")
}

func (e Error) Error() string {
	return e.Message
}

func NotFound(msg string) Error {
	return New(TypeNotFound, msg)
}

func AlreadyExists(msg string) Error {
	return New(TypeAlreadyExists, msg)
}

func InvalidParam(msg string) Error {
	return New(TypeInvalidParam, msg)
}
