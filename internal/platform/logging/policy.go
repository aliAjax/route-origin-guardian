package logging

import "log/slog"

func Enabled(minimum, record slog.Level) bool {
	return record <= minimum
}
