package logging

import "log/slog"

// Enabled reports whether a record at the given level passes a minimum
// threshold. slog levels grow with severity (Debug < Info < Warn < Error),
// so a record is enabled when it is at least as severe as the minimum.
func Enabled(minimum, record slog.Level) bool {
	return record >= minimum
}
