package main

import "fmt"

func main() {
	fmt.Println("\n Testing Testing Testing")

	filename := "test.dat"
	password := "Pygme"

	// generate testing - This should generate a file with the name "test.dat" and if you open that, it should gibberish (encoded)
	//Remember to delete the file for before handing in.
	fmt.Println("\n Make that wallet! it's supposed to be slow...")
	publicKey := Generate(filename, password)
	fmt.Println("Public Key: ", publicKey)

	// Sign testing. If the signign doesn't work, the test should crash/panic or return an error, if it outputs fine, it's an success
	fmt.Println("\n Signing...")
	message := []byte("What's for Dinner?")
	signature := Sign(filename, password, message)

	fmt.Println("Message:", string(message))
	fmt.Println("Signature:", signature.String()[:10])
}
