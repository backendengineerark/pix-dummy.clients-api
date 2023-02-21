package customerrors

const (
	INVALID_PARAM         = "INVALID_PARAM"
	CLIENT_ALREADY_EXISTS = "CLIENT_ALREADY_EXISTS"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewError(code string, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
