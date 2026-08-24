package adapter

import (
	"encoding/binary"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/bgp/domain"
)

type Attribute struct {
	Flags byte
	Code  byte
	Value []byte
}

func ParseAttributes(payload []byte) ([]Attribute, error) {
	var out []Attribute
	for len(payload) > 0 {
		if len(payload) < 3 {
			return nil, fmt.Errorf("attribute header truncated")
		}
		flags, code := payload[0], payload[1]
		n := int(payload[2])
		off := 3
		if flags&0x10 != 0 {
			if len(payload) < 4 {
				return nil, fmt.Errorf("extended attribute length truncated")
			}
			n = int(binary.BigEndian.Uint16(payload[2:4]))
			off = 4
		}
		if n > 4096 || len(payload) < off+n {
			return nil, fmt.Errorf("attribute length %d", n)
		}
		out = append(out, Attribute{flags, code, append([]byte(nil), payload[off:off+n]...)})
		payload = payload[off+n:]
	}
	return out, nil
}
func ValidatePrefix(prefix string) error {
	bits, e := domain.PrefixBits(prefix)
	if e != nil {
		return e
	}
	if bits == 0 && prefix != "0.0.0.0/0" {
		return fmt.Errorf("noncanonical default prefix")
	}
	return nil
}
