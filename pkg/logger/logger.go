// Package logger provides a simple logging interface using zap
package logger

import (
	"go.uber.org/zap"
)

// Logger is a simple wrapper around zap.SugaredLogger
type Logger struct {
	*zap.SugaredLogger
}

// Global logger instance
var globalLogger *Logger

// Initialize the global logger
func init() {
	// Create a minimal production logger that only shows important information
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.DisableStacktrace = true
	config.OutputPaths = []string{"stdout"}
	config.Encoding = "console"
	
	zapLogger, _ := config.Build()
	globalLogger = &Logger{zapLogger.Sugar()}
}

// New returns a new logger with fields
func New() *Logger {
	return globalLogger
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{l.With(key, value)}
}

// Error logs an error message
func (l *Logger) Error(err error, msg string, args ...interface{}) {
	// Always include the error
	l.With("error", err).Errorw(msg, args...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.SugaredLogger.Sync()
}
