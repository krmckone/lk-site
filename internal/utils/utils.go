package utils

import (
	"log"
	"os"
	"strconv"
	"time"
)

// ReadFile wrapper for ioutil.ReadFile
func ReadFile(path string) []byte {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return []byte{}
	}
	return bytes
}

// WriteFile wrapper for ioutil.WriteFile
func WriteFile(path string, b []byte) {
	if err := os.WriteFile(path, b, 0644); err != nil {
		log.Fatal(err)
	}
}

// Mkdir wrapper for os.MkdirAll
func Mkdir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatal(err)
	}
}

// Clean cleans the directory at path
func Clean(path string) {
	if err := os.RemoveAll(path); err != nil {
		log.Fatal(err)
	}
}

// Returns the current eastern timestamp
func GetCurrentEasternTime() string {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}
	return time.Now().In(location).Format(time.RFC822)
}

// Returns the current year such as "2024"
func GetCurrentYear() string {
	return strconv.Itoa(time.Now().Year())
}
