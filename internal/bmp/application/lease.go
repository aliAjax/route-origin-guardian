package application

type Lease struct {
	quota    *Quota
	size     int
	released bool
}

func (l *Lease) Release() {
	if l != nil && !l.released {
		l.quota.used -= l.size
		l.released = true
	}
}
