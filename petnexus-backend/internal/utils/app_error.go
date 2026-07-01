package utils

// AppError is a safe service error that handlers can translate into the
// standard PetNexus error response.
type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Details    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap preserves the internal cause for logs and errors.Is/errors.As.
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError builds an application error without exposing its internal cause
// in the HTTP response.
func NewAppError(status int, code, message, details string, cause error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
		Details:    details,
		Err:        cause,
	}
}
