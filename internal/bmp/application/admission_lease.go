package application

import "sync"

type AdmissionLease struct {
	once    sync.Once
	release func() error
	err     error
}

func NewAdmissionLease(release func() error) *AdmissionLease {
	return &AdmissionLease{release: release}
}

func (l *AdmissionLease) Release() error {
	if l == nil {
		return nil
	}
	l.once.Do(func() {
		if l.release != nil {
			l.err = l.release()
		}
	})
	return l.err
}
