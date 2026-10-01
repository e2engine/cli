package errors

type UserError struct {
	Message string
	Err     error
}

func NewUserError(message string, err error) *UserError {
	return &UserError{
		Message: message,
		Err:     err,
	}
}

func (e *UserError) Error() string {
	return e.Message
}

func (e *UserError) Unwrap() error {
	return e.Err
}
