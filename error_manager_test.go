package errmsg_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"api/pkg/errmsg"
	"api/test/testutil"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestRespond_CustomError(t *testing.T) {
	ctx, rec := testutil.NewEchoCtx(http.MethodGet, "/", nil)

	err := errmsg.New(errmsg.CodeNotFound, "feature missing")
	_ = errmsg.Respond(ctx, err)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var body errmsg.ErrorResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	require.Equal(t, errmsg.CodeNotFound, body.Error.Code)
	require.Equal(t, "Resource not found.", body.Error.Message)
	require.Equal(t, "feature missing", body.Error.Details)
}

func TestRespond_GenericError(t *testing.T) {
	ctx, rec := testutil.NewEchoCtx(http.MethodGet, "/", nil)

	_ = errmsg.Respond(ctx, echo.NewHTTPError(http.StatusTeapot, "any error"))

	require.Equal(t, http.StatusInternalServerError, rec.Code)

	var body errmsg.ErrorResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	require.Equal(t, errmsg.CodeInternalServer, body.Error.Code)
}
