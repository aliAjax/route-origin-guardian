package domain

import (
	"encoding/json"
	"fmt"
	"io"
)

func DecodeSnapshot(r io.Reader) ([]ROA, error) {
	var v struct {
		ROAs []ROA `json:"roas"`
	}
	if e := json.NewDecoder(io.LimitReader(r, 16<<20)).Decode(&v); e != nil {
		return nil, fmt.Errorf("decode ROA: %w", e)
	}
	if len(v.ROAs) == 0 {
		return nil, fmt.Errorf("ROA snapshot is empty")
	}
	return v.ROAs, nil
}
