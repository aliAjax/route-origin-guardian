package domain

import (
	"fmt"
	"net/netip"
	"strings"
)

type Message struct {
	Type    uint8
	Payload []byte
}

func Encode(m Message) []byte {
	l := 19 + len(m.Payload)
	b := make([]byte, l)
	for i := 0; i < 16; i++ {
		b[i] = 0xff
	}
	b[16] = byte(l >> 8)
	b[17] = byte(l)
	b[18] = m.Type
	copy(b[19:], m.Payload)
	return b
}
func Decode(b []byte) (Message, error) {
	if len(b) < 19 {
		return Message{}, fmt.Errorf("bgp message too short")
	}
	for i := 0; i < 16; i++ {
		if b[i] != 0xff {
			return Message{}, fmt.Errorf("invalid marker")
		}
	}
	l := int(b[16])<<8 | int(b[17])
	if l < 19 || l > 4096 || l != len(b) {
		return Message{}, fmt.Errorf("invalid bgp length %d", l)
	}
	return Message{Type: b[18], Payload: append([]byte(nil), b[19:]...)}, nil
}
func PrefixBits(prefix string) (int, error) {
	p, e := netip.ParsePrefix(strings.TrimSpace(prefix))
	if e != nil {
		return 0, fmt.Errorf("invalid prefix: %w", e)
	}
	return p.Bits(), nil
}
