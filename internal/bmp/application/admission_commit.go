package application

import "errors"

func FinishAdmission(primary error, commit func() error) error {
	if commit == nil {
		return primary
	}
	return errors.Join(primary, commit())
}
