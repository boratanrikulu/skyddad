package skyddad

import (
  "testing"
)

func TestHello(t *testing.T) {
  got := Hello()
  if want := "Hello"; got != want {
    t.Errorf("Error.")
  }
}
