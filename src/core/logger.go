package core

import (
	"fmt"

	"github.com/fatih/color"
)

type ILogger interface {
	Debug(msg string, args ...any)
	Debugnl(msg string, args ...any)
	Info(msg string, args ...any)
	Infonl(msg string, args ...any)
	Warn(msg string, args ...any)
	Warnnl(msg string, args ...any)
	Error(msg string, args ...any)
	Errornl(msg string, args ...any)
	Log(msg string, args ...any)
	Lognl(msg string, args ...any)
}

type logger struct {
	ILogger
}

var (
	logLevelDebug = "DBG"
	logLevelInfo  = "INF"
	logLevelWarn  = "WRN"
	logLevelError = "ERR"
)

var logLevelColorMap = map[string]func(format string, a ...any) string{
	logLevelDebug: color.MagentaString,
	logLevelInfo:  color.CyanString,
	logLevelWarn:  color.YellowString,
	logLevelError: color.RedString,
}

func resolveLogLevelStr(level *string) string {
	if level == nil {
		return ""
	}

	colorMap, colorMapExists := logLevelColorMap[*level]

	if !colorMapExists {
		colorMap = func(format string, a ...any) string {
			return fmt.Sprintf(format, a)
		}
	}

	return colorMap(fmt.Sprintf("%s: ", *level))
}

func log(level *string, nl bool, msg string, args ...any) {
	levelStr := resolveLogLevelStr(level)

	fmt.Print(levelStr)
	fmt.Printf(msg, args...)

	if nl {
		fmt.Println("")
	}
}

func (logger) Debug(msg string, args ...any) {
	log(&logLevelDebug, false, msg, args...)
}

func (logger) Info(msg string, args ...any) {
	log(&logLevelInfo, false, msg, args...)
}

func (logger) Warn(msg string, args ...any) {
	log(&logLevelWarn, false, msg, args...)
}

func (logger) Error(msg string, args ...any) {
	log(&logLevelError, false, msg, args...)
}

func (logger) Log(msg string, args ...any) {
	log(nil, false, msg, args...)
}

func (logger) Debugnl(msg string, args ...any) {
	log(&logLevelDebug, true, msg, args...)
}

func (logger) Infonl(msg string, args ...any) {
	log(&logLevelInfo, true, msg, args...)
}

func (logger) Warnnl(msg string, args ...any) {
	log(&logLevelWarn, true, msg, args...)
}

func (logger) Errornl(msg string, args ...any) {
	log(&logLevelError, true, msg, args...)
}

func (logger) Lognl(msg string, args ...any) {
	log(nil, true, msg, args...)
}

func MakeLogger() ILogger {
	return logger{}
}
