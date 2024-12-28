package perror

type Error struct {
	Type    Type
	Message string
}

func New(t Type, msg string) *Error {
	return &Error{
		Type:    t,
		Message: msg,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func NotFound(msg string) *Error {
	return New(TypeNotFound, msg)
}

func AlreadyExists(msg string) *Error {
	return New(TypeAlreadyExists, msg)
}

func InvalidParam(msg string) *Error {
	return New(TypeInvalidParam, msg)
}
