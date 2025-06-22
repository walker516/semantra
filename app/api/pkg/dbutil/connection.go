package dbutil

import (
	"api/config"
	"api/pkg/logutil"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DatabaseConnections struct {
	PostGIS *sqlx.DB
}

func OpenPostGISDatabases() (*DatabaseConnections, error) {
	cfg := config.ConfigData.DatabaseConfigs["postgis"]

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
		cfg.TimeZone,
	)

	var db *sqlx.DB
	var err error
	const maxRetries = 10
	const retryDelay = 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}
		logutil.Error("PostGIS connection failed (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostGIS after retries: %w", err)
	}

	configureDatabase(db)
	return &DatabaseConnections{PostGIS: db}, nil
}

func (d *DatabaseConnections) Close() {
	if d.PostGIS != nil {
		d.PostGIS.Close()
	}
}

func configureDatabase(db *sqlx.DB) {
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
}
