package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Log is the default logger
	Log *logrus.Logger
)

// Config holds the configuration for the logger
type Config struct {
	// LogLevel is the minimum level of severity for logging messages
	LogLevel string
	// LogFile is the file path where logs should be written
	LogFile string
	// EnableConsole determines if logs should also be written to stdout
	EnableConsole bool
	// MaxSize is the maximum size in megabytes of the log file before it gets rotated
	MaxSize int
	// MaxBackups is the maximum number of old log files to retain
	MaxBackups int
	// MaxAge is the maximum number of days to retain old log files
	MaxAge int
	// Compress determines if the rotated log files should be compressed using gzip
	Compress bool
}

// Initialize sets up the logger with the given configuration
func Initialize(config Config) error {
	Log = logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %v", err)
	}
	Log.SetLevel(level)

	// Set log format
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := filepath.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

	// Enable caller information
	Log.SetReportCaller(true)

	// Set output writers
	var writers []io.Writer

	// Add file output if specified
	if config.LogFile != "" {
		// Create logs directory if it doesn't exist
		logDir := filepath.Dir(config.LogFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}

		// Configure log rotation
		rotateLogger := &lumberjack.Logger{
			Filename:   config.LogFile,
			MaxSize:    config.MaxSize,    // megabytes
			MaxBackups: config.MaxBackups, // number of backups
			MaxAge:     config.MaxAge,     // days
			Compress:   config.Compress,   // compress rotated files
		}
		writers = append(writers, rotateLogger)
	}

	// Add console output if enabled
	if config.EnableConsole {
		writers = append(writers, os.Stdout)
	}

	// Set the output to be both file and console if multiple writers
	if len(writers) > 0 {
		Log.SetOutput(io.MultiWriter(writers...))
	}

	Info("Logger initialized", Fields{
		"level":       config.LogLevel,
		"file":        config.LogFile,
		"max_size":    config.MaxSize,
		"max_backups": config.MaxBackups,
		"max_age":     config.MaxAge,
		"compress":    config.Compress,
	})

	return nil
}

// Fields type is an alias for logrus.Fields
type Fields logrus.Fields

// Debug logs a message at level Debug
func Debug(msg string, fields Fields) {
	if fields == nil {
		Log.Debug(msg)
	} else {
		Log.WithFields(logrus.Fields(fields)).Debug(msg)
	}
}

// Info logs a message at level Info
func Info(msg string, fields Fields) {
	if fields == nil {
		Log.Info(msg)
	} else {
		Log.WithFields(logrus.Fields(fields)).Info(msg)
	}
}

// Warn logs a message at level Warn
func Warn(msg string, fields Fields) {
	if fields == nil {
		Log.Warn(msg)
	} else {
		Log.WithFields(logrus.Fields(fields)).Warn(msg)
	}
}

// Error logs a message at level Error
func Error(msg string, fields Fields) {
	if fields == nil {
		Log.Error(msg)
	} else {
		Log.WithFields(logrus.Fields(fields)).Error(msg)
	}
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1
func Fatal(msg string, fields Fields) {
	if fields == nil {
		Log.Fatal(msg)
	} else {
		Log.WithFields(logrus.Fields(fields)).Fatal(msg)
	}
}

// WithError creates an entry from the standard logger and adds an error to it
func WithError(err error) *logrus.Entry {
	return Log.WithError(err)
}

// WithField creates an entry from the standard logger and adds a field to it
func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple fields to it
func WithFields(fields Fields) *logrus.Entry {
	return Log.WithFields(logrus.Fields(fields))
}
