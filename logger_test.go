package logutil_test

import (
	"api/pkg/logutil"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestInit_CreatesLogDirAndFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, logutil.Init(dir, zapcore.InfoLevel))

	logutil.Info("test-log")

	_, err := ioutil.ReadFile(filepath.Join(dir, "app.log"))
	require.NoError(t, err)
}

func TestStructuredLogging_Observer(t *testing.T) {
	core, recorded := observer.New(zapcore.DebugLevel)
	logutil.SetLogger(zap.New(core))

	logutil.Info("structured", zap.Int("key", 42))

	entries := recorded.All()
	require.Len(t, entries, 1)
	require.Equal(t, "structured", entries[0].Message)

	val, ok := entries[0].ContextMap()["key"]
	require.True(t, ok)

	// int → int64 で比較
	require.Equal(t, int64(42), val)
}

func TestConvenienceWrappers_AllLevels(t *testing.T) {
	core, observed := observer.New(zapcore.DebugLevel)
	logutil.SetLogger(zap.New(core))

	logutil.Debug("d")
	logutil.Info("i")
	logutil.Warn("w")
	logutil.Error("e")

	expected := []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel}
	entries := observed.All()
	require.Len(t, entries, 4)
	for i, ent := range entries {
		require.Equal(t, expected[i], ent.Level)
	}
}

func TestLogger_WithFields(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	logutil.SetLogger(zap.New(core))

	logutil.L().With(zap.String("user_id", "abc123")).Info("ctx")

	ent := recorded.All()[0]
	require.Equal(t, "ctx", ent.Message)
	require.Equal(t, "abc123", ent.ContextMap()["user_id"])
}
