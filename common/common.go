package common

type ContextKey string

const RequestIDKey ContextKey = "requestID"

const UserIDKey ContextKey = "userID"

const ErrorKey ContextKey = "error"

const (
	UNABLE_TO_SAVE          = "UNABLE_TO_SAVE"
	UNABLE_TO_FIND_RESOURCE = "UNABLE_TO_FIND_RESOURCE"
	UNABLE_TO_READ          = "UNABLE_TO_READ"
	UNAUTHORIZED            = "UNAUTHORIZED"
)

type CustomError struct {
	Type       string
	Message    string
	StatusCode int
}

// func (e *CustomError) Error() string {
// 	return fmt.Sprintf("%s: %s: %d", e.Type, e.Message, e.StatusCode)
// }
