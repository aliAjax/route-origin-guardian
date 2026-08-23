package logging

import (
	"log/slog"
	"os"
)

func New(level string) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: ParseLevel(level)}))
}
