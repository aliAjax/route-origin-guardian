package application

func CanTransitionAdmission(from, to AdmissionState) bool {
	switch from {
	case AdmissionQueued:
		return to == AdmissionRunning
	case AdmissionRunning:
		return to == AdmissionRetrying || to == AdmissionSucceeded || to == AdmissionFailed
	case AdmissionRetrying:
		return to == AdmissionRunning || to == AdmissionFailed
	default:
		return false
	}
}
