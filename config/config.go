package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// DefaultConfigDir represents a default configuration directory of this application.
// Can be overridden by the user.
var DefaultConfigDir = filepath.Join(defaultUserConfigRoot, "bearpush")

// Config stores configuration of this application.
type Config struct {
	Path string
}

// Load the app configuration.
// If the config files do not exist on disk, they will be created.
func Load(dir string) (*Config, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot determine absolute path for '%s': %w", dir, err)
	}
	dirInfo, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0740); err != nil {
			return nil, fmt.Errorf("config directory '%s' cannot be created: %w", dir, err)
		}
	} else if !dirInfo.IsDir() {
		return nil, fmt.Errorf("an entry for '%s' exists, but it is not a directory", dir)
	} else if err != nil {
		return nil, fmt.Errorf("cannot retrieve info about directory '%s': %w", dir, err)
	}

	c := Config{
		Path: dir,
	}

	return &c, nil
}
