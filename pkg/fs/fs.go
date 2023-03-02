package fs

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// ReadFile reads from a file.
func ReadFile(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error while reading file: %v", err)
	}

	return content
}

// WriteToFile writes to a file.
func WriteToFile(content []byte, path string) error {
	if err := ioutil.WriteFile(path, content, 0o600); err != nil {
		log.Fatalf("error while writing to file: %v", err)
	}

	return nil
}

// CreateDirIfNotExist creates a directory if it does not exist.
func CreateDirIfNotExist(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			log.Fatalf("cannot create directory %s: %v", path, err)
		}
	}
}

// RenameFile renames a given file.
func RenameFile(oldName, newName string) {
	if err := os.Rename(oldName, newName); err != nil {
		log.Fatalf("cannot rename file %s: %v", oldName, err)
	}
}
