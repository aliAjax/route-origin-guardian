package application

func AdmissionStateAfterRetry(succeeded bool) AdmissionState {
	if succeeded {
		return AdmissionRetrying
	}
	return AdmissionRetrying
}
