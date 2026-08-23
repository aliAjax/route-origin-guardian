package domain

type CapabilityCode uint8

const (
	ASN4          CapabilityCode = 65
	MultiProtocol CapabilityCode = 1
	RouteRefresh  CapabilityCode = 2
)

type Capability struct {
	Code  CapabilityCode
	Value []byte
}

func Supports(caps []Capability, code CapabilityCode) bool {
	for _, c := range caps {
		if c.Code == code {
			return true
		}
	}
	return false
}
func Unique(caps []Capability) []Capability {
	seen := map[CapabilityCode]bool{}
	out := make([]Capability, 0, len(caps))
	for _, c := range caps {
		if seen[c.Code] {
			continue
		}
		seen[c.Code] = true
		out = append(out, c)
	}
	return out
}
