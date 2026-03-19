package main

import (
	"encoding/json"
	"fmt"
)

type secretKey struct {
	D []byte
	N []byte
}

func Generate(filename string, password string) string {

	keySize := 2000
	n, e, d, err := KeyGen(keySize)
	if err != nil {
		panic(fmt.Sprint("Error: ", err))
	}

	private := secretKey{
		D: d.Bytes(),
		N: n.Bytes(),
	}

	plaintext, err := json.Marshal(private)
	if err != nil {
		panic(fmt.Sprint("Marshal error: ", err))
	}

	//make a salt for protection against dictionary attacks, prevents effective guessing

	//hash like there is no tomorrow for protect against brute fuckers

	//Encrypt using AES (blockcipher)

	// Save to a file (os.WriteFile) remember to use the FileName given in the method

	return fmt.Sprint("Modulo N: ", n.String(), "\n Public key E: ", e.String())
}

func Sign(filename string, password string, msg []byte) string {
	return "Signature"
}
