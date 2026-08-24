package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
)

type JSONMessage struct {
	Operation string            `json:"operation"`
	Route     domain.Route      `json:"route"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

func DecodeJSON(line []byte) (JSONMessage, error) {
	var m JSONMessage
	if len(line) > 1<<20 {
		return m, fmt.Errorf("JSON message exceeds 1MiB")
	}
	if e := json.Unmarshal(line, &m); e != nil {
		return m, fmt.Errorf("decode collector message: %w", e)
	}
	if m.Operation == "" {
		m.Operation = "announce"
	}
	return m, nil
}
func EncodeJSON(m JSONMessage) ([]byte, error) {
	b, e := json.Marshal(m)
	if e != nil {
		return nil, fmt.Errorf("encode collector message: %w", e)
	}
	return append(b, '\n'), nil
}
