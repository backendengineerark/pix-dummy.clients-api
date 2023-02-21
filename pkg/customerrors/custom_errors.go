package customerrors

const (
	INVALID = "INVALID"
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
