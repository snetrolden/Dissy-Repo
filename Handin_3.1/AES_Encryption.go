package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// Function we encrypt a number using AES in CTR mode, we needed aes and cipher packages from go.
// Returns the IV and the ciphertext.
func EncryptNumber(key []byte, plaintext []byte) (iv []byte, ciphertext []byte, err error) {

	// Step 1
	// Build the AES block cipher from the secret key K
	// (AES key must be 16, 24, or 32 bytes.)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// Step 2
	// Choose a random IV (initialization vector) the one where we say + 0, 1, 2
	// In CTR mode this is the starting counter value
	iv = make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, nil, err
	}

	// Step 3
	// CTR mode generates a keystream using
	// AES(key, IV), AES(key, IV+1), AES(key, IV+2), ...
	// the package does it for us, is just does. Very happy student.
	stream := cipher.NewCTR(block, iv)

	// Step 4
	// Ciphertext = Plaintext XOR Keystream
	// We XOR the Plaintext with the keystream, this is bascially
	// the encryption. If we have the same key and the same IV
	// we will get the same stream every time, but only we know the
	// key so people dont know what stream we get.
	//We create the ciphertext byt xoring the plaintext and the stream.
	ciphertext = make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	return iv, ciphertext, nil
}

// Time to decrypt by using the same key
// DecryptNumber decrypts ciphertext using AES CTR.
func DecryptNumber(key []byte, iv []byte, ciphertext []byte) ([]byte, error) {

	// Step 5
	// Create AES block cipher again using the same key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Step 6
	// Recreate the same CTR keystream using the same IV
	stream := cipher.NewCTR(block, iv)

	// Step 7
	// Plaintext = Ciphertext XOR Keystream
	// A traceback.
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}
