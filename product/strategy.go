package product

import (
	"fmt"
	"strings"
)

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

	return fmt.Errorf("\"%s\" is not a valid token strategy", j)
}
