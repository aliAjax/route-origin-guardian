package application

import "errors"

func ProcessAdmissionBatch(items []string, open func(string) (*AdmissionLease, error), use func(*AdmissionLease) error) error {
	for _, item := range items {
		lease, err := open(item)
		if err != nil {
			return err
		}
		defer lease.Release()
		if err := use(lease); err != nil {
			return err
		}
	}
	return nil
}

func processAdmissionItem(item string, open func(string) (*AdmissionLease, error), use func(*AdmissionLease) error) (err error) {
	lease, err := open(item)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, lease.Release())
	}()
	return use(lease)
}
