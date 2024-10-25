package logging

import (
	"io"
	"log"
	"sync"
)

type LogLevel int
type (
	ILogger interface {
		Debug(message string, args ...any)
		Info(message string, args ...any)
		Warning(message string, args ...any)
		Error(message string, args ...any)
		SetOutput(io.Writer)
		SetFlags(int)
	}
	SimpleLogger struct {
		lock          sync.RWMutex
		level         LogLevel
		debugLogger   *log.Logger
		infoLogger    *log.Logger
		warningLogger *log.Logger
		errorLogger   *log.Logger
	}
)

// Debug implements ILogger.
func (l *SimpleLogger) Debug(message string, args ...any) {
	if l.level&DebugLogLevel == 1 {
		l.lock.RLock()
		l.debugLogger.Printf(message, args...)
		l.lock.RUnlock()
	}
}

// Error implements ILogger.
func (l *SimpleLogger) Error(message string, args ...any) {
	if l.level&ErrorLogLevel == 1 {
		l.lock.RLock()
		l.errorLogger.Printf(message, args...)
		l.lock.RUnlock()
	}
}

// Info implements ILogger.
func (l *SimpleLogger) Info(message string, args ...any) {
	if l.level&InfoLogLevel == 1 {
		l.lock.RLock()
		l.infoLogger.Printf(message, args...)
		l.lock.RUnlock()
	}
}

// Warning implements ILogger.
func (l *SimpleLogger) Warning(message string, args ...any) {
	if l.level&WarningLogLevel == 1 {
		l.lock.RLock()
		l.warningLogger.Printf(message, args...)
		l.lock.RUnlock()
	}
}

func (l *SimpleLogger) SetOutput(w io.Writer) {
	l.lock.Lock()
	l.debugLogger.SetOutput(w)
	l.infoLogger.SetOutput(w)
	l.warningLogger.SetOutput(w)
	l.errorLogger.SetOutput(w)
	l.lock.Unlock()
}

func (l *SimpleLogger) SetFlags(flags int) {
	l.lock.Lock()
	l.debugLogger.SetFlags(flags)
	l.infoLogger.SetFlags(flags)
	l.warningLogger.SetFlags(flags)
	l.errorLogger.SetFlags(flags)
	l.lock.Unlock()
}

const (
	DebugLogLevel   LogLevel = 1 << iota
	InfoLogLevel             // 2
	WarningLogLevel          // 4
	ErrorLogLevel            // 8
)

// LogLevels returns all log levels which can be passed to the New(...) function
func LogLevels() LogLevel {
	return DebugLogLevel | InfoLogLevel | WarningLogLevel | ErrorLogLevel
}

// New creates and returns a new basic logger interface which exposes a few base methods from the `log.Logger` structure
// that can be set for logs later on such as flags or writers
func New(writer io.Writer, level LogLevel) ILogger {
	return &SimpleLogger{
		debugLogger:   log.New(writer, "DEBUG", log.Flags()),
		infoLogger:    log.New(writer, "INFO", log.Flags()),
		warningLogger: log.New(writer, "WARNING", log.Flags()),
		errorLogger:   log.New(writer, "ERROR", log.Flags()),
		level:         level,
	}
}
