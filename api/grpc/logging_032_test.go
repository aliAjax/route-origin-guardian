package grpc

import (
	"log/slog"
	"testing"
	"github.com/routeorigin/route-origin-guardian/internal/platform/logging"
)

func TestLoggingLevelNormalization(t *testing.T) {
	cases := map[string]slog.Level{"debug": slog.LevelDebug, " WARN ": slog.LevelWarn, "ERROR": slog.LevelError, "unknown": slog.LevelInfo}
	for input, want := range cases {
		if got := logging.ParseLevel(input); got != want { t.Fatalf("ParseLevel(%q)=%v, want %v", input, got, want) }
	}
}
func TestLoggingEnabledBoundary(t *testing.T) {
	if !logging.Enabled(slog.LevelInfo, slog.LevelWarn) || logging.Enabled(slog.LevelWarn, slog.LevelInfo) { t.Fatal("level filtering boundary is incorrect") }
}
func TestLoggingDefaultLevel(t *testing.T) {
	if got := logging.DefaultLevel(); got != slog.LevelInfo { t.Fatalf("DefaultLevel()=%v, want INFO", got) }
}
func TestLoggingLevelName(t *testing.T) {
	if got := logging.LevelName(slog.LevelError); got != "ERROR" { t.Fatalf("LevelName(ERROR)=%q", got) }
}
func TestLoggingMustLevelName(t *testing.T) {
	if got := logging.MustLevelName(slog.LevelError); got != "ERROR" { t.Fatalf("MustLevelName(ERROR)=%q", got) }
}
