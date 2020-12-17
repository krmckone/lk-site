package utils

import (
	"io/ioutil"
	"log"
	"os"
)

// ReadFile wrapper for ioutil.ReadFile
func ReadFile(path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}
	}
	return bytes
}

// WriteFile wrapper for ioutil.WriteFile
func WriteFile(path string, b []byte) {
	if err := ioutil.WriteFile(path, b, 0644); err != nil {
		log.Fatal(err)
	}
}

// Mkdir wrapper for os.Mkdir
func Mkdir(path string) {
	if err := os.Mkdir("build", 0755); err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}
}

// Clean cleans the directory at path
func Clean(path string) {
	if err := os.RemoveAll("build"); err != nil {
		log.Fatal(err)
	}
}
