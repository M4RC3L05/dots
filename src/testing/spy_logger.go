package testing

import (
	"github.com/m4rc3l05/dots/src/core"
	"github.com/onsi/gomega"
)

type SpyLoggerCalls struct {
	Debug   []core.SpyCallNoRt
	Debugnl []core.SpyCallNoRt
	Error   []core.SpyCallNoRt
	Errornl []core.SpyCallNoRt
	Info    []core.SpyCallNoRt
	Infonl  []core.SpyCallNoRt
	Log     []core.SpyCallNoRt
	Lognl   []core.SpyCallNoRt
	Warn    []core.SpyCallNoRt
	Warnnl  []core.SpyCallNoRt
}

type SpyLoggerCallNumber struct {
	Debug   int
	Debugnl int
	Error   int
	Errornl int
	Info    int
	Infonl  int
	Log     int
	Lognl   int
	Warn    int
	Warnnl  int
}

type SpyLogger struct {
	core.ILogger

	Calls SpyLoggerCalls
}

func (sl *SpyLogger) Debug(msg string, args ...any) {
	sl.Calls.Debug = append(sl.Calls.Debug, core.SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Error(msg string, args ...any) {
	sl.Calls.Error = append(sl.Calls.Error, core.SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Info(msg string, args ...any) {
	sl.Calls.Info = append(sl.Calls.Info, core.SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Log(msg string, args ...any) {
	sl.Calls.Log = append(sl.Calls.Log, core.SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Warn(msg string, args ...any) {
	sl.Calls.Warn = append(sl.Calls.Warn, core.SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Debugnl(msg string, args ...any) {
	sl.Calls.Debugnl = append(sl.Calls.Debugnl, core.SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Errornl(msg string, args ...any) {
	sl.Calls.Errornl = append(sl.Calls.Errornl, core.SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Infonl(msg string, args ...any) {
	sl.Calls.Infonl = append(sl.Calls.Infonl, core.SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Lognl(msg string, args ...any) {
	sl.Calls.Lognl = append(sl.Calls.Lognl, core.SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Warnnl(msg string, args ...any) {
	sl.Calls.Warnnl = append(sl.Calls.Warnnl, core.SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func MakeSpyLogger() *SpyLogger {
	return &SpyLogger{}
}

func AssertSpyLoggerCalls(logger SpyLogger, callNumber *SpyLoggerCallNumber) {
	var callNumberVal SpyLoggerCallNumber

	if callNumber != nil {
		callNumberVal = *callNumber
	} else {
		callNumberVal = SpyLoggerCallNumber{}
	}

	gomega.Expect(logger.Calls.Debug).To(gomega.HaveLen(callNumberVal.Debug))
	gomega.Expect(logger.Calls.Error).To(gomega.HaveLen(callNumberVal.Error))
	gomega.Expect(logger.Calls.Info).To(gomega.HaveLen(callNumberVal.Info))
	gomega.Expect(logger.Calls.Log).To(gomega.HaveLen(callNumberVal.Log))
	gomega.Expect(logger.Calls.Warn).To(gomega.HaveLen(callNumberVal.Warn))

	gomega.Expect(logger.Calls.Debugnl).To(gomega.HaveLen(callNumberVal.Debugnl))
	gomega.Expect(logger.Calls.Errornl).To(gomega.HaveLen(callNumberVal.Errornl))
	gomega.Expect(logger.Calls.Infonl).To(gomega.HaveLen(callNumberVal.Infonl))
	gomega.Expect(logger.Calls.Lognl).To(gomega.HaveLen(callNumberVal.Lognl))
	gomega.Expect(logger.Calls.Warnnl).To(gomega.HaveLen(callNumberVal.Warnnl))
}
