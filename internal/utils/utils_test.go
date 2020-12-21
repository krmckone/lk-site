package utils

import (
	"fmt"
	"os"
	"testing"
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
