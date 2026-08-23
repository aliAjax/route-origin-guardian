package adapter

import (
	"encoding/binary"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/bgp/domain"
)

func EncodeOpen(asn uint32, hold uint16, bgpID [4]byte) domain.Message {
	p := make([]byte, 29)
	p[0] = 4
	binary.BigEndian.PutUint16(p[1:3], uint16(asn))
	binary.BigEndian.PutUint16(p[3:5], hold)
	copy(p[5:9], bgpID[:])
	return domain.Message{Type: 1, Payload: p}
}
func EncodeNotification(code, subcode byte, data []byte) domain.Message {
	if len(data) > 4000 {
		data = data[:4000]
	}
	p := append([]byte{code, subcode}, data...)
	return domain.Message{Type: 3, Payload: p}
}
func ParseOpen(m domain.Message) (uint32, timeFields, error) {
	if m.Type != 1 || len(m.Payload) < 9 {
		return 0, timeFields{}, fmt.Errorf("invalid open")
	}
	return uint32(binary.BigEndian.Uint16(m.Payload[1:3])), timeFields{Hold: binary.BigEndian.Uint16(m.Payload[3:5])}, nil
}

type timeFields struct{ Hold uint16 }
