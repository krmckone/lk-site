package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// ReadFile wrapper for os.ReadFile
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile wrapper for ioutil.WriteFile
func WriteFile(path string, b []byte) {
	if err := os.WriteFile(path, b, 0644); err != nil {
		log.Fatalf("Error writing file at %s: %s", path, err)
	}
}

// Mkdir wrapper for os.MkdirAll
func Mkdir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalf("Error making directory at %s: %s", path, err)
	}
}

// Clean cleans the directory at path
func Clean(path string) {
	if err := os.RemoveAll(path); err != nil {
		log.Fatalf("Error cleaning path %s: %s", path, err)
	}
}

// Returns the current eastern timestamp
func GetCurrentEasternTime() string {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("Error getting the current EST time: %s", err)
	}
	return time.Now().In(location).Format(time.RFC822)
}

// Returns the current year such as "2024"
func GetCurrentYear() string {
	return strconv.Itoa(time.Now().Year())
}

// Copies files and directories from srcPath to dstPath
func CopyFiles(srcPath, dstPath string) {
	entries, err := os.ReadDir(srcPath)
	if err != nil {
		log.Fatalf("Error reading directory at %s: %s", srcPath, err)
	}
	for _, entry := range entries {
		if !entry.Type().IsDir() {
			copyFile(
				fmt.Sprintf("%s/%s", srcPath, entry.Name()),
				fmt.Sprintf("%s/%s", dstPath, entry.Name()),
			)
		} else {
			Mkdir(fmt.Sprintf("%s/%s", dstPath, entry.Name()))
			CopyFiles(
				fmt.Sprintf("%s/%s", srcPath, entry.Name()),
				fmt.Sprintf("%s/%s", dstPath, entry.Name()),
			)
		}
	}
}

func copyFile(srcPath, dstPath string) {
	sourceFileStat, err := os.Stat(srcPath)
	if err != nil {
		log.Fatalf("Error getting description for file at %s: %s", srcPath, err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		log.Fatal(fmt.Errorf("%s is not a regular file", srcPath))
	}

	source, err := os.Open(srcPath)
	if err != nil {
		log.Fatalf("Error opening file at %s: %s", srcPath, err)
	}
	defer source.Close()

	destination, err := os.Create(dstPath)
	if err != nil {
		log.Fatalf("Error creating file at %s: %s", dstPath, err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		log.Fatalf("Error copying file from %s to %s: %s", srcPath, dstPath, err)
	}
}
