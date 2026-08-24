package application

func AdmissionStateAfterRetry(succeeded bool) AdmissionState {
	if succeeded {
		return AdmissionSucceeded
	}
	return AdmissionRetrying
}
