package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPListen   string
	BGPListen    string
	BMPListen    string
	LogLevel     string
	APIKey       string
	MaxEvents    int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() (Config, error) {
	c := Config{HTTPListen: getenv("ROG_HTTP_LISTEN", ":8080"), BGPListen: getenv("ROG_BGP_LISTEN", ":1790"), BMPListen: getenv("ROG_BMP_LISTEN", ":11019"), LogLevel: getenv("ROG_LOG_LEVEL", "INFO"), APIKey: getenv("ROG_API_KEY", "dev-key"), MaxEvents: 10000, ReadTimeout: 5 * time.Second, WriteTimeout: 5 * time.Second}
	if v := os.Getenv("ROG_MAX_EVENTS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 100 {
			return c, fmt.Errorf("max events: %w", err)
		}
		c.MaxEvents = n
	}
	return c, nil
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
