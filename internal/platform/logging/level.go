package logging

import "log/slog"

// ParseLevel translates deployment configuration into slog levels.
func ParseLevel(level string) slog.Level {
	switch normalizeLevel(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelError
	case "ERROR":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
