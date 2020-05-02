package crypto

import (
	"bytes"
	"image"
	"log"

	"github.com/auyer/steganography"
)

// EncodeSteganography encodes the image with given secret message.
// It returns the encoded image file as a byte array.
func EncodeSteganography(message string, plainImage []byte) []byte {
	img, _, err := image.Decode(bytes.NewReader(plainImage))
	if err != nil {
		log.Fatalf("Error opening file %v", err)
	}

	encodedImg := new(bytes.Buffer)
	err = steganography.Encode(encodedImg, img, []byte(message)) // Calls library and Encodes the message into a new buffer
	if err != nil {
		log.Fatalf("Error encoding message into file  %v", err)
	}

	return encodedImg.Bytes()
}

// DecodeSteganography decodes the image.
// Get the secret message and returns it as a string.
func DecodeSteganography(secretImage []byte) string {
	if len(secretImage) == 0 {
		log.Fatal("Secret image can not be empty.")
	}

	img, _, err := image.Decode(bytes.NewReader(secretImage))
	if err != nil {
		log.Fatal("error decoding file", img)
	}

	sizeOfMessage := steganography.GetMessageSizeFromImage(img) // Uses the library to check the message size

	msg := steganography.Decode(sizeOfMessage, img) // Read the message from the picture file

	return string(msg)
}
