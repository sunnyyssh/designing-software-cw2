package errs

import "fmt"

type HTTPError struct {
	Code    int
	Message string
}

func (e HTTPError) Error() string { return e.Message }

func BadRequest(format string, args ...any) HTTPError {
	return HTTPError{
		Code:    400,
		Message: fmt.Sprintf(format, args...),
	}
}

func NotFound(format string, args ...any) HTTPError {
	return HTTPError{
		Code:    404,
		Message: fmt.Sprintf(format, args...),
	}
}
