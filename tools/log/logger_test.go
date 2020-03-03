package log

import "testing"

func TestNewLogger(t *testing.T) {
	log1 := NewLogger("info").With("test")
	log1.Info().Str("foo", "bar").Msg("Hello World")
	log2 := log1.With("test1")
	log2.Info().Str("foo", "bar").Msg("Hello World")
	log1.Info().Str("foo", "bar").Msg("Hello World")
	log2.Info().Str("foo", "bar").Msg("Hello World")
}
