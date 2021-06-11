package config

import (
	"log"
	"os"
	"path/filepath"
)

// Default configuration directory of this application.
// Can be overriden by the user.
var DefaultConfigDir = filepath.Join(defaultUserConfigRoot, "bearpush")

type Config struct {
	Path string
}

// Loads the app configuration.
// If the config files do not exist on disk, they will be created.
func Load(dir string) *Config {
	dirInfo, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0740); err != nil {
			log.Fatalf("Config directory '%s' cannot be created: %s", dir, err)
		}
	} else if !dirInfo.IsDir() {
		log.Fatalf("Config exists, but it is not a directory. Aborting")
	} else if err != nil {
		log.Fatalf("An error occured while retrieving config directory info: %s", err)
	}

	return &Config{
		Path: dir,
	}
}
