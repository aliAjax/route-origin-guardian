package logging

import "strings"

func normalizeLevel(level string) string {
	return strings.ToUpper(strings.TrimSpace(level))
}
