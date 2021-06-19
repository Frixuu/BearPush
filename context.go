package bearpush

import (
	"github.com/ReneKroon/ttlcache/v2"
	"go.uber.org/zap"
)

// Context stores current application state.
type Context struct {
	Config   *Config
	Logger   *zap.SugaredLogger
	Products map[string]*Product
	Cache    *ttlcache.Cache
}

// ContextFromConfig constructs a Context object from a loaded Config.
func ContextFromConfig(c *Config) (*Context, error) {
	p, err := LoadAllProducts(c.Path)
	if err != nil {
		return nil, err
	}

	cache := ttlcache.NewCache()
	cache.SkipTTLExtensionOnHit(true)

	return &Context{
		Config:   c,
		Logger:   zap.NewNop().Sugar(),
		Products: p,
		Cache:    cache,
	}, nil
}
