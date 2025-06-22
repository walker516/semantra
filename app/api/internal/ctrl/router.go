package ctrl

import (
	"api/pkg/dbutil"

	"github.com/labstack/echo/v4"
)

type Router struct {
	echo     *echo.Echo
	db       *dbutil.DatabaseConnections
}

func NewRouter(e *echo.Echo, dbs *dbutil.DatabaseConnections) *Router {
	return &Router{
		echo: e,
		db:   dbs,
	}
}

func (r *Router) SetupRoutes() {
	r.echo.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	api := r.echo.Group("/api/v1")
	r.setupFeatureRoutes(api)
}

func (r *Router) setupFeatureRoutes(g *echo.Group) {
	// 実装予定: r.handlers.Feature.Get, Create, etc.
}
