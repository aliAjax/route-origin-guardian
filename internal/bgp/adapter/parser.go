package adapter

import (
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
	cursor := 0
	for len(payload) > 0 {
		header, err := parseAttributeHeader(payload, cursor)
		if err != nil {
			return nil, err
		}
		if header.code == 0 {
			return nil, fmt.Errorf("reserved attribute code at offset %d", cursor+1)
		}
		n, valueOffset, err := parseAttributeLength(payload, header, cursor)
		if err != nil {
			return nil, err
		}
		value, rest, err := parseAttributeValue(payload, n, valueOffset, cursor)
		if err != nil {
			return nil, err
		}
		out = append(out, Attribute{header.flags, header.code, value})
		consumed := valueOffset + n
		cursor += consumed
		payload = rest
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
