package bearpush

import (
	"fmt"
	"strings"
)

// TokenStrategy determines how the user tokens get validated.
type TokenStrategy int

const (
	noStrategy TokenStrategy = 1 << iota
	// StaticToken means there is a single token that will not change at runtime.
	StaticToken
	// RetrieveToken means the possible tokens will be provided by a script.
	RetrieveToken
	// VerifyToken means each token will get verified by a script.
	VerifyToken
)

// String converts a Strategy to a string.
func (s TokenStrategy) String() string {
	return strategyToString[s]
}

var toStrategy = map[string]TokenStrategy{
	"none":     noStrategy,
	"static":   StaticToken,
	"retrieve": RetrieveToken,
	"generate": RetrieveToken,
	"verify":   VerifyToken,
}

var strategyToString = map[TokenStrategy]string{
	noStrategy:    "none",
	StaticToken:   "static",
	RetrieveToken: "retrieve",
	VerifyToken:   "verify",
}

// UnmarshalYAML unmarshals YAML value to a TokenStrategy,
func (s *TokenStrategy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var t string
	err := unmarshal(&t)
	if err != nil {
		return err
	}

	if val, ok := toStrategy[strings.ToLower(t)]; ok {
		*s = val
		return nil
	}

	return fmt.Errorf("'%s' is not a valid token strategy", t)
}

// VerifyToken checks whether a token can be considered valid,
// according to current strategy.
func (p *Product) VerifyToken(token string, c *Context) (bool, error) {
	switch p.TokenSettings.Strategy {
	case StaticToken:
		return *p.TokenSettings.Value == token, nil
	case RetrieveToken:
		return false, nil
	}

	panic(fmt.Sprintf("Token strategy %v is not implemented", p.TokenSettings.Strategy))
}
