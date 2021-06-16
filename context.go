package bearpush

import "go.uber.org/zap"

// Context stores current application state.
type Context struct {
	Config   *Config
	Logger   *zap.SugaredLogger
	Products map[string]*Product
}

// ContextFromConfig constructs a Context object from a loaded Config.
func ContextFromConfig(c *Config) (*Context, error) {
	p, err := LoadAllProducts(c.Path)
	if err != nil {
		return nil, err
	}

	return &Context{
		Config:   c,
		Logger:   zap.NewNop().Sugar(),
		Products: p,
	}, nil
}
