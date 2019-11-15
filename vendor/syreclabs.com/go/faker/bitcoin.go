package faker

import (
	"crypto/sha256"
	"math/big"
)

type FakeBitcoin interface {
	Address() string // => "1GpEKM5UvD4XDLMirpNLoDnRVrGutogMj2"
	String() string  // String is an alias for Address.
}

type fakeBitcoin struct{}

func Bitcoin() FakeBitcoin {
	return fakeBitcoin{}
}

func (b fakeBitcoin) Address() string {
	v := make([]byte, 20)
	localRand.Read(v)

	return string(encodeBase58Check(v))
}

func encodeBase58Check(val []byte) []byte {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	// payload
	p := make([]byte, 1+len(val)) // version byte (0x00) + val
	copy(p[1:], val)
	h1 := sha256.Sum256(p)
	h2 := sha256.Sum256(h1[:])

	// value as []byte
	v := make([]byte, len(p)+4) // payload + first 4 bytes of h2
	copy(v, p)
	copy(v[len(p):], h2[:4])

	var res []byte
	x := new(big.Int).SetBytes(v)
	y := big.NewInt(58)
	m, zero := new(big.Int), new(big.Int)

	// convert to base58
	for x.Cmp(zero) > 0 {
		x, m = x.DivMod(x, y, m)
		res = append(res, alphabet[m.Int64()])
	}
	// append '1' for each leading zero byte in value
	for i := 0; v[i] == 0; i++ {
		res = append(res, alphabet[0])
	}
	// reverse
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}

	return res
}

func (b fakeBitcoin) String() string {
	return b.Address()
}
