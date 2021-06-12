package product

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/facebookgo/symwalk"
	"gopkg.in/yaml.v2"
)

// Product describes a type of entity Bearpush can process.
type Product struct {
	Script        string        `yaml:"process-script"`
	TokenSettings TokenSettings `yaml:"token"`
}

// TokenStrategy determines how the user tokens get validated.
type TokenStrategy int

const (
	None TokenStrategy = 1 << iota
	Static
	Retrieve
	Verify
)

// String converts a Strategy to a string.
func (s TokenStrategy) String() string {
	return strategyToString[s]
}

var toStrategy = map[string]TokenStrategy{
	"none":     None,
	"static":   Static,
	"retrieve": Retrieve,
	"generate": Retrieve,
	"verify":   Verify,
}

var strategyToString = map[TokenStrategy]string{
	None:     "none",
	Static:   "static",
	Retrieve: "retrieve",
	Verify:   "verify",
}

// UnmarshalYAML unmarshals YAML value to a TokenStrategy,
func (s *TokenStrategy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	err := unmarshal(&j)
	if err != nil {
		return err
	}

	if val, ok := toStrategy[strings.ToLower(j)]; ok {
		*s = val
		return nil
	}

	return errors.New(fmt.Sprintf("\"%s\" is not a valid token strategy", j))
}

// TokenSettings describes how a Product validates incoming auth tokens.
type TokenSettings struct {
	Strategy TokenStrategy `yaml:"strategy"`
	Value    *string       `yaml:"static-value"`
	Script   *string       `yaml:"token-script"`
}

// VerifyToken checks whether a token can be considered valid,
// according to current strategy.
func (p *Product) VerifyToken(token string) bool {
	switch p.TokenSettings.Strategy {
	case Static:
		return *p.TokenSettings.Value == token
	}

	panic(fmt.Sprintf("Token strategy %v is not implemented", p.TokenSettings.Strategy))
}

// Loads product manifest from a file under a provided path.
func LoadFromFile(path string) (*Product, error) {
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

// LoadAll parses all available product manifests.
func LoadAll(basePath string) (map[string]*Product, error) {
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
		p, err := LoadFromFile(path)
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
