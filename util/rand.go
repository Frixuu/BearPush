package util

import (
	"crypto/rand"
	"math/big"
)

// Generates a new string of n random chars.
// The user should not assume what characters can or cannot appear in the token.
func GenerateRandomToken(n int) (string, error) {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	max := big.NewInt(int64(len(chars)))
	token := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		token[i] = chars[num.Int64()]
	}

	return string(token), nil
}
