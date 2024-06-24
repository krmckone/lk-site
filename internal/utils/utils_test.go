package utils

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestWriteFile(t *testing.T) {
	expect := "TEST0"
	expectFile := "TEST_FILE.txt"
	WriteFile(expectFile, []byte(expect))
	b, err := ReadFile(expectFile)
	if err != nil {
		t.Errorf("Unexpected error from ReadFile: %s", err)
	}
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

func TestCopyFiles(t *testing.T) {
	dstPath := "build/test_path"
	srcPath := "../../assets/test"

	if _, err := os.Stat(dstPath); !os.IsNotExist(err) {
		t.Errorf("Expected %s to not exist and it does", dstPath)
	}
	Mkdir(dstPath)
	CopyFiles(srcPath, dstPath)
	dir, err := os.ReadDir(dstPath)
	if err != nil {
		t.Errorf("Unable to read dir %s", dstPath)
	}
	if len(dir) != 1 {
		t.Errorf("Unexpected entry in dir %s", dstPath)
	}
	nested := fmt.Sprintf("%s/%s", dstPath, dir[0].Name())
	nestedDir, err := os.ReadDir(nested)
	if err != nil {
		t.Errorf("Unable to read dir %s", nested)
	}
	expected := []string{"post_0.md", "post_1.md", "post_2.md"}
	actual := []string{}
	for _, file := range nestedDir {
		actual = append(actual, file.Name())
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
	Clean(dstPath)
}
