package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetRepoRoot(t *testing.T) {
	expected, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Errorf("Unexpected error from filepath.Abs: %s", err)
	}
	actual := GetRepoRoot()
	if !strings.HasPrefix(actual, expected) {
		t.Errorf("Expected %s, got: %s", expected, actual)
	}
}

func TestMakePath(t *testing.T) {
	expected, err := filepath.Abs(filepath.Join("..", "..", "test"))
	if err != nil {
		t.Errorf("Unexpected error from filepath.Abs: %s", err)
	}
	actual := MakePath("test")
	if actual != expected {
		t.Errorf("Expected %s, got: %s", expected, actual)
	}
}

func TestSetupBuild(t *testing.T) {
	if err := SetupBuild(); err != nil {
		t.Errorf("Unexpected error from SetupBuild: %s", err)
	}
	dir, err := os.ReadDir(MakePath("build"))
	if err != nil {
		t.Errorf("Unexpected error from os.ReadDir: %s", err)
	}
	expected := []string{"images", "js", "shaders"}
	actual := []string{}
	for _, d := range dir {
		actual = append(actual, d.Name())
	}
	if !slices.Equal(expected, actual) {
		t.Errorf("Expected %s, got: %s", expected, actual)
	}
	if err := Clean("build"); err != nil {
		t.Errorf("Unexpected error from Clean: %s", err)
	}
}

func TestWriteFile(t *testing.T) {
	expectBody := "TEST0"
	expectFile := "TEST_FILE.txt"
	if err := WriteFile(expectFile, []byte(expectBody)); err != nil {
		t.Errorf("Unexpected error from WriteFile: %s", err)
	}
	b, err := ReadFile(expectFile)
	if err != nil {
		t.Errorf("Unexpected error from ReadFile: %s", err)
	}
	if string(b) != expectBody {
		t.Errorf("Expected %s: %s, got: %s", expectFile, expectBody, string(b))
	}
	if err := os.Remove(MakePath(expectFile)); err != nil {
		t.Errorf("Unexpected error from os.Remove: %s", err)
	}
	if err := Clean(MakePath(expectFile)); err != nil {
		t.Errorf("Unexpected error from Clean: %s", err)
	}
}

func TestMkdir(t *testing.T) {
	expect := "test_path"
	if err := Mkdir(expect); err != nil {
		t.Errorf("Unexpected error from Mkdir: %s", err)
	}
	if _, err := os.Stat(MakePath(expect)); os.IsNotExist(err) {
		t.Errorf("Expected %s to exist and it does not", MakePath(expect))
	}
	if err := Clean(expect); err != nil {
		t.Errorf("Unexpected error from Clean: %s", err)
	}
	if _, err := os.Stat(MakePath(expect)); !os.IsNotExist(err) {
		t.Errorf("Expected %s to not exist and it does", MakePath(expect))
	}
	if err := os.RemoveAll(MakePath(expect)); err != nil {
		t.Errorf("Unexpected error from os.RemoveAll: %s", err)
	}
}

func TestClean(t *testing.T) {
	expectPath := "test_path"
	expectFile := "TEST_FILE.txt"
	if err := Mkdir(expectPath); err != nil {
		t.Errorf("Unexpected error from Mkdir: %s", err)
	}
	if err := WriteFile(filepath.Join(expectPath, expectFile), []byte{0, 1, 2}); err != nil {
		t.Errorf("Unexpected error from WriteFile: %s", err)
	}
	if _, err := os.Stat(MakePath(filepath.Join(expectPath, expectFile))); os.IsNotExist(err) {
		t.Errorf("Expected %s to exist and it does not", MakePath(expectPath))
	}
	if err := Clean(expectPath); err != nil {
		t.Errorf("Unexpected error from Clean: %s", err)
	}
	if _, err := os.Stat(MakePath(filepath.Join(expectPath, expectFile))); !os.IsNotExist(err) {
		t.Errorf("Expected %s to not exist and it does", MakePath(filepath.Join(expectPath, expectFile)))
	}
	if _, err := os.Stat(MakePath(expectPath)); !os.IsNotExist(err) {
		t.Errorf("Expected %s to not exist and it does", MakePath(expectPath))
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
	srcPath := "assets/test"
	dstPath := "build/test_path"

	if err := Clean(dstPath); err != nil {
		t.Errorf("Unexpected error from Clean: %s", err)
	}
	if err := Mkdir(dstPath); err != nil {
		t.Errorf("Unexpected error from Mkdir: %s", err)
	}
	if _, err := os.Stat(MakePath(dstPath)); os.IsNotExist(err) {
		t.Errorf("Expected %s to exist and it does not", MakePath(dstPath))
	}
	if err := CopyFiles(srcPath, dstPath); err != nil {
		t.Errorf("Unexpected error from CopyFiles: %s", err)
	}

	dir, err := os.ReadDir(MakePath(dstPath))
	if err != nil {
		t.Errorf("unable to read dir %s: %s", MakePath(dstPath), err)
	}
	if len(dir) != 1 {
		t.Errorf("unexpected entry in dir %s: %s", MakePath(dstPath), dir)
	}

	nestedPath := filepath.Join(MakePath(dstPath), dir[0].Name())
	nestedDir, err := os.ReadDir(nestedPath)
	if err != nil {
		t.Errorf("unable to read dir %s: %s", nestedPath, err)
	}
	expected := []string{"post_0.md", "post_1.md", "post_2.md"}
	actual := []string{}
	for _, file := range nestedDir {
		actual = append(actual, file.Name())
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
	if err := Clean(dstPath); err != nil {
		t.Errorf("Unexpected error from Clean: %s", err)
	}
	if _, err := os.Stat(MakePath(dstPath)); !os.IsNotExist(err) {
		t.Errorf("Expected %s to not exist and it does", MakePath(dstPath))
	}

	if err := Clean(MakePath(dstPath)); err != nil {
		t.Errorf("Unexpected error from Clean: %s", err)
	}
}
