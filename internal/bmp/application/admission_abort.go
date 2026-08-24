package application

import "errors"

func AbortAdmission(primary error, release func() error) error {
	if release == nil {
		return primary
	}
	return errors.Join(primary, release())
}
