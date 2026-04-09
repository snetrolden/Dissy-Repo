package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"math/big"
	"time"
)

// To keep track of each user and differentiate between the ID and keys for each user.
type User struct {
	ID string
	N  *big.Int
	E  *big.Int
	D  *big.Int
}

//Right now, this is just a copy of the previous main method, so ignore until later!!

func main() {

	nr := 10 // Number of Peers

	fmt.Println("---Starting 10 Peers---")
	//create and start Peers
	var peers []*Peer
	for i := 0; i < nr; i++ {
		p := newPeer(3000 + i)
		peers = append(peers, p)
		go p.Start()
	}
	time.Sleep(1 * time.Second)

	fmt.Println("---Constructing Network of ALL Peers---")
	for i := 1; i < nr; i++ {
		peers[i].Connect("127.0.0.1", 3000)
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(3 * time.Second) // Wait for network to fully form

	// Check lenght of network to check if it's fully formed
	fullLenght := len(peers) - 1 //Each peer should be connected to 9 others in network of 10 hence -1
	allConnected := true
	for _, p := range peers {
		p.lock.Lock()

		if len(p.peers) != fullLenght {
			allConnected = false
			fmt.Println("Network not fully connected")
		}
		p.lock.Unlock()

	}
	if allConnected {
		fmt.Println("Fully connected network")
	}

	fmt.Println("--- Creating accounts---")
	//Could have been made with a loop, but that would make it harder to use cool names
	Grace := createAccounts(2000)
	Leon := createAccounts(2000)
	Nathan := createAccounts(2000)
	Emily := createAccounts(2000)

	// sending transactions
	fmt.Println("---Sending Transactions---")

	// Grace sends 500 to Leon
	go sendValidTransactions(peers[0], Grace, Leon, 500, "tx-1")
	time.Sleep(2 * time.Second)

	// Nathan sends 1000 to Emily
	go sendValidTransactions(peers[1], Nathan, Emily, 1000, "tx-2")
	time.Sleep(2 * time.Second)

	// Leon sends 500 to Grace but Grace hacks the message (Grace is gonna get rich!)
	go sendInvalidTransactions(peers[2], Leon, Grace, 500, "tx-invalid-1")
	time.Sleep(2 * time.Second)

	fmt.Println("---Check Ledgers---")

	peers[0].ledger.lock.Lock()
	baseLedger := make(map[string]int)
	maps.Copy(baseLedger, peers[0].ledger.Accounts)
	peers[0].ledger.lock.Unlock()

	peers[0].ledger.lock.Lock()
	fmt.Println("-Value of Ledger-")
	for ID, balance := range peers[0].ledger.Accounts {
		name := "Bruh"
		//Don't look at my shame of else-if statements!
		if ID == Grace.ID {
			name = "Grace"
		} else if ID == Leon.ID {
			name = "Leon"
		} else if ID == Nathan.ID {
			name = "Nathan"
		} else if ID == Emily.ID {
			name = "Emily"
		}
		fmt.Printf(" - %s: %d\n", name, balance)
	}

	allIdentical := true

	for _, p := range peers {
		p.ledger.lock.Lock()

		if !maps.Equal(baseLedger, p.ledger.Accounts) {
			allIdentical = false
			fmt.Println("Not identical Ledger")
		}

		p.ledger.lock.Unlock()
	}
	if allIdentical {

		fmt.Println("All Ledgers match!")
	}

}

// Sends a valid transactions through the network for given peer P
func sendValidTransactions(p *Peer, sender *User, receiver *User, amount int, txID string) {
	st := &SerializedTransaction{
		TxID:        txID,
		FromAccount: sender.ID,
		ToAccount:   receiver.ID,
		Amount:      amount,
	}

	// Marshal and Sign the transaction
	data := st.Serialization()
	sig := RSASign(data, sender.D, sender.N)

	tx := &SignedTransaction{
		Data:      data,
		Signature: sig,
	}
	fmt.Println("Valid Tx Flooding network")
	//flood the network
	p.FloodTransaction(tx)
}

// sends an invalid transaction
func sendInvalidTransactions(p *Peer, sender *User, receiver *User, amount int, txID string) {

	//true original transactions
	st := &SerializedTransaction{
		TxID:        txID,
		FromAccount: sender.ID,
		ToAccount:   receiver.ID,
		Amount:      amount,
	}
	// Marshal and Sign the transaction
	originalData := st.Serialization()
	sig := RSASign(originalData, sender.D, sender.N)

	stMalicious := &SerializedTransaction{
		TxID:        txID,
		FromAccount: sender.ID,
		ToAccount:   receiver.ID,
		Amount:      58008, // Maliciously inflated amount
	}
	maliciousData := stMalicious.Serialization()

	// make transaction with malicious data but real signature
	tx := &SignedTransaction{
		Data:      maliciousData,
		Signature: sig,
	}

	fmt.Println("Invalid Tx Flooding network")
	//flood the network, should be ignored
	p.FloodTransaction(tx)
}

// Helper method for creating keypairs for each account/user, saving the paris in a list
func createAccounts(k int) *User {

	n, e, d, _ := KeyGen(k)

	//Consider finding a better solution instead of having 2 structs for handling user information
	pk := &PublicKey{
		N: n,
		E: e,
	}

	id, _ := json.Marshal(pk)

	user := &User{
		ID: string(id),
		N:  n,
		E:  e,
		D:  d,
	}

	return user

}
