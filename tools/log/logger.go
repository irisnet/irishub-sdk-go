package log

import (
	"github.com/rs/zerolog"
	"os"
)

type Logger struct {
	zerolog.Logger
}

//type CallerHook struct{}
//
//func (hook CallerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
//	pc := make([]uintptr, 10)
//	runtime.Callers(5, pc)
//	f := runtime.FuncForPC(pc[0])
//	e.Str("method", f.Name())
//}

func NewLogger() *Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05.000"}
	log := zerolog.New(output).With().
		Caller().
		Timestamp().Logger()
	//Hook(CallerHook{})
	return &Logger{
		log,
	}
}
