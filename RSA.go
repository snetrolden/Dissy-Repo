package main

import (
	"crypto/rand"
	"math/big"
)

func KeyGen(k int) (n, e, d *big.Int, err error) {

	e = big.NewInt(3)

	one := big.NewInt(1)
	for {
		pLen := k / 2
		qLen := k - pLen

		var p, q *big.Int

		//generate p
		for {
			p, err = rand.Prime(rand.Reader, pLen)

			if err != nil {
				return nil, nil, nil, err
			}

			pMin := new(big.Int).Sub(p, one)

			if new(big.Int).GCD(nil, nil, e, pMin).Cmp(one) == 0 {
				break
			}
		}
		//generate q
		for {
			q, err = rand.Prime(rand.Reader, qLen)

			if err != nil {
				return nil, nil, nil, err
			}

			qMin := new(big.Int).Sub(q, one)

			if new(big.Int).GCD(nil, nil, e, qMin).Cmp(one) == 0 {
				break
			}
		}

		//n =p*q
		n = new(big.Int).Mul(p, q)

		if n.BitLen() == k {
			//(p-1)(q-1)
			pMin := new(big.Int).Sub(p, one)
			qMin := new(big.Int).Sub(q, one)
			phi := new(big.Int).Mul(pMin, qMin)

			d = new(big.Int).ModInverse(e, phi)
			if d != nil {
				break
			}
		}

	}
	return n, e, d, nil
}

func Encrypt(m, e, n *big.Int) *big.Int {
	return new(big.Int).Exp(m, e, n)
}

func Decrypt(c, d, n *big.Int) *big.Int {
	return new(big.Int).Exp(c, d, n)
}
