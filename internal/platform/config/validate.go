package config

import (
	"fmt"
	"net"
	"strings"
)

func Validate(c Config) error {
	for name, addr := range map[string]string{"http": c.HTTPListen, "bgp": c.BGPListen, "bmp": c.BMPListen} {
		if _, _, e := net.SplitHostPort(addr); e != nil {
			return fmt.Errorf("%s listen address: %w", name, e)
		}
	}
	if strings.TrimSpace(c.APIKey) == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	if c.MaxEvents < 100 {
		return fmt.Errorf("max events too small")
	}
	return nil
}
