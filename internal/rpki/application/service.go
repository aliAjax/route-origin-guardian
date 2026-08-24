package application

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/rpki/domain"
)

type Snapshot struct {
	Serial uint64       `json:"serial"`
	ROAs   []domain.ROA `json:"roas"`
	Digest string       `json:"digest,omitempty"`
}
type Service struct {
	index *domain.Index
	last  Snapshot
}

func NewService(i *domain.Index) *Service { return &Service{index: i} }
func (s *Service) Import(_ context.Context, snap Snapshot) error {
	b, e := json.Marshal(struct {
		Serial uint64
		ROAs   []domain.ROA
	}{snap.Serial, snap.ROAs})
	if e != nil {
		return fmt.Errorf("marshal snapshot: %w", e)
	}
	d := sha256.Sum256(b)
	if snap.Digest != "" && snap.Digest != fmt.Sprintf("%x", d) {
		return fmt.Errorf("snapshot digest mismatch")
	}
	if e = s.index.Replace(snap.ROAs); e != nil {
		return fmt.Errorf("replace ROA index: %w", e)
	}
	s.last = snap
	return nil
}
func (s *Service) Validate(prefix string, asn uint32) domain.Result {
	return s.index.Validate(prefix, asn)
}
func (s *Service) Snapshot() Snapshot { return s.last }
