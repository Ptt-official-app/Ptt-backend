package main

import (
	"log/slog"
	"testing"
)

func TestSlogLevelForLogLevel(t *testing.T) {
	tests := []struct {
		logLevel uint64
		want     slog.Level
	}{
		{0, slog.LevelError},
		{3, slog.LevelError},
		{4, slog.LevelWarn},
		{5, slog.LevelWarn},
		{6, slog.LevelInfo},
		{7, slog.LevelDebug},
		{99, slog.LevelDebug},
	}

	for _, test := range tests {
		if got := slogLevelForLogLevel(test.logLevel); got != test.want {
			t.Errorf("slogLevelForLogLevel(%d) = %v, want %v", test.logLevel, got, test.want)
		}
	}
}
