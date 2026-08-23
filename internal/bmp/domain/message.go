package domain

type Type uint8

const (
	Initiation      Type = 1
	PeerUp          Type = 3
	PeerDown        Type = 4
	RouteMonitoring Type = 0
	Statistics      Type = 1
)

type Message struct {
	Type    Type
	Router  string
	Peer    string
	Payload []byte
}
