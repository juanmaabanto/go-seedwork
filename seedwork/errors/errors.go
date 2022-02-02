package errors

type ErrorType struct {
	t string
}

var (
	ErrorTypeUnknown    = ErrorType{"unknown"}
	ErrorTypeBadRequest = ErrorType{"bad-request"}
	ErrorTypeNotFound   = ErrorType{"not-found"}
	ErrorTypeValidation = ErrorType{"Validation Failure"}
)

type ApplicationError struct {
	message   string
	title     string
	errorType ErrorType
	errors    map[string]string
}

func (app ApplicationError) Error() string {
	return app.message
}

func (app ApplicationError) Errors() map[string]string {
	return app.errors
}

func (app ApplicationError) Title() string {
	return app.title
}

func (app ApplicationError) ErrorType() ErrorType {
	return app.errorType
}

func NewApplicationError(error string, title string) ApplicationError {
	return ApplicationError{
		message:   error,
		title:     title,
		errorType: ErrorTypeUnknown,
	}
}

func NewBadRequestError(error string) ApplicationError {
	return ApplicationError{
		message:   error,
		title:     "Bad Request",
		errorType: ErrorTypeBadRequest,
	}
}

func NewNotFoundError(error string) ApplicationError {
	return ApplicationError{
		message:   error,
		title:     "Not Found",
		errorType: ErrorTypeNotFound,
	}
}

func NewValidationError(errors map[string]string) ApplicationError {
	return ApplicationError{
		message:   "One or more validation errors occurred.",
		title:     "Validation Failure",
		errorType: ErrorTypeValidation,
		errors:    errors,
	}
}
