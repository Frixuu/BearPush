package util

import (
	"log"
	"os"
)

// TryRemoveDir attempts to remove a directory, logging and swallowing potential errors.
func TryRemoveDir(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Printf("Cannot remove directory %s: %s", path, err)
	}
}
