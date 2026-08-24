package adapter

import (
	"encoding/binary"
	"fmt"
)

func parseAttributeLength(payload []byte, header attributeHeader, cursor int) (length int, valueOffset int, err error) {
	if header.flags&0x10 == 0 {
		return int(payload[2]), 3, nil
	}
	if len(payload) < 4 {
		return 0, 0, fmt.Errorf("extended attribute length truncated at offset %d: need 4 bytes, have %d", cursor, len(payload))
	}
	return int(binary.BigEndian.Uint16(payload[2:4])), 4, nil
}
