package application

func AdmissionInFlight(state AdmissionState) bool {
	return state == AdmissionQueued || state == AdmissionRunning
}
