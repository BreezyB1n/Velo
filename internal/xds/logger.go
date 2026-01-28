package xds

import "log"

// Logger implements go-control-plane's logging interface.
type Logger struct {
	base *log.Logger
}

func NewLogger(base *log.Logger) *Logger {
	return &Logger{base: base}
}

func (l *Logger) Debugf(format string, args ...any) {
	l.base.Printf("[DEBUG] "+format, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.base.Printf("[INFO] "+format, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.base.Printf("[WARN] "+format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.base.Printf("[ERROR] "+format, args...)
}
