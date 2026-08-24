package application

import "errors"

func CloseAdmission(primary error, closers ...func() error) error {
	errList := []error{primary}
	for _, closeFn := range closers {
		if closeFn != nil {
			errList = append(errList, closeFn())
		}
	}
	return errors.Join(errList...)
}
