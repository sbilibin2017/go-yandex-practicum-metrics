package logger

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Error
)

func (l LogLevel) String() string {
	switch l {
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Error:
		return "error"
	default:
		return "debug"
	}
}

type StringBuilder interface {
	WriteString(s string) (n int, err error)
	String() string
}

type DefaultStringBuilder struct {
	strings.Builder
}

type Logger struct {
	level LogLevel
	sb    StringBuilder
}

func NewLogger(level LogLevel, sb StringBuilder) *Logger {
	if sb == nil {
		sb = &DefaultStringBuilder{}
	}
	return &Logger{level: level, sb: sb}
}

func (l *Logger) logMessage(level LogLevel, msg string, extraFields ...interface{}) string {
	if level < l.level {
		return ""
	}

	_, file, line, _ := runtime.Caller(2)
	traceID := generateTraceID()
	timestamp := time.Now().Format(time.RFC3339)

	l.sb.WriteString("{")
	l.sb.WriteString(fmt.Sprintf("\"@timestamp\": \"%s\", ", timestamp))
	l.sb.WriteString(fmt.Sprintf("\"level\": \"%s\", ", level.String()))
	l.sb.WriteString(fmt.Sprintf("\"message\": \"%s\", ", msg))
	l.sb.WriteString(fmt.Sprintf("\"trace_id\": \"%s\", ", traceID))
	l.sb.WriteString(fmt.Sprintf("\"caller\": \"%s:%d\", ", file, line))

	for i := 0; i < len(extraFields); i += 2 {
		key, ok := extraFields[i].(string)
		if !ok || i+1 >= len(extraFields) {
			continue
		}
		value := extraFields[i+1]
		l.sb.WriteString(fmt.Sprintf("\"%s\": \"%v\", ", key, value))
	}

	logMessage := l.sb.String()
	if len(logMessage) > 1 {
		logMessage = strings.TrimSuffix(logMessage, ", ")
	}

	l.sb.WriteString("}")

	fmt.Println(logMessage)
	return logMessage
}

func (l *Logger) Debug(msg string, extraFields ...interface{}) string {
	return l.logMessage(Debug, msg, extraFields...)
}

func (l *Logger) Info(msg string, extraFields ...interface{}) string {
	return l.logMessage(Info, msg, extraFields...)
}

func (l *Logger) Error(msg string, extraFields ...interface{}) string {
	return l.logMessage(Error, msg, extraFields...)
}

func generateTraceID() string {
	return fmt.Sprintf("%x", rand.Int63())
}
