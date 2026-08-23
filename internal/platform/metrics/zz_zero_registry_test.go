package metrics

import "testing"

func TestZeroRegistryCounterDoesNotPanic(t *testing.T) {
	var registry Registry
	registry.Counter("embedded_requests").Inc()
	if got := registry.Counter("embedded_requests").Get(); got != 1 {
		t.Fatalf("counter value = %d", got)
	}
	if text := registry.Text(); text != "embedded_requests 1\n" {
		t.Fatalf("counter text = %q", text)
	}
}
