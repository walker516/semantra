package dbutil_test

import (
	"api/config"
	"api/pkg/dbutil"
	"api/pkg/logutil"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/stretchr/testify/require"
)

func TestMustPostgres_FailsFastAndLogs(t *testing.T) {
	// Use invalid config to force failure
	config.Cfg.Postgres.Host = "localhost"
	config.Cfg.Postgres.Port = 5432
	config.Cfg.Postgres.User = "invalid"
	config.Cfg.Postgres.Password = "invalid"
	config.Cfg.Postgres.DBName = "invalid"
	config.Cfg.Postgres.SSLMode = "disable"
	config.Cfg.Postgres.TimeZone = "UTC"

	// Zap observer to capture logs
	core, recorded := observer.New(zapcore.DebugLevel)
	logutil.SetLogger(zap.New(core))

	// Patch: disable sleep + 1 retry only
	dbutil.SleepFunc = func(_ time.Duration) {}
	dbutil.MaxRetries = 1

	defer func() {
		dbutil.SleepFunc = time.Sleep
		dbutil.MaxRetries = 20
	}()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustPostgres should panic when connection fails")
		}
	}()

	dbutil.MustPostgres()

	entries := recorded.All()
	require.NotEmpty(t, entries)
	require.Contains(t, entries[0].Message, "PG connection attempt failed")
}
