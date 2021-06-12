package templates

import (
	_ "embed" // For embedding template files
	"log"

	"github.com/Frixuu/BearPush/v2/util"
)

//go:embed product.yml
var templateProduct string

// Generates contents of a file describing a new product.
func GenerateProductFile(name string) string {
	tokenLength := util.RandInt(28, 35)
	token, err := util.GenerateRandomToken(tokenLength)
	if err != nil {
		log.Fatalf("Error while generating random token: %s", err)
	}

	return util.Expand(templateProduct, map[string]string{
		"TOKEN": token,
	})
}
