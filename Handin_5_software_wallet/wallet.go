package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
)

type SecretKey struct {
	D []byte
	N []byte
}

// Helper function for hashing the password, prevent/slow bruteforce
func MasterHasher(password string, salt []byte) []byte {

	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(salt)
	bs := hash.Sum(nil)

	// repeat hash the password + salt
	for i := 0; i < 100000; i++ {
		hash := sha256.New()
		hash.Write(bs)
		bs = hash.Sum(nil)
	}

	return bs
}

func Generate(filename string, password string) string {

	keySize := 2000
	n, e, d, err := KeyGen(keySize)
	if err != nil {
		panic(fmt.Sprint("Error: ", err))
	}

	// struct to keep the secret key together
	private := SecretKey{
		D: d.Bytes(),
		N: n.Bytes(),
	}

	// Marshal the struct into data, for easier handling
	plaintext, err := json.Marshal(private)
	if err != nil {
		panic(fmt.Sprint("Marshal error: ", err))
	}

	//make a salt and hash for protection against dictionary attacks, prevents effective guessing
	//hash like there is no tomorrow for protect against brute fuckers
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		panic("Can't make salt number")
	}

	aes := MasterHasher(password, salt)

	//Encrypt using AES (blockcipher)
	iv, ciphertext, err := EncryptNumber(aes, plaintext)
	if err != nil {
		panic(fmt.Sprint("AES encryption failed: ", err))
	}

	// Combine the hashed thingies
	saveFile := append(salt, append(iv, ciphertext...)...)

	//Write to disk
	err = os.WriteFile(filename, saveFile, 0600)
	//Return public key N E
	return fmt.Sprint("Modulo N: ", n.String(), "\n Public key E: ", e.String())
}

// it unlocks a wallet and signs a message
func Sign(filename string, password string, msg []byte) string {
	//read the bytes from the file
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprint("Can't read file: ", err))
	}

	//get the components of the file
	salt := data[:16] //the salt should be the first 16 bit
	iv := data[16:32]
	ciphertext := data[32:] //evyerhitng after the 32 bits

	//unmarshall the sketchy RSA keys
	aes := MasterHasher(password, salt)

	//hjælp kan ikke finde ud af at push

	return "Signature"
}
