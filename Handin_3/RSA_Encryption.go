package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// KeyGen generates an RSA key-pair (n ,e ,d) given a bit-lenght k
func KeyGen(k int) (n, e, d *big.Int, err error) {

	e = big.NewInt(3)
	one := big.NewInt(1)

	for {

		pLen := k / 2
		qLen := k - pLen

		// preventing error with variables
		var p, q *big.Int
		var err error

		//Finding gcd(3,p-1) = 1
		for {
			p, err = rand.Prime(rand.Reader, pLen)
			if err != nil {
				fmt.Print("Prime could not be generated for p", err)
				break
			}

			pMin := new(big.Int).Sub(p, one)
			// Cmp returns 0 if GCD == 1
			if new(big.Int).GCD(nil, nil, e, pMin).Cmp(one) == 0 {
				break
			}
		}

		//Finding  gcd(3, q-1) = 1
		for {
			q, err = rand.Prime(rand.Reader, qLen)
			if err != nil {
				fmt.Print("Prime could not be generated for q", err)
			}

			qMin := new(big.Int).Sub(q, one)
			if new(big.Int).GCD(nil, nil, e, qMin).Cmp(one) == 0 {
				break
			}
		}

		n = new(big.Int).Mul(p, q)

		// Bit lenght check
		if n.BitLen() == k {
			qMin := new(big.Int).Sub(q, one)
			pMin := new(big.Int).Sub(p, one)
			phi := new(big.Int).Mul(pMin, qMin)

			// get d
			d := new(big.Int).ModInverse(e, phi)
			if d != nil {
				return n, e, d, nil
			}
		}

	}
	// return key pair
}

func Encrypt(m, e, n *big.Int) *big.Int {
	return new(big.Int).Exp(m, e, n)

}

func Decrypt(c, d, n *big.Int) *big.Int {
	return new(big.Int).Exp(c, d, n)

}
