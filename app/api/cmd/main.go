package main

import (
	"api/config"
	"api/internal/ctrl"
	"api/pkg/dbutil"
	"api/pkg/logutil"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Config initialization error: %v", err)
	}
	logutil.InitLoggers()

	dbs, err := dbutil.OpenPostGISDatabases()
	if err != nil {
		logutil.Error("Database connection error: %v", err)
		return
	}
	defer dbs.Close()

	e := echo.New()
	router := ctrl.NewRouter(e, dbs)
	router.SetupRoutes()

	address := fmt.Sprintf(":%s", config.ConfigData.ServerConfig.Port)
	if err := e.Start(address); err != nil && err != http.ErrServerClosed {
		logutil.Error("Server error: %v", err)
	}
}
