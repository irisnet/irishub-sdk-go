package log

import (
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
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

	//override default CallerMarshalFunc
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		index := strings.LastIndex(file, string(os.PathSeparator))
		rs := []rune(file)
		fileName := string(rs[index+1:])
		return fileName + ":" + strconv.Itoa(line)
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

func (l *Logger) With(component string) *Logger {
	logger := *l
	logger.Logger = logger.Logger.With().Str("component", component).Logger()
	return &logger
}
