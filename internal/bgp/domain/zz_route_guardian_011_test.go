package domain

import "testing"

func TestCapabilityPrefixValidation(t *testing.T) {
	if bits, err := PrefixBits("999.1.1.1/24"); err == nil || bits != 0 {
		t.Fatalf("invalid IPv4 prefix accepted: bits=%d err=%v", bits, err)
	}
	if bits, err := PrefixBits("192.0.2.0/24"); err != nil || bits != 24 {
		t.Fatalf("valid prefix rejected: bits=%d err=%v", bits, err)
	}
	value := []byte{1, 2}
	unique := Unique([]Capability{{Code: ASN4, Value: value}})
	value[0] = 9
	if unique[0].Value[0] != 1 {
		t.Fatalf("capability value aliases input: %v", unique[0].Value)
	}
}
