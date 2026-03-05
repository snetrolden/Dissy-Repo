package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"
)

func main() {

	fmt.Println("Testing Signature")
	k := 2000
	n, e, d, err := KeyGen(k)
	if err != nil {
		fmt.Print("error, KeyGen Error", err)
		return
	}

	originalMsg := []byte("Hello my name is Markiplier!")
	signature := RSASign(originalMsg, d, n)

	check := RSAVerify(originalMsg, signature, e, n)
	fmt.Print("Checking if message is correct:", check, "\n")

	fakeMessage := []byte("Hello my name is NOT Markiplier!")
	checkFake := RSAVerify(fakeMessage, signature, e, n)
	fmt.Print("Checking fake message:", checkFake, "\n")

	//---Measure hashing - AI was partially used for the measurements exercises. Espeacially the print statements
	fmt.Print("Measureing hashing speed \n")
	kb10 := make([]byte, 10*1024)
	rand.Read(kb10)

	iterations := 5000
	startHash := time.Now()
	for range iterations {
		h := sha256.New()
		h.Write(kb10)
		h.Sum(nil)
	}
	hashDuration := time.Since(startHash).Seconds()

	totalBitsHashed := float64(10 * 1024 * 8 * iterations)
	hashSpeedBps := totalBitsHashed / hashDuration
	fmt.Printf("Total time to hash %d times: %.4f seconds\n", iterations, hashDuration)
	fmt.Printf("SHA-256 Hashing Speed: %.2f bits/sec\n", hashSpeedBps)

	//---RSA Signing speed
	fmt.Print("RSA signing speed \n")
	signatureIterations := 100
	startRSATime := time.Now()
	for range signatureIterations {
		RSASign(kb10, d, n)
	}
	RSADuration := time.Since(startRSATime).Seconds()

	tpaSignature := RSADuration / float64(signatureIterations)
	fmt.Printf("Time for one 2000-bit RSA signature: %f seconds\n", tpaSignature)

	//--- Throughput thingy
	//Since we use 2000 bits each operation is ~ 2000 bits when using RSA
	rsaDirectThroughput := 2000.0 / tpaSignature
	fmt.Printf("Estimated RSA Direct Throughput: %.2f bits/sec\n", rsaDirectThroughput)
	fmt.Printf("Speed Multiplier (Hash vs RSA): %.2fx faster\n", hashSpeedBps/rsaDirectThroughput)
	//Discussion for 4. does hashing help?
	//Since hashing is much faster than RSA signing, and it can be used to reduce the size of the message signed, it can improve -
	// the performance of RSA signigng by enabling us to sign a smaller sized hash instead of the entire message, which becomes apparent for larger messages.
	// So in short yes, hashing improves efficiency significantly.
}
