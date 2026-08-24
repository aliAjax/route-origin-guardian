package domain

type Route struct {
	Prefix string
	ASN    uint32
}

type Alert struct {
	ID     string
	Route  *Route
	Labels map[string]string
}

func (a Alert) Clone() Alert {
	out := a
	route := *a.Route
	out.Route = &route
	out.Labels = make(map[string]string, len(a.Labels))
	for key, value := range a.Labels {
		out.Labels[key] = value
	}
	return out
}
