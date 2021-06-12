package main

import (
	"github.com/Frixuu/BearPush/v2/config"
	"github.com/Frixuu/BearPush/v2/product"
)

// Context stores current application state.
type Context struct {
	Config   *config.Config
	Products map[string]*product.Product
}

// ContextFromConfig constructs a Context object from a loaded Config.
func ContextFromConfig(c config.Config) (*Context, error) {
	p, err := product.LoadAll(c.Path)
	if err != nil {
		return nil, err
	}

	return &Context{
		Config:   &c,
		Products: p,
	}, nil
}
