package application

func CloseAdmission(primary error, closers ...func() error) error {
	result := primary
	for _, closeFn := range closers {
		if closeFn != nil {
			if err := closeFn(); err != nil {
				result = err
			}
		}
	}
	return result
}
