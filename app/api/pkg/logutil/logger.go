package logutil

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	logDir      = "logs"
	mu          sync.Mutex
	currentDate string
)

// InitLoggers initializes the loggers and manages log rotation.
func InitLoggers() {
	mu.Lock()
	defer mu.Unlock()

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Initialize loggers for the current date
	currentDate = time.Now().Format("2006-01-02")
	setupLoggers(currentDate)

	// Start a goroutine to manage daily log rotation
	go func() {
		for {
			time.Sleep(time.Hour) // Check every hour
			checkAndRotateLogFile()
		}
	}()
}

// setupLoggers sets up info and error loggers for a given date.
func setupLoggers(date string) {
	infoLogFile := createLogFile(filepath.Join(logDir, fmt.Sprintf("info-%s.log", date)))
	errorLogFile := createLogFile(filepath.Join(logDir, fmt.Sprintf("error-%s.log", date)))

	infoLogger = log.New(infoLogFile, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(errorLogFile, "ERROR: ", log.Ldate|log.Ltime)
}

// createLogFile creates or opens a log file for appending.
func createLogFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	return file
}

// checkAndRotateLogFile checks if the date has changed and rotates the log files if necessary.
func checkAndRotateLogFile() {
	mu.Lock()
	defer mu.Unlock()

	newDate := time.Now().Format("2006-01-02")
	if currentDate != newDate {
		currentDate = newDate
		setupLoggers(currentDate)
	}
}

// Info logs an informational message.
func Info(format string, v ...interface{}) {
	message := sanitizeMessage(format, v...)
	infoLogger.Print(message)
}

// Error logs an error message.
func Error(format string, v ...interface{}) {
	message := sanitizeMessage(format, v...)
	errorLogger.Print(message)
}

// Warn logs a warning message (logged to info file).
func Warn(format string, v ...interface{}) {
	message := sanitizeMessage(format, v...)
	infoLogger.Print("WARN: " + message)
}


// sanitizeMessage removes newlines and extra spaces from log messages.
func sanitizeMessage(format string, v ...interface{}) string {
	message := fmt.Sprintf(format, v...)
	message = strings.ReplaceAll(message, "\r\n", " ")
	message = strings.Join(strings.Fields(message), " ")
	return message
}
