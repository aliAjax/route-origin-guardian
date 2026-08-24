package application

import (
	"errors"
	"fmt"
	"testing"
)

func TestAdmissionBatchReleasesEachLease(t *testing.T) {
	active := 0
	open := func(item string) (*AdmissionLease, error) {
		if active != 0 {
			return nil, fmt.Errorf("previous lease still active before %s", item)
		}
		active++
		return NewAdmissionLease(func() error { active--; return nil }), nil
	}
	if err := ProcessAdmissionBatch([]string{"a", "b", "c"}, open, func(*AdmissionLease) error { return nil }); err != nil {
		t.Fatal(err)
	}
	if active != 0 {
		t.Fatalf("active leases = %d", active)
	}
}

func TestFinishAdmissionPreservesBothErrors(t *testing.T) {
	primary := errors.New("validation failed")
	commitErr := errors.New("commit failed")
	err := FinishAdmission(primary, func() error { return commitErr })
	if !errors.Is(err, primary) || !errors.Is(err, commitErr) {
		t.Fatalf("combined error = %v", err)
	}
}

func TestAbortAdmissionAlwaysReleases(t *testing.T) {
	primary := errors.New("admission rejected")
	released := false
	err := AbortAdmission(primary, func() error { released = true; return nil })
	if !released || !errors.Is(err, primary) {
		t.Fatalf("released=%v err=%v", released, err)
	}
}

func TestCloseAdmissionPreservesEveryError(t *testing.T) {
	primary := errors.New("read failed")
	closeOne := errors.New("peer close failed")
	closeTwo := errors.New("budget close failed")
	err := CloseAdmission(primary, func() error { return closeOne }, func() error { return closeTwo })
	for _, want := range []error{primary, closeOne, closeTwo} {
		if !errors.Is(err, want) {
			t.Fatalf("%v missing from %v", want, err)
		}
	}
}
