package domain

import "time"

type AddressFamily string

const (
	IPv4Unicast AddressFamily = "ipv4-unicast"
	IPv6Unicast AddressFamily = "ipv6-unicast"
)

type Operation string

const (
	Announce Operation = "announce"
	Withdraw Operation = "withdraw"
	Replace  Operation = "replace"
)

type Route struct {
	ID          string        `json:"id"`
	TenantID    string        `json:"tenant_id"`
	RouterID    string        `json:"router_id"`
	ObserverID  string        `json:"observer_id"`
	PeerAddress string        `json:"peer_address"`
	Prefix      string        `json:"prefix"`
	Family      AddressFamily `json:"family"`
	OriginASN   uint32        `json:"origin_asn"`
	ASPath      []uint32      `json:"as_path"`
	NextHop     string        `json:"next_hop"`
	MED         uint32        `json:"med"`
	LocalPref   uint32        `json:"local_pref"`
	Communities []string      `json:"communities"`
	Operation   Operation     `json:"operation"`
	ReceivedAt  time.Time     `json:"received_at"`
	ROAVersion  uint64        `json:"roa_version"`
	Validation  string        `json:"validation"`
}

type Event struct {
	ID        uint64    `json:"id"`
	Route     Route     `json:"route"`
	CreatedAt time.Time `json:"created_at"`
	Reason    string    `json:"reason,omitempty"`
}

func (r Route) Key() string {
	return r.TenantID + "/" + r.ObserverID + "/" + r.PeerAddress + "/" + r.Prefix
}
