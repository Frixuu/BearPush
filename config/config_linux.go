package config

import (
	"os"
	"path/filepath"
	"strings"
)

var defaultGlobalConfigRoots []string
var defaultUserConfigRoot string

func init() {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		defaultUserConfigRoot = os.Getenv("XDG_CONFIG_HOME")
	} else {
		defaultUserConfigRoot = filepath.Join(os.Getenv("HOME"), ".config")
	}
	if os.Getenv("XDG_CONFIG_DIRS") != "" {
		defaultGlobalConfigRoots = strings.Split(os.Getenv("XDG_CONFIG_DIRS"), ":")
	} else {
		defaultGlobalConfigRoots = []string{"/etc/xdg"}
	}
}
