package config

import (
	"log"
	"os"
	"path/filepath"
)

// Default configuration directory of this application.
// Can be overridden by the user.
var DefaultConfigDir = filepath.Join(defaultUserConfigRoot, "bearpush")

type Config struct {
	Path string
}

// Load the app configuration.
// If the config files do not exist on disk, they will be created.
func Load(dir string) (Config, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("Error while loading config: Cannot determine absolute path for %s", dir)
	}
	dirInfo, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0740); err != nil {
			log.Fatalf("Config directory '%s' cannot be created: %s", dir, err)
		}
	} else if !dirInfo.IsDir() {
		log.Fatalf("An entry for %s exists, but it is not a directory. Aborting", dir)
	} else if err != nil {
		log.Fatalf("An error occurred while retrieving info for dir %s: %s", dir, err)
	}

	c := Config{
		Path: dir,
	}

	return c, nil
}
