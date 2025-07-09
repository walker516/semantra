package route

import (
	"api/internal/ctrl"
	"api/internal/flow"
	authmw "api/internal/mw/auth"
	corsmw "api/internal/mw/cors"
	"api/internal/repo"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Router struct {
	echo       *echo.Echo
	db         *sqlx.DB
	ctrls      *ctrl.Ctrls
	middleware *Middleware
}

type Middleware struct {
	Auth echo.MiddlewareFunc
	CORS echo.MiddlewareFunc
}

func NewRouter(e *echo.Echo, db *sqlx.DB) *Router {
	return &Router{
		echo: e,
		db:   db,
	}
}

func (r *Router) SetupRoutes() {
	r.setupMiddleware()
	r.setupHandlers()
	r.registerRoutes()
}

func (r *Router) setupMiddleware() {
	r.middleware = &Middleware{
		Auth: authmw.RequireUserID(),
		CORS: corsmw.NewCORSConfig(),
	}
	r.echo.Use(r.middleware.CORS)
}

func (r *Router) setupHandlers() {
	repos := repo.NewRepositories(r.db)
	flow := flow.NewFlows(repos)
	r.ctrls = ctrl.NewCtrls(flow)
}

func (r *Router) registerRoutes() {
	r.echo.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	// バージョン付きルート
	api := r.echo.Group("/api/v1")

	// ------------------------------
	// Feature Resource
	// ------------------------------
	api.OPTIONS("/features", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })
	api.OPTIONS("/features/:id", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	api.GET("/features/:id", r.ctrls.Feature.GetByID)
	api.POST("/features", r.ctrls.Feature.Create)
	api.PATCH("/features/:id", r.ctrls.Feature.Update)
	api.DELETE("/features/:id", r.ctrls.Feature.Delete)
	// ------------------------------
	// Tag Resource
	// ------------------------------
	// api.OPTIONS("/tags", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	// api.GET("/tags", r.ctrls.Tag.List)
}
