package main

import (
	"encoding/json"
	"math/big"
	"sync"
)

type Ledger struct {
	Accounts map[string]int      // known accounts
	SeenIDs  map[string]struct{} // Seend transactions
	lock     sync.Mutex
}

// A struct to combine the public key (n, e), helps with string encoding.
// Should act as the account ID
type PublicKey struct {
	N, E *big.Int
}

// Transaction hold serialized data
type SignedTransaction struct {
	Data      []byte   // This should always be a serialized transaction
	Signature *big.Int //RSA signature
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

	//unserialize the data
	var st SerializedTransaction
	err := json.Unmarshal(t.Data, &st)
	if err != nil {
		return //panic or something
	}

	// get public key from data (account identifier)
	var pk PublicKey
	json.Unmarshal([]byte(st.FromAccount), &pk)

	//Verify signature
	if !RSAVerify(t.Data, t.Signature, pk.E, pk.N) {
		return //goofy ah return (ignore transaction)
	}

	//
	if _, seen := l.SeenIDs[st.TxID]; seen {
		return // ignore transaction
	}

	l.SeenIDs[st.TxID] = struct{}{}
	l.Accounts[st.FromAccount] -= st.Amount
	l.Accounts[st.ToAccount] += st.Amount
}

// Helper method that serializes transaction data. A peer should call this function before signing or verifying0
func (st *SerializedTransaction) Serialization() []byte {
	data := SerializedTransaction{
		TxID:        st.TxID,
		FromAccount: st.FromAccount,
		ToAccount:   st.ToAccount,
		Amount:      st.Amount,
	}

	bytes, _ := json.Marshal(data)
	return bytes
}
