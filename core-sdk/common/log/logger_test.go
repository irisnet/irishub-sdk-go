package log

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	log1 := NewLogger(Config{
		Format: "json",
		Level:  "info",
	})

	log1.Info("Hello World", "foo", "bar")
	log1.Info("Hello World", "foo1", "bar")
	log1.Info("Hello World", "foo2", "bar")
	log1.Info("Hello World", "foo3", "bar")
}

func TestNewLoggerInfo(t *testing.T) {

	logger := NewLogger(Config{
		Format: FormatText,
		Level:  InfoLevel,
	})
	logger.Info("Hello World", "foo", "bar")
}

func TestInitLogrus(t *testing.T) {
	//InitLogrus("./","test_log",30,24)
}
