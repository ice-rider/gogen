package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	level      Level
	mu         sync.Mutex
	writer     io.Writer
	logFile    *os.File
	enableFile bool
}

func NewLogger(level Level, enableFile bool) *Logger {
	l := &Logger{
		level:      level,
		writer:     os.Stdout,
		enableFile: enableFile,
	}

	if enableFile {
		if err := l.openLogFile(); err != nil {
			log.Printf("Warning: failed to open log file: %v", err)
		}
	}

	return l
}

func (l *Logger) openLogFile() error {
	logPath := ".generator.log"

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	l.logFile = f
	return nil
}

func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.log("DEBUG", format, args...)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.log("INFO", format, args...)
	}
}

func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.log("WARN", format, args...)
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= LevelError {
		l.log("ERROR", format, args...)
	}
}

func (l *Logger) Success(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(format, args...)

	fmt.Fprintf(l.writer, "✓ %s\n", msg)

	if l.logFile != nil {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(l.logFile, "[%s] SUCCESS: %s\n", timestamp, msg)
	}
}

func (l *Logger) log(level, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	fmt.Fprintf(l.writer, "[%s] %s: %s\n", timestamp, level, msg)

	if l.logFile != nil {
		fmt.Fprintf(l.logFile, "[%s] %s: %s\n", timestamp, level, msg)
	}
}

func (l *Logger) Section(title string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	separator := "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

	fmt.Fprintf(l.writer, "\n%s\n", separator)
	fmt.Fprintf(l.writer, "%s\n", title)
	fmt.Fprintf(l.writer, "%s\n\n", separator)

	if l.logFile != nil {
		fmt.Fprintf(l.logFile, "\n=== %s ===\n\n", title)
	}
}
