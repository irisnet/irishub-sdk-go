package log

import (
	"fmt"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	log1 := NewLogger("info").With("test")

	w, err := os.Create("1.txt")
	if err != nil {
		fmt.Println(err.Error())
	}

	log1.SetOutput(w)

	log1.Info().Str("foo", "bar").Msg("Hello World")
	log2 := log1.With("test1")
	log2.Info().Str("foo1", "bar").Msg("Hello World")
	log1.Info().Str("foo2", "bar").Msg("Hello World")
	log2.Info().Str("foo3", "bar").Msg("Hello World")

}
