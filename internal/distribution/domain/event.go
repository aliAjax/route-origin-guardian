package domain

import "github.com/routeorigin/route-origin-guardian/internal/rib/domain"

type Delivery struct {
	Event    domain.Event
	Attempts int
	Status   string
}
