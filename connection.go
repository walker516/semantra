package dbutil

import (
	"fmt"
	"time"

	"api/config"
	"api/pkg/logutil"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	SleepFunc  = time.Sleep
	MaxRetries = 20
)

// MustPostgres connects with retries & logs failures.
func MustPostgres() *sqlx.DB {
	pg := config.Cfg.Postgres
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone='%s'",
		pg.Host, pg.Port, pg.User, pg.Password, pg.DBName, pg.SSLMode, pg.TimeZone)

	var db *sqlx.DB
	var err error

	for i := 0; i < MaxRetries; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}

		logutil.Error("PG connection attempt failed",
			zap.Int("attempt", i+1),
			zap.Int("maxRetries", MaxRetries),
			zap.Error(err),
		)

		SleepFunc(5 * time.Second)
	}

	if err != nil {
		panic(fmt.Sprintf("PG connection failed after %d retries: %v", MaxRetries, err))
	}

	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}
