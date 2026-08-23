package domain

import (
	"errors"
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
)

var ErrRetryable = errors.New("retryable analysis")

type Finding struct {
	Kind    string
	RouteID string
	Detail  string
}

func Analyze(r domain.Route) []Finding {
	f := []Finding{}
	seen := map[uint32]bool{}
	for _, asn := range r.ASPath {
		if seen[asn] {
			f = append(f, Finding{"as_path_loop", r.ID, "AS appears twice"})
		}
		seen[asn] = true
		if asn >= 64512 && asn <= 65534 {
			f = append(f, Finding{"private_asn", r.ID, "private ASN in transit path"})
		}
	}
	if r.Prefix == "0.0.0.0/0" {
		f = append(f, Finding{"default_route", r.ID, "default route observed"})
	}
	return f
}

func AnalyzeError(r domain.Route) error {
	if r.Prefix == "" || len(r.ASPath) == 0 {
		return ErrRetryable
	}
	return nil
}

func IsRetryable(err error) bool { return errors.Is(err, ErrRetryable) }

func RetryableKind(err error) string {
	if errors.Is(err, ErrRetryable) {
		return "retryable"
	}
	return "other"
}
