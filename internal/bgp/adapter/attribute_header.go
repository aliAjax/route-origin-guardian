package adapter

import "fmt"

type attributeHeader struct {
	flags byte
	code  byte
}

func parseAttributeHeader(payload []byte, cursor int) (attributeHeader, error) {
	if len(payload) < 3 {
		return attributeHeader{}, fmt.Errorf("attribute header truncated at offset %d: need 3 bytes, have %d", cursor, len(payload))
	}
	return attributeHeader{flags: payload[0], code: payload[1]}, nil
}
