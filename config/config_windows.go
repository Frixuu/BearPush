package config

import "os"

var defaultGlobalConfigRoots = []string{os.Getenv("PROGRAMDATA")}
var defaultUserConfigRoot = os.Getenv("APPDATA")
