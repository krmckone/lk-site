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

	"github.com/gofrs/flock"
)

func NewTestRuntime() RuntimeConfig {
	return RuntimeConfig{
		AssetsPath:  "test/assets",
		BuildPath:   "test/build",
		ConfigsPath: "test/configs",
	}
}

func TestReadDir(t *testing.T) {
	actual, err := ReadDir(filepath.Join("assets", "components"))
	if err != nil {
		t.Errorf("Unexpected error from ReadDir: %s", err)
	}
	expected := []string{filepath.Join("assets", "components", "steam_deck_top_50.html")}
	for i, file := range actual {
		if file != expected[i] {
			t.Errorf("Expected: %s, actual: %s", expected[i], file)
		}
	}
}

func TestGetBasePageFiles(t *testing.T) {
	runtime := NewTestRuntime()
	actual := GetBasePageFiles(runtime)
	expected := []string{
		filepath.Join(MakePath(runtime.AssetsPath), "base_page.html"),
		filepath.Join(MakePath(runtime.AssetsPath), "header.html"),
		filepath.Join(MakePath(runtime.AssetsPath), "footer.html"),
		filepath.Join(MakePath(runtime.AssetsPath), "topnav.html"),
	}
	for i, file := range actual {
		if file != expected[i] {
			t.Errorf("Expected: %s, actual: %s", expected[i], file)
		}
	}
}

func TestGetComponentFiles(t *testing.T) {
	runtime := NewTestRuntime()
	actual, err := GetComponentFiles(runtime)
	if err != nil {
		t.Errorf("Unexpected error from GetComponentFiles: %s", err)
	}
	expected := []string{
		filepath.Join(MakePath(runtime.AssetsPath), "components", "test_component.html"),
	}
	for i, file := range actual {
		if file != expected[i] {
			t.Errorf("Expected: %s, actual: %s", expected[i], file)
		}
	}
}

func TestGetRepoRoot(t *testing.T) {
	runtime := NewTestRuntime()
	t.Cleanup(func() {
		if err := Clean(MakePath(runtime.BuildPath)); err != nil {
			t.Errorf("Cleanup Unexpected error from Clean: %s", err)
		}
	})
	expected, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Errorf("Unexpected error from filepath.Abs: %s", err)
	}
	actual := GetRepoRoot()
	if !strings.HasPrefix(actual, expected) {
		t.Errorf("Expected %s, got: %s", expected, actual)
	}
}

func makeExpectedPath(input string) (string, error) {
	return filepath.Abs(filepath.Join("..", "..", input))
}

func TestMakePath(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{}
	for _, input := range []string{"test", "test/test"} {
		expected, err := makeExpectedPath(input)
		if err != nil {
			t.Errorf("Unexpected error from makeExpected: %s", err)
		}
		cases = append(cases, struct {
			input    string
			expected string
		}{input: input, expected: expected})
	}
	for _, c := range cases {
		actual := MakePath(c.input)
		if actual != c.expected {
			t.Errorf("Expected %s, got: %s", c.expected, actual)
		}
	}
}

func TestSetupBuild(t *testing.T) {
	runtime := NewTestRuntime()
	lock := flock.New(MakePath(runtime.BuildPath))
	defer lock.Unlock()

	locked, err := lock.TryLock()
	if err != nil {
		t.Errorf("Unexpected error from TryLock: %s", err)
	}
	if !locked {
		t.Errorf("Expected lock to be acquired")
	}

	t.Cleanup(func() {
		if err := Clean(MakePath(runtime.BuildPath)); err != nil {
			t.Errorf("Cleanup Unexpected error from Clean: %s", err)
		}
	})

	if err := SetupBuild(runtime); err != nil {
		t.Errorf("Unexpected error from SetupBuild: %s", err)
	}
	dir, err := os.ReadDir(MakePath(runtime.BuildPath))
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
}

func TestWriteFile(t *testing.T) {
	runtime := NewTestRuntime()
	t.Cleanup(func() {
		if err := Clean(MakePath(runtime.BuildPath)); err != nil {
			t.Errorf("Cleanup Unexpected error from Clean: %s", err)
		}
	})
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
	runtime := NewTestRuntime()
	t.Cleanup(func() {
		if err := Clean(MakePath(runtime.BuildPath)); err != nil {
			t.Errorf("Cleanup Unexpected error from Clean: %s", err)
		}
	})
	srcPath := filepath.Join(runtime.AssetsPath, "pages")
	dstPath := filepath.Join(runtime.BuildPath, "pages")

	if err := Clean(dstPath); err != nil {
		t.Errorf("Unexpected error from Clean: %s", err)
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

func TestMakeNavTitleFromHref(t *testing.T) {
	runtime := NewTestRuntime()
	t.Cleanup(func() {
		if err := Clean(MakePath(runtime.BuildPath)); err != nil {
			t.Errorf("Unexpected error from Clean: %s", err)
		}
	})
	cases := []struct {
		in, expect string
	}{
		{
			"posts/page_0", "Page 0",
		}, {
			"posts/page 1", "Page 1",
		}, {
			"index_page_zero", "Index Page Zero",
		},
	}
	for _, c := range cases {
		actual := makeNavTitleFromHref(c.in)
		if actual != c.expect {
			t.Errorf("Expected: %s, actual: %s", c.expect, actual)
		}
	}
}

func TestMakeHref(t *testing.T) {
	runtime := NewTestRuntime()
	t.Cleanup(func() {
		if err := Clean(MakePath(runtime.BuildPath)); err != nil {
			t.Errorf("Unexpected error from Clean: %s", err)
		}
	})
	cases := []struct {
		assetName, originalPath, expect string
	}{
		{
			"testing_1_2_3", "test/posts", "/posts/testing_1_2_3",
		},
		{
			"file_name", "test123/xyz", "/xyz/file_name",
		},
		{
			"index", "website/test123/files", "/files/index",
		},
	}
	for _, c := range cases {
		actual := makeHref(c.assetName, c.originalPath)
		if actual != c.expect {
			t.Errorf("Expected: %s, actual: %s", c.expect, actual)
		}
	}
}

func TestMakeHrefs(t *testing.T) {
	runtime := NewTestRuntime()
	t.Cleanup(func() {
		if err := Clean(MakePath(runtime.BuildPath)); err != nil {
			t.Errorf("Unexpected error from Clean: %s", err)
		}
	})
	cases := []struct {
		path   string
		expect []string
	}{
		{
			filepath.Join(MakePath(runtime.AssetsPath), "test", "pages"),
			[]string{"/pages/post_0", "/pages/post_1", "/pages/post_2"},
		},
	}
	for _, c := range cases {
		actual, err := makeHrefs(c.path)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("Expected: %s, actual: %s", c.expect, actual)
		}
	}
}
