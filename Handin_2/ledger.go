package main

import (
	"sync"
)

type Ledger struct {
	Accounts map[string]int //(key)string --> (val)int 
	lock     sync.Mutex
}

type Transaction struct {
	ID     string // Transaction identifier
	From   string // The account to debit
	To     string // The account to credit
	Amount int 
}

func MakeLedger() *Ledger {
	ledger := new(Ledger)
	ledger.Accounts = make(map[string]int)

	return ledger
}

func (l *Ledger) Transaction(t *Transaction) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.Accounts[t.From] -= t.Amount
	l.Accounts[t.To] += t.Amount
}
