package application

import "testing"

func TestAdmissionTerminalStates(t *testing.T) {
	if !AdmissionSucceeded.Terminal() || !AdmissionFailed.Terminal() || AdmissionRetrying.Terminal() {
		t.Fatal("terminal admission states are inconsistent")
	}
}

func TestAdmissionRetryTransitionAllowed(t *testing.T) {
	if !CanTransitionAdmission(AdmissionRetrying, AdmissionSucceeded) {
		t.Fatal("successful retry cannot reach succeeded")
	}
}

func TestAdmissionRetryWritesSucceeded(t *testing.T) {
	if got := AdmissionStateAfterRetry(true); got != AdmissionSucceeded {
		t.Fatalf("retry state = %s", got)
	}
}

func TestAdmissionRetryVisibleInFlight(t *testing.T) {
	if !AdmissionInFlight(AdmissionRetrying) {
		t.Fatal("retrying admission disappeared from in-flight view")
	}
}
