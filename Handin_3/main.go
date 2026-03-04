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
		fmt.Print("Success! RSA decrypted Msg matches the original Msg")
	} else {
		fmt.Print("Failure! RSA decrypted Msg don't match the original Msg")
	}

}
