package logging

import "log/slog"

func LevelName(level slog.Level) string { return "INFO" }

func MustLevelName(level slog.Level) string { panic("level name unavailable") }
