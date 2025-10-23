package writer

import (
	"log"

	"github.com/995933447/fastlog/logger"
	"github.com/995933447/fastlog/logger/fmts"
)

var _ logger.Writer = (*StdoutWriter)(nil)

func NewStdoutWriter(level logger.Level, modName string, skipCall int) *StdoutWriter {
	return &StdoutWriter{
		level: level,
		fmt: fmts.NewTraceFormatter(
			modName,
			skipCall,
			fmts.FormatText,
			false,
			false,
			func() int32 {
				return 0
			},
			func() int32 {
				return 0
			},
		),
	}
}

type StdoutWriter struct {
	level logger.Level
	fmt   logger.Formatter
}

func (w *StdoutWriter) IsLoggable(level logger.Level) bool {
	return w.level <= level
}

func (w *StdoutWriter) DisableCacheCaller(disabled bool) {
	w.fmt.DisableCacheCaller(disabled)
}

func (w *StdoutWriter) EnableStdoutPrinter() {
}

func (w *StdoutWriter) DisableStdoutPrinter() {
}

func (w *StdoutWriter) Write(level logger.Level, args ...interface{}) error {
	logContent, err := w.fmt.Sprintf(level, args...)
	if err != nil {
		return err
	}
	log.Print(string(logContent))
	return nil
}

func (w *StdoutWriter) WriteBySkipCall(level logger.Level, skipCall int, args ...interface{}) error {
	if !w.IsLoggable(level) {
		return nil
	}

	fm := w.fmt
	if w.fmt.GetSkipCall() != skipCall {
		fm = w.fmt.Copy()
		fm.SetSkipCall(skipCall)
	}

	logContent, err := fm.Sprintf(level, args...)
	if err != nil {
		return err
	}

	log.Print(string(logContent))

	return nil
}

func (w *StdoutWriter) WriteMsg(msg *logger.Msg) error {
	if !w.IsLoggable(msg.Level) {
		return nil
	}

	log.Print(string(msg.Formatted))

	return nil
}

func (w *StdoutWriter) GetMsg(level logger.Level, args ...interface{}) (*logger.Msg, error) {
	formatted, err := w.fmt.Sprintf(level, args...)
	if err != nil {
		return nil, err
	}

	return &logger.Msg{
		Level:     level,
		SkipCall:  w.fmt.GetSkipCall(),
		Formatted: formatted,
	}, nil
}

func (w *StdoutWriter) GetMsgBySkipCall(level logger.Level, skipCall int, args ...interface{}) (*logger.Msg, error) {
	fm := w.fmt
	if w.fmt.GetSkipCall() != skipCall {
		fm = w.fmt.Copy()
		fm.SetSkipCall(skipCall)
	}

	formatted, err := fm.Sprintf(level, args...)
	if err != nil {
		return nil, err
	}

	return &logger.Msg{
		Level:     level,
		SkipCall:  skipCall,
		Formatted: formatted,
	}, nil
}

func (w *StdoutWriter) GetSkipCall() int {
	return w.fmt.GetSkipCall()
}

func (w *StdoutWriter) Flush() error {
	return nil
}
