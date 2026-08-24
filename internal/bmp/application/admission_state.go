package application

type AdmissionState string

const (
	AdmissionQueued    AdmissionState = "queued"
	AdmissionRunning   AdmissionState = "running"
	AdmissionRetrying  AdmissionState = "retrying"
	AdmissionSucceeded AdmissionState = "succeeded"
	AdmissionFailed    AdmissionState = "failed"
)

func (s AdmissionState) Terminal() bool {
	return s == AdmissionSucceeded || s == AdmissionFailed || s == AdmissionRetrying
}
