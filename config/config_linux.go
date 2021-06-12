package config

import (
	"os"
	"path/filepath"
)

var defaultUserConfigRoot = determineUserDirLinux()

func determineUserDirLinux() string {
	if val, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		return val
	}
	return filepath.Join(os.Getenv("HOME"), ".config")
}
