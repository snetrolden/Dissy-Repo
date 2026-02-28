package main

import (
	"fmt"
	"maps"
	"strconv"
	"time"
)

func main() {

	n := 10 // Number of Peers

	fmt.Println("---Starting 10 Peers---")
	//create and start Peers
	var peers []*Peer
	for i := 0; i < n; i++ {
		p := newPeer(3000 + i)
		peers = append(peers, p)
		go p.Start()
	}
	time.Sleep(1 * time.Second)

	fmt.Println("---Constructing Network of ALL Peers---")
	for i := 1; i < n; i++ {
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

	fmt.Println("---Sending Transactions---")
	for i := 0; i < 10; i++ {

		go sendTransactions(peers[i]) // go-routines for concurrency
	}
	time.Sleep(5 * time.Second)

	fmt.Println("---Check Ledgers---")

	peers[0].ledger.lock.Lock()
	baseLedger := make(map[string]int)
	maps.Copy(baseLedger, peers[0].ledger.Accounts)
	peers[0].ledger.lock.Unlock()

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

// Sends 10 transaction for a given Peer p
func sendTransactions(p *Peer) {
	accounts := []string{"Acc1", "Acc2", "Acc3", "Acc4", "Acc5"}

	for j := 0; j < 10; j++ {
		tx := &Transaction{
			ID:     "tx-" + strconv.Itoa(p.port) + "-msg-" + strconv.Itoa(j),
			From:   accounts[j%5],     // loop via modulo Shoutout to Theo
			To:     accounts[(j+1)%5], // Always go next
			Amount: 10,
		}
		p.FloodTransaction(tx)
		time.Sleep(10 * time.Millisecond) //
	}

}
