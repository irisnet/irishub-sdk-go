package log

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger("info")
}

type Logger struct {
	zerolog.Logger
}

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

	log := zerolog.New(prettyWriter(os.Stdout)).
		Level(l).
		With().
		Caller().
		Timestamp().Logger()

	logger := &Logger{log}
	defaultLogger = logger
	return logger
}

func (l *Logger) SetOutput(w io.Writer) {
	l.Logger = zerolog.New(prettyWriter(w)).
		Level(l.GetLevel()).
		With().
		Timestamp().Logger()
}

func GetLogger() *Logger {
	return defaultLogger
}

func prettyWriter(w io.Writer) io.Writer {
	output := zerolog.ConsoleWriter{Out: w, TimeFormat: time.RFC3339}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("****%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("[%s]:", i)
	}
	return output
}
