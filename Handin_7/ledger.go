package main

import (
	"encoding/json"
	"math/big"
	"sync"
)

type Ledger struct {
	Accounts map[string]int //(key)string --> (val)int
	lock     sync.Mutex
}

// A struct to combine the public key (n, e), helps with string encoding
type Account struct {
	N, E *big.Int
}

// Transaction hold serialized data
type SignedTransaction struct {
	Data      *SerializedTransaction // Separate the serialized data
	Signature *big.Int               //RSA signature
}

// Sensitive transaction information (to be serialized)
type SerializedTransaction struct {
	TxID        string // Transaction identifier
	FromAccount string
	ToAccount   string
	Amount      int
}

func MakeLedger() *Ledger {
	ledger := new(Ledger)
	ledger.Accounts = make(map[string]int)

	return ledger
}

func (l *Ledger) Transaction(t *SignedTransaction) {
	l.lock.Lock()
	defer l.lock.Unlock()

	//add verification

	//check if txID has been seen before

	//if signature is invalid ignore transaction

	// if txID has already been seen/executed , ignore transaction (replay protection)

	l.Accounts[t.Data.FromAccount] -= t.Data.Amount
	l.Accounts[t.Data.ToAccount] += t.Data.Amount
}

// Helper method that serializes transaction data. A peer should call this function before signing or verifying0
func (st *SerializedTransaction) Serialization() []byte {
	data := SerializedTransaction{
		TxID:        st.TxID,
		FromAccount: st.ToAccount,
		ToAccount:   st.ToAccount,
		Amount:      st.Amount,
	}

	bytes, _ := json.Marshal(data)
	return bytes
}
