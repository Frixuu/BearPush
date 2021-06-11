package templates

import (
	"crypto/rand"
	_ "embed"
	"log"
	"math/big"
	"strings"
)

//go:embed product.yml
var templateProduct string

func generateRandomToken(n int) (string, error) {
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

// Generates contents of a file describing a new product.
func GenerateProductFile(name string) string {
	token, err := generateRandomToken(32)
	if err != nil {
		log.Fatalf("Error while generating random token: %s", err)
	}

	contents := strings.ReplaceAll(templateProduct, "${{TOKEN}}", token)
	return contents
}
