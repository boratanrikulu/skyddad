package crypto

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestEncodeSteganography(t *testing.T) {
	EncodeSteganography("A secret message.", getTestingPlainImage())
}

func TestDecodeSteganography(t *testing.T) {
	want := "A secret message."
	encodedImage := EncodeSteganography(want, getTestingPlainImage())
	got := DecodeSteganography(encodedImage)

	if got != want {
		t.Errorf("Input and output are not same.")
	}
}

// Private methods

func getTestingPlainImage() []byte {
	b, err := ioutil.ReadFile("steganography_test.png")
	if err != nil {
		log.Fatalf("Error occur while reading test image: %v", err)
	}

	return b
}
