package utils

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestReadFileErr(t *testing.T) {
	b := ReadFile("$$CANNOTEXIST$$")
	if len(b) > 0 {
		t.Errorf("Expected bytes length 0, got: %d", len(b))
	}
}

func TestWriteFile(t *testing.T) {
	expect := "TEST0"
	expectFile := "TEST_FILE.txt"
	WriteFile(expectFile, []byte(expect))
	b := ReadFile(expectFile)
	if string(b) != expect {
		t.Errorf("Expected %s: %s, got: %s", expectFile, expect, string(b))
	}
	os.Remove(expectFile)
}

func TestMkdir(t *testing.T) {
	expect := "test_path"
	Mkdir(expect)
	if _, err := os.Stat(expect); os.IsNotExist(err) {
		t.Errorf("Expected %s to exist and it does not", expect)
	}
	os.Remove(expect)
}

func TestClean(t *testing.T) {
	expectPath := "test_path"
	expectFile := "TEST_FILE.txt"
	Mkdir(expectPath)
	WriteFile(fmt.Sprintf("%s/%s", expectPath, expectFile), []byte{})
	Clean(expectPath)
	if _, err := os.Stat(expectFile); !os.IsNotExist(err) {
		t.Errorf("Expected %s to not exist and it does", expectFile)
	}
	if _, err := os.Stat(expectPath); !os.IsNotExist(err) {
		t.Errorf("Expected %s to not exist and it does", expectPath)
	}
}

func TestGetCurrentEasternTime(t *testing.T) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Errorf("Error in setting up location: %s", err)
	}
	expected := time.Now().In(location).Format(time.RFC822)
	actual := GetCurrentEasternTime()
	if expected != actual {
		t.Errorf("Expected: %s, got: %s", expected, actual)
	}
}
func TestGetCurrentYear(t *testing.T) {
	expected := strconv.Itoa(time.Now().Year())
	actual := GetCurrentYear()
	if expected != actual {
		t.Errorf("Expected: %s, got :%s", expected, actual)
	}
}
