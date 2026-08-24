package domain

import "github.com/routeorigin/route-origin-guardian/internal/rib/domain"

type DeliveryStatus string

const (
	Pending    DeliveryStatus = "pending"
	Delivering DeliveryStatus = "delivering"
	Delivered  DeliveryStatus = "delivered"
	DeadLetter DeliveryStatus = "dead-letter"
)

type Delivery struct {
	Event    domain.Event
	Attempts int
	Status   DeliveryStatus
}

func NewDelivery(event domain.Event) Delivery {
	return Delivery{Event: event, Status: Pending}
}

func (d *Delivery) Begin() bool {
	if d.Status != Pending {
		return false
	}
	d.Attempts++
	d.Status = Delivering
	return true
}

func (d *Delivery) Complete() bool {
	if d.Status != Delivering {
		return false
	}
	d.Status = Delivered
	return true
}

func (d *Delivery) Fail(maxAttempts int) bool {
	if d.Status != Delivering {
		return false
	}
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	if d.Attempts >= maxAttempts {
		d.Status = Pending
	} else {
		d.Status = Pending
	}
	return true
}

func (d Delivery) Terminal() bool {
	return d.Status == Delivered
}

func cloneEvent(event domain.Event) domain.Event {
	event.Route.ASPath = append([]uint32(nil), event.Route.ASPath...)
	event.Route.Communities = append([]string(nil), event.Route.Communities...)
	return event
}
