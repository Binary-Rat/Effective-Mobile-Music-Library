package appErr

type Error struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func New(statusCode int, message string) *Error {
	return &Error{
		StatusCode: statusCode,
		Message:    message,
	}
}
