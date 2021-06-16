package bearpush

import (
	"log"
	"os"
	"path/filepath"

	"github.com/facebookgo/symwalk"
	"gopkg.in/yaml.v2"
)

// Product describes a type of entity Bearpush can process.
type Product struct {
	Script        string        `yaml:"process-script"`
	TokenSettings TokenSettings `yaml:"token"`
}

// TokenSettings describes how a Product validates incoming auth tokens.
type TokenSettings struct {
	Strategy TokenStrategy `yaml:"strategy"`
	Value    *string       `yaml:"static-value"`
	Script   *string       `yaml:"token-script"`
}

// LoadProductFromFile loads product manifest from a file under a provided path.
func LoadProductFromFile(path string) (*Product, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var p Product
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// LoadAllProducts parses all available product manifests.
func LoadAllProducts(basePath string) (map[string]*Product, error) {
	m := make(map[string]*Product)
	dir := filepath.Join(basePath, "products")
	err := symwalk.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		name := filepath.Base(path)
		ext := filepath.Ext(name)
		if ext != ".yml" {
			return nil
		}

		base := name[:len(name)-4]
		p, err := LoadProductFromFile(path)
		if err != nil {
			log.Printf("Cannot load product %s: %s\n", base, err)
			return nil
		}

		m[base] = p
		return nil
	})

	if err != nil {
		return nil, err
	}

	return m, nil
}
