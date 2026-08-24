package application

func FinishAdmission(primary error, commit func() error) error {
	if commit == nil {
		return primary
	}
	return commit()
}
