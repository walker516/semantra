package authmw

import (
	"api/pkg/errmsg"
	"api/pkg/logutil"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RequireUserID is a middleware that checks if usr_id is present in query parameters
func RequireUserID() echo.MiddlewareFunc {
	fmt.Println("requre id")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			usrID := c.QueryParam("usr_id")
			if usrID == "" {
				return errmsg.Respond(c, errmsg.New("ERR_UNAUTHORIZED", "usr_id is required"))
			}
			logutil.Info("Request authenticated for user: %s",
				zap.String("user", usrID))
			return next(c)
		}
	}
}
