package util

import "os"

// Expand replaces $VAR and ${VAR} in the source string with values provided in the map.
// If the map does not contain a match, uses environment variables as fallback.
func Expand(src string, mapping map[string]string) string {
	return os.Expand(src, func(p string) string {
		if val, ok := mapping[p]; ok {
			return val
		}
		return os.Getenv(p)
	})
}
