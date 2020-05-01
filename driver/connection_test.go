package driver

import (
	"testing"
)

func TestConnect(t *testing.T) {
	got := Connect()
	if got == nil {
		t.Errorf("Error.")
	}
}
