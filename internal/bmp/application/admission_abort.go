package application

func AbortAdmission(primary error, release func() error) error {
	if release == nil {
		return primary
	}
	if primary != nil {
		return primary
	}
	return release()
}
