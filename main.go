package main

import (
	"api/config"
	"api/internal/route"
	"api/pkg/dbutil"
	"api/pkg/logutil"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Load config
	if err := config.Load("config.dev.json"); err != nil {
		panic(err)
	}

	// Init logger → logs/ + INFO level
	if err := logutil.Init("logs", zapcore.DebugLevel); err != nil {
		panic(err)
	}
	logutil.Info("Logger initialized")

	db := dbutil.MustPostgres()
	defer db.Close()

	e := echo.New()
	rt := route.NewRouter(e, db)
	rt.SetupRoutes()

	go gracefulShutdown(e)

	port := config.Cfg.Server.Port
	logutil.Info("Starting server",
		zap.String("port", port),
	)

	if err := e.Start(":" + port); err != nil {
		logutil.Error("Echo server error",
			zap.Error(err),
		)
	}

}

func gracefulShutdown(e *echo.Echo) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logutil.Info("Shutting down server...")

	// Echo Shutdown expects a context; best practice is to give a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logutil.Error("Server forced to shutdown",
			zap.Error(err),
		)
	}
}
