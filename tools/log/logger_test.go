package log

import "testing"

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	logger.Info().Str("foo", "bar").Msg("Hello World")
}
