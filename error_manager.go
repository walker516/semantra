package errmsg

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Code is a machine‑readable error code.
// Define as typed string for extra type‑safety.
type Code string

const (
	CodeBadRequest     Code = "ERR_BAD_REQUEST"
	CodeUnauthorized   Code = "ERR_UNAUTHORIZED"
	CodeNotFound       Code = "ERR_NOT_FOUND"
	CodeInternalServer Code = "ERR_INTERNAL_SERVER"
)

// mapping from Code → (status,msg).  Central place for i18n later.
var def = map[Code]struct {
	status int
	msg    string
}{
	CodeBadRequest:     {http.StatusBadRequest, "Invalid request parameters."},
	CodeUnauthorized:   {http.StatusUnauthorized, "Authentication required."},
	CodeNotFound:       {http.StatusNotFound, "Resource not found."},
	CodeInternalServer: {http.StatusInternalServerError, "Unexpected server error."},
}

// CustomError implements error and carries structured fields for JSON API.
// It is *always* created through New / Wrap so that Code is never empty.
type CustomError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Err     error  `json:"-"`
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %s → %v", e.Code, e.Message, e.Details, e.Err)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
}
func (e *CustomError) Unwrap() error { return e.Err }

// JSON wire‑format returned to clients.
type ErrorResponse struct {
	Error *CustomError `json:"error"`
}

// New returns CustomError with default message looked‑up by Code.
func New(code Code, details string) *CustomError {
	d := def[code]
	return &CustomError{Code: code, Message: d.msg, Details: details}
}

// Wrap enriches an underlying error while preserving stack.
func Wrap(code Code, details string, err error) *CustomError {
	d := def[code]
	return &CustomError{Code: code, Message: d.msg, Details: details, Err: err}
}

// HTTPStatus maps Code→ status.
func (c Code) HTTPStatus() int { return def[c].status }

// Respond uses Echo to emit JSON API errors.  Always returns nil so that
// Echo does not attempt a second write.
func Respond(c echo.Context, err error) error {
	var cust *CustomError
	if !errors.As(err, &cust) {
		// fall‑back – keep original error string for server log, but do not expose it.
		cust = New(CodeInternalServer, "")
	}
	return c.JSON(cust.Code.HTTPStatus(), ErrorResponse{cust})
}
