package httpx

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
)

func RequestID(r *http.Request) string {
	if v := strings.TrimSpace(r.Header.Get("X-Request-ID")); v != "" {
		return v
	}
	b := make([]byte, 8)
	if _, e := rand.Read(b); e != nil {
		return "generated"
	}
	return hex.EncodeToString(b)
}
func ClientIP(r *http.Request) string {
	if v := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]); v != "" {
		return v
	}
	return r.RemoteAddr
}
func IsJSON(r *http.Request) bool {
	return strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "application/json")
}
