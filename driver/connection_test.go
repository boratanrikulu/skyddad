package driver

import (
	"testing"
)

// Check if the connect is working.
func TestConnect(t *testing.T) {
	got := Connect()
	if got == nil {
		t.Errorf("Error.")
	}
}
