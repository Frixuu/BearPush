package templates

import (
	_ "embed" // For embedding template files

	"github.com/Frixuu/BearPush/util"
)

//go:embed product.yml
var templateProduct string

// Generates contents of a file describing a new product.
func GenerateProductFile(name string) string {
	tokenLength := util.RandInt(28, 35)
	token, _ := util.GenerateRandomToken(tokenLength)

	return util.Expand(templateProduct, map[string]string{
		"TOKEN": token,
	})
}
