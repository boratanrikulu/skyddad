package main

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "encoding/hex"
  "log"
  "io"
)

func encrypt(body string, keya string) string {
  // TO DO: change the key value.
  key, _ := hex.DecodeString(keya)

  // TO DO: change the message value.
  plaintext := []byte(body)

  // Create aes cipher by using key.
  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err)
  }


	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

  // Creates stream value and XORs it.
  stream := cipher.NewCFBEncrypter(block, iv)
  stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

  return hex.EncodeToString(ciphertext)
}

func decrypt(encryptedBody, keya string) string {
  // Load your secret key from a safe place and reuse it across multiple
  // NewCipher calls. (Obviously don't use this example key for anything
  // real.) If you want to convert a passphrase to a key, use a suitable
  // package like bcrypt or scrypt.
  key, _ := hex.DecodeString(keya)
  ciphertext, _ := hex.DecodeString(encryptedBody)

  block, err := aes.NewCipher(key)
  if err != nil {
    panic(err)
  }

  // The IV needs to be unique, but not secure. Therefore it's common to
  // include it at the beginning of the ciphertext.
  if len(ciphertext) < aes.BlockSize {
    log.Fatal("ciphertext too short")
    return encryptedBody
  }
  iv := ciphertext[:aes.BlockSize]
  ciphertext = ciphertext[aes.BlockSize:]

  stream := cipher.NewCFBDecrypter(block, iv)

  // XORKeyStream can work in-place if the two arguments are the same.
  stream.XORKeyStream(ciphertext, ciphertext)
  return string(ciphertext)
}
