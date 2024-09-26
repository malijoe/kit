package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type statusError int

var (
	ErrBadRequest statusError = http.StatusBadRequest
	ErrNotFound   statusError = http.StatusNotFound
	ErrConflict   statusError = http.StatusConflict
	ErrNoContent  statusError = http.StatusNoContent
	ErrInternal   statusError = http.StatusInternalServerError
)

func (se statusError) Error() string {
	return http.StatusText(int(se))
}

func StatusError(status int, msg string) error {
	return fmt.Errorf("%w: %s", statusError(status), msg)
}

func StatusErrorf(status int, format string, args ...any) error {
	return StatusError(status, fmt.Sprintf(format, args...))
}

func JoinStatusError(status int, errs ...error) error {
	return errors.Join(append([]error{statusError(status)}, errs...)...)
}

func NotFoundError(msg string) error {
	return StatusError(http.StatusBadRequest, msg)
}

func NotFoundErrorf(format string, args ...any) error {
	return StatusErrorf(http.StatusNotFound, format, args...)
}

func JoinNotFoundError(errs ...error) error {
	return JoinStatusError(http.StatusNotFound, errs...)
}

func BadRequestError(msg string) error {
	return StatusError(http.StatusBadRequest, msg)
}

func BadRequestErrorf(format string, args ...any) error {
	return StatusErrorf(http.StatusBadRequest, format, args...)
}

func JoinBadRequestError(errs ...error) error {
	return JoinStatusError(http.StatusBadRequest, errs...)
}
