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

func NewLogger(level string) *Logger {
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		l = zerolog.InfoLevel
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05.000"}
	log := zerolog.New(output).
		Level(l).
		With().
		Caller().
		Timestamp().Logger()
	//Hook(CallerHook{})
	return &Logger{
		log,
	}
}
