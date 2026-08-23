package logging

import "log/slog"

// levelName returns the canonical deployment name for the four well-known
// slog levels. Unknown or custom levels map to the empty string, signalling
// the absence of a stable name.
func levelName(level slog.Level) (string, bool) {
	switch level {
	case slog.LevelDebug:
		return "DEBUG", true
	case slog.LevelInfo:
		return "INFO", true
	case slog.LevelWarn:
		return "WARN", true
	case slog.LevelError:
		return "ERROR", true
	default:
		return "", false
	}
}

// LevelName returns the canonical name for a level, falling back to "INFO"
// for unknown or custom levels so it never panics.
func LevelName(level slog.Level) string {
	if name, ok := levelName(level); ok {
		return name
	}
	return "INFO"
}

// MustLevelName returns the canonical name for a known level and panics when
// the level has no stable mapping, surfacing misconfiguration early.
func MustLevelName(level slog.Level) string {
	name, ok := levelName(level)
	if !ok {
		panic("level name unavailable")
	}
	return name
}
