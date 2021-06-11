package config

import "os"

var defaultGlobalConfigRoots = []string{"/Library/Application Support"}
var defaultUserConfigRoot = os.Getenv("HOME") + "/Library/Application Support"
