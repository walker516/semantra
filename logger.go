package logutil

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var std *zap.Logger // global production logger

func init() { std = zap.NewNop() }

// Init configures the global logger.
//
//	dir   – log directory (default "logs")
//	level – zapcore.Level (zapcore.InfoLevel, zapcore.DebugLevel,…)
//
// It emits **JSON** to both stdout *and* a rotated file (app.log).
func Init(dir string, level zapcore.Level) error {
	if dir == "" {
		dir = "logs"
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "time"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEnc := zapcore.NewJSONEncoder(encCfg)

	fileSync := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(dir, "app.log"),
		MaxSize:    50, // MB
		MaxBackups: 7,
		MaxAge:     30, // days
		Compress:   true,
	})

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEnc, zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(jsonEnc, fileSync, level),
	)

	// std = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	std = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	return nil
}

// L returns the production *zap.Logger*.
func L() *zap.Logger { return std }

// S returns a SugaredLogger for ergonomic key‑value calls.
func S() *zap.SugaredLogger { return std.Sugar() }

// TEST ONLY: override the global logger.
func SetLogger(l *zap.Logger) { std = l }

// Thin wrappers ------------------------------------------------------------
func Debug(msg string, fields ...zap.Field) { std.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { std.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { std.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { std.Error(msg, fields...) }
