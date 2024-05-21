package utils

import (
	"fmt"
	"io"
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

// Copies files and directories from srcPath to dstPath
func CopyFiles(srcPath, dstPath string) {
	entries, err := os.ReadDir(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		fmt.Printf("%s/%s\n", srcPath, entry.Name())
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
		log.Fatal(err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		log.Fatal(fmt.Errorf("%s is not a regular file", srcPath))
	}

	source, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()

	destination, err := os.Create(dstPath)
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		log.Fatal(err)
	}
}
