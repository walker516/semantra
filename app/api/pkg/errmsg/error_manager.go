package errmsg

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CustomError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Err     error  `json:"-"`
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %s -> %v", e.Code, e.Message, e.Details, e.Err)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
}

func (e *CustomError) Unwrap() error {
	return e.Err
}

type ErrorResponse struct {
	Error CustomError `json:"error"`
}

const (
	ErrBadRequest     = "ERR_BAD_REQUEST"
	ErrNotFound       = "ERR_NOT_FOUND"
	ErrInternalServer = "ERR_INTERNAL_SERVER"
)

var errorDefs = map[string]struct {
	Status int
	Msg    string
}{
	ErrBadRequest:     {http.StatusBadRequest, "Invalid request parameters."},
	ErrNotFound:       {http.StatusNotFound, "Resource not found."},
	ErrInternalServer: {http.StatusInternalServerError, "Unexpected server error."},
}

func New(code, details string) error {
	def := errorDefs[code]
	return &CustomError{
		Code:    code,
		Message: def.Msg,
		Details: details,
	}
}

func Wrap(code, details string, err error) error {
	def := errorDefs[code]
	return &CustomError{
		Code:    code,
		Message: def.Msg,
		Details: details,
		Err:     err,
	}
}

func Respond(c echo.Context, err error) error {
	var cerr *CustomError
	if !errors.As(err, &cerr) {
		cerr = &CustomError{
			Code:    ErrInternalServer,
			Message: errorDefs[ErrInternalServer].Msg,
			Err:     err,
		}
	}
	status := http.StatusInternalServerError
	if def, ok := errorDefs[cerr.Code]; ok {
		status = def.Status
	}
	return c.JSON(status, ErrorResponse{Error: *cerr})
}
