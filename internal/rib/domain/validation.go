package domain

import (
	"fmt"
	"net/netip"
	"strings"
)

func ValidateFamily(f AddressFamily) error {
	switch f {
	case IPv4Unicast, IPv6Unicast:
		return nil
	default:
		return fmt.Errorf("unsupported address family %q", f)
	}
}
func ValidatePeer(peer string) error {
	a, e := netip.ParseAddr(peer)
	if e != nil {
		return fmt.Errorf("peer address: %w", e)
	}
	if a.IsUnspecified() || a.IsMulticast() {
		return fmt.Errorf("peer address is not unicast")
	}
	return nil
}
func CanonicalPrefix(prefix string) (string, error) {
	p, e := netip.ParsePrefix(strings.TrimSpace(prefix))
	if e != nil {
		return "", fmt.Errorf("prefix: %w", e)
	}
	return p.Masked().String(), nil
}
func ValidateASPath(path []uint32) error {
	if len(path) > 255 {
		return fmt.Errorf("AS path too long")
	}
	for _, asn := range path {
		if asn == 0 {
			return fmt.Errorf("AS path contains zero ASN")
		}
	}
	return nil
}
func ValidateRoute(r Route) error {
	if ValidateFamily(r.Family) != nil {
		return ValidateFamily(r.Family)
	}
	if ValidatePeer(r.PeerAddress) != nil {
		return ValidatePeer(r.PeerAddress)
	}
	if _, e := CanonicalPrefix(r.Prefix); e != nil {
		return e
	}
	return ValidateASPath(r.ASPath)
}
