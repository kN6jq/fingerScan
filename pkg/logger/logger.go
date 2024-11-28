package logger

import (
	"fmt"
	"github.com/gookit/color"
	"log"
	"os"
)

var (
	colorError   = "237,64,35"  // 红色
	colorWarning = "244,211,49" // 黄色
	colorInfo    = "0,128,0"    // 绿色

	defaultLogger *Logger
)

func init() {
	defaultLogger = NewLogger(false)
}

type Logger struct {
	stdLogger *log.Logger
	debug     bool
}

func NewLogger(debug bool) *Logger {
	return &Logger{
		stdLogger: log.New(os.Stdout, "", log.LstdFlags),
		debug:     debug,
	}
}

// Global functions
func Error(format string, v ...interface{}) {
	defaultLogger.Error(format, v...)
}

func Warning(format string, v ...interface{}) {
	defaultLogger.Warning(format, v...)
}

func Info(format string, v ...interface{}) {
	defaultLogger.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	defaultLogger.Debug(format, v...)
}

// Logger methods
func (l *Logger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	color.RGBStyleFromString(colorError).Println("[ERROR] " + msg)
	if l.debug {
		l.stdLogger.Printf("[ERROR] %s\n", msg)
	}
}

func (l *Logger) Warning(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	color.RGBStyleFromString(colorWarning).Println("[WARNING] " + msg)
	if l.debug {
		l.stdLogger.Printf("[WARNING] %s\n", msg)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	color.RGBStyleFromString(colorInfo).Println("[INFO] " + msg)
	if l.debug {
		l.stdLogger.Printf("[INFO] %s\n", msg)
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.debug {
		msg := fmt.Sprintf(format, v...)
		l.stdLogger.Printf("[DEBUG] %s\n", msg)
	}
}

func (l *Logger) SetDebug(debug bool) {
	l.debug = debug
}
