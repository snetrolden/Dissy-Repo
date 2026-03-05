package main

import (
	"fmt"
	"math/big"
)

func main() {
	fmt.Println("Testing RSA")
	k := 2048
	n, e, d, err := KeyGen(k)
	if err != nil {
		fmt.Print("error, KeyGen Error", err)
		return
	}

	fmt.Printf("Generated %d-bit RSA key. \n", n.BitLen())

	origMsg := big.NewInt(123456)

	ciphertext := Encrypt(origMsg, e, n)
	decryptMsg := Decrypt(ciphertext, d, n)

	if origMsg.Cmp(decryptMsg) == 0 {
		fmt.Println("Success! RSA decrypted Msg matches the original Msg")
	} else {
		fmt.Println("Failure! RSA decrypted Msg don't match the original Msg")
	}

	// AES TEST WITH SIMPLE NUMBER

	fmt.Println("\nTesting AES CTR with a number")

	// AES key (must be 16, 24, or 32 bytes)
	key := []byte("examplekey123456") // 16 bytes

	// Example number
	number := byte(12)

	fmt.Println("Original number:", number)

	// Convert number to plaintext bytes
	plaintext := []byte{number}

	// Encrypt
	iv, aesCipher, err := EncryptNumber(key, plaintext)
	if err != nil {
		fmt.Println("AES encryption error:", err)
		return
	}

	fmt.Println("Ciphertext:", aesCipher)

	// Decrypt
	decrypted, err := DecryptNumber(key, iv, aesCipher)
	if err != nil {
		fmt.Println("AES decryption error:", err)
		return
	}

	fmt.Println("Decrypted number:", decrypted[0])
}
