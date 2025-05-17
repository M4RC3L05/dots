package core

import (
	"github.com/onsi/gomega"
)

type SpyCallNoRt struct {
	Args []any
}

type SpyCalls struct {
	Debug   []SpyCallNoRt
	Debugnl []SpyCallNoRt
	Error   []SpyCallNoRt
	Errornl []SpyCallNoRt
	Info    []SpyCallNoRt
	Infonl  []SpyCallNoRt
	Log     []SpyCallNoRt
	Lognl   []SpyCallNoRt
	Warn    []SpyCallNoRt
	Warnnl  []SpyCallNoRt
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
	ILogger

	Calls SpyCalls
}

func (sl *SpyLogger) Debug(msg string, args ...any) {
	sl.Calls.Debug = append(sl.Calls.Debug, SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Error(msg string, args ...any) {
	sl.Calls.Error = append(sl.Calls.Error, SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Info(msg string, args ...any) {
	sl.Calls.Info = append(sl.Calls.Info, SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Log(msg string, args ...any) {
	sl.Calls.Log = append(sl.Calls.Log, SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Warn(msg string, args ...any) {
	sl.Calls.Warn = append(sl.Calls.Warn, SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Debugnl(msg string, args ...any) {
	sl.Calls.Debugnl = append(sl.Calls.Debugnl, SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Errornl(msg string, args ...any) {
	sl.Calls.Errornl = append(sl.Calls.Errornl, SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Infonl(msg string, args ...any) {
	sl.Calls.Infonl = append(sl.Calls.Infonl, SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Lognl(msg string, args ...any) {
	sl.Calls.Lognl = append(sl.Calls.Lognl, SpyCallNoRt{Args: append([]any{msg}, args...)})
}

func (sl *SpyLogger) Warnnl(msg string, args ...any) {
	sl.Calls.Warnnl = append(sl.Calls.Warnnl, SpyCallNoRt{Args: append([]any{msg}, args...)})
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
