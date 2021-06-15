package product

import (
	"fmt"
	"strings"
)

// TokenStrategy determines how the user tokens get validated.
type TokenStrategy int

const (
	// None is a noop strategy and should not be used.
	None TokenStrategy = 1 << iota
	// Static means there is a single token that will not change at runtime.
	Static
	// Retrieve means the possible tokens will be provided by a script.
	Retrieve
	// Verify means each token will get verified by a script.
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
