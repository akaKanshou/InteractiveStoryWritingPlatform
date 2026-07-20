package errors

type DBError struct {
	code            int
	responseError   string
	underlyingError error
}

func NewDBError(code int, responseError string, underlyingError error) *DBError {
	return &DBError{
		code:            code,
		responseError:   responseError,
		underlyingError: underlyingError,
	}
}

func (e *DBError) ResponseError() string {
	return e.responseError
}

func (e *DBError) Code() int {
	return e.code
}

func (e *DBError) Error() string {
	return e.underlyingError.Error()
}
