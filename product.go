package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Frixuu/BearPush/v2/config"
	"github.com/facebookgo/symwalk"
	"gopkg.in/yaml.v2"
)

type Product struct {
	TokenSettings TokenSettings `yaml:"token"`
}

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

type TokenSettings struct {
	Strategy TokenStrategy `yaml:"strategy"`
	Value    *string       `yaml:"static-value"`
	Script   *string       `yaml:"script"`
}

func (p *Product) VerifyToken(token string) bool {
	return true
}

// Loads product manifest from a file under a provided path.
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

// Parses all available product manifests.
func LoadAllProducts(c *config.Config) (map[string]*Product, error) {
	m := make(map[string]*Product)
	dir := filepath.Join(c.Path, "products")
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
