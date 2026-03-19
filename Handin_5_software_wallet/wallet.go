package main

import (
	"crypto/sha256"
	"math/big"
)

type SecretKey struct {
	D *big.Int
	N *big.Int
}

func MasterHasher(password string) []byte {
	beenHashed := []byte(password)
	for i := 0; i < 100000; i++ {
		hash := sha256.Sum256(beenHashed)
	}
	return beenHashed
}

func Generate(filename string, password string) string {

	//generate a RSA key
	n, e, d, err := KeyGen(2048)
	return "Fake-it"

	//salting

	//hash 100k times

	//aes encrypt
}

func Sign(filename string, password string, msg []byte) string {
	return "Signature"
}
