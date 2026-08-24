package adapter

import "fmt"

func parseAttributeValue(payload []byte, length int, valueOffset int, cursor int) ([]byte, []byte, error) {
	available := len(payload) - valueOffset
	if length > 4096 || available < length {
		return nil, nil, fmt.Errorf("attribute value truncated at offset %d: need %d bytes, have %d", cursor+valueOffset, length, available)
	}
	value := append([]byte(nil), payload[valueOffset:valueOffset+length]...)
	return value, payload[valueOffset+length:], nil
}
