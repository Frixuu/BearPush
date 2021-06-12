package templates

import (
	_ "embed"
	"log"

	"github.com/Frixuu/BearPush/v2/util"
)

//go:embed product.yml
var templateProduct string

// Generates contents of a file describing a new product.
func GenerateProductFile(name string) string {
	token, err := util.GenerateRandomToken(32)
	if err != nil {
		log.Fatalf("Error while generating random token: %s", err)
	}

	return util.Expand(templateProduct, map[string]string{
		"TOKEN": token,
	})
}
