package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

func Encrypt(keyStr string, textStr string) ([]byte, error) {
	hash := sha256.Sum256([]byte(keyStr))
	text := []byte(textStr)
	key := hash[:32]

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logger.Error(err)

		return nil, err
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.

	return gcm.Seal(nonce, nonce, text, nil), nil
}

func Decrypt(keyStr string, textRow []byte) ([]byte, error) {
	hash := sha256.Sum256([]byte(keyStr))
	key := hash[:32]
	ciphertext := textRow

	c, err := aes.NewCipher(key)
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	return plaintext, nil
}
