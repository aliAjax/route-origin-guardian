package domain

import (
	"fmt"
	"net/netip"
	"sort"
	"sync"
)

type Status string

const (
	Valid    Status = "valid"
	Invalid  Status = "invalid"
	NotFound Status = "not_found"
)

type ROA struct {
	Prefix    string `json:"prefix"`
	ASN       uint32 `json:"asn"`
	MaxLength int    `json:"max_length"`
}
type Result struct {
	Status    Status `json:"status"`
	Prefix    string `json:"prefix"`
	OriginASN uint32 `json:"origin_asn"`
	Matched   *ROA   `json:"matched_roa,omitempty"`
	Reason    string `json:"reason"`
	Version   uint64 `json:"roa_version"`
}
type Index struct {
	mu      sync.RWMutex
	version uint64
	roas    []ROA
}

func NewIndex() *Index { return &Index{} }
func (i *Index) Replace(roas []ROA) error {
	for _, r := range roas {
		p, e := netip.ParsePrefix(r.Prefix)
		if e != nil {
			return fmt.Errorf("roa prefix %s: %w", r.Prefix, e)
		}
		if r.MaxLength < p.Bits() {
			return fmt.Errorf("max length for %s", r.Prefix)
		}
	}
	cp := append([]ROA(nil), roas...)
	sort.Slice(cp, func(a, b int) bool { return len(cp[a].Prefix) > len(cp[b].Prefix) })
	i.mu.Lock()
	i.version++
	i.roas = cp
	i.mu.Unlock()
	return nil
}
func (i *Index) Version() uint64 { i.mu.RLock(); defer i.mu.RUnlock(); return i.version }
func (i *Index) Validate(prefix string, asn uint32) Result {
	i.mu.RLock()
	defer i.mu.RUnlock()
	res := Result{Status: NotFound, Prefix: prefix, OriginASN: asn, Version: i.version, Reason: "no covering ROA"}
	p, e := netip.ParsePrefix(prefix)
	if e != nil {
		res.Status = Invalid
		res.Reason = "invalid prefix"
		return res
	}
	for _, r := range i.roas {
		rp, e := netip.ParsePrefix(r.Prefix)
		if e != nil || !rp.Contains(p.Addr()) {
			continue
		}
		if p.Bits() > r.MaxLength {
			res.Status = Invalid
			res.Matched = &r
			res.Reason = "prefix length exceeds ROA max length"
			return res
		}
		if r.ASN == asn {
			res.Status = Valid
			res.Matched = &r
			res.Reason = "origin authorized"
			return res
		}
		res.Status = Invalid
		res.Matched = &r
		res.Reason = "origin ASN mismatch"
		return res
	}
	return res
}
