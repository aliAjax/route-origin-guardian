package logging

import "log/slog"

// DefaultLevel is the level used when deployment configuration omits or
// cannot be mapped to a known level.
func DefaultLevel() slog.Level { return slog.LevelInfo }
