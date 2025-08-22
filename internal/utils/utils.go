package utils

// TODO: Update these to return errors rather than log
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Parameterizes specific values needed to load assets and configuration
// at runtime
type RuntimeConfig struct {
	AssetsPath  string
	ConfigsPath string
	BuildPath   string
}

var (
	repoRoot     string
	repoRootOnce sync.Once
)

// ReadDir returns the list of files in the directory, including subdirectories
// starting with path as a relative path from the repo root
func ReadDir(path string) ([]string, error) {
	paths := []string{}
	entries, err := os.ReadDir(MakePath(path))
	if err != nil {
		return []string{}, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			paths = append(paths, filepath.Join(path, entry.Name()))
		} else {
			subPaths, err := ReadDir(filepath.Join(path, entry.Name()))
			if err != nil {
				return []string{}, err
			}
			paths = append(paths, subPaths...)
		}
	}
	return paths, nil
}

// GetBasePageFiles returns the list of base page files in the assets directory
// This is statically defined and does not support subdirectories since these
// should not change often
func GetBasePageFiles(runtime RuntimeConfig) []string {
	files := []string{
		"base_page.html",
		"header.html",
		"footer.html",
		"topnav.html",
	}
	for i, file := range files {
		files[i] = filepath.Join(MakePath(runtime.AssetsPath), file)
	}
	return files
}

// GetComponentFiles returns the list of component files in the assets/components directory
func GetComponentFiles(runtime RuntimeConfig) ([]string, error) {
	return ReadDir(filepath.Join(MakePath(runtime.AssetsPath), "components"))
}

// GetRepoRoot returns the root directory of the repository. This value is used
// to enable consistent relative paths for file copying/creating throughout the
// templating process
func GetRepoRoot() string {
	repoRootOnce.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current working directory: %s", err)
		}

		for {
			_, err := os.Stat(filepath.Join(wd, "go.mod"))
			if err == nil {
				repoRoot = wd
				return
			}
			parent := filepath.Dir(wd)
			if parent == wd {
				log.Fatalf("could not find repo root")
			}
			wd = parent
		}
	})
	return repoRoot
}

func MakePath(path string) string {
	if strings.HasPrefix(path, GetRepoRoot()) {
		return path
	}
	return filepath.Join(GetRepoRoot(), path)
}

// SetupBuild generates the directories for the output artifacts and puts
// assets that do not need processing in the build directory; these assets
// are referred to by the output artifacts
func SetupBuild(runtime RuntimeConfig) error {
	assetDirs := []string{"images", "js", "shaders"}
	dirs := []string{}
	for _, dir := range assetDirs {
		dirs = append(dirs, filepath.Join(runtime.BuildPath, dir))
	}
	for _, dir := range dirs { // Maybe we could combine these loops
		if err := Clean(dir); err != nil {
			return fmt.Errorf("error cleaning directory %s: %s", dir, err)
		}
		if err := Mkdir(dir); err != nil {
			return fmt.Errorf("error making directory %s: %s", dir, err)
		}
	}
	for _, dir := range assetDirs {
		if err := CopyAssetToBuild(runtime, dir); err != nil {
			return fmt.Errorf("error copying asset to build: %s", err)
		}
	}
	return nil
}

// ReadFile wrapper for os.ReadFile
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(MakePath(path))
}

// WriteFile wrapper for ioutil.WriteFile
func WriteFile(path string, b []byte) error {
	if err := os.WriteFile(MakePath(path), b, 0644); err != nil {
		return fmt.Errorf("error writing file at %s: %s", MakePath(path), err)
	}
	return nil
}

// Mkdir wrapper for os.MkdirAll
func Mkdir(path string) error {
	if err := os.MkdirAll(MakePath(path), 0755); err != nil {
		return fmt.Errorf("error making directory at %s: %s", MakePath(path), err)
	}
	return nil
}

// Clean cleans the directory at path
func Clean(path string) error {
	if err := os.RemoveAll(MakePath(path)); err != nil {
		return fmt.Errorf("error cleaning path %s: %s", MakePath(path), err)
	}
	return nil
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

func CopyAssetToBuild(runtime RuntimeConfig, srcName string) error {
	return CopyFiles(
		filepath.Join(runtime.AssetsPath, srcName),
		filepath.Join(runtime.BuildPath, srcName),
	)
}

// Copies files and directories from srcPath to dstPath
func CopyFiles(srcPath, dstPath string) error {
	repoSrcPath := srcPath
	if !strings.HasPrefix(srcPath, GetRepoRoot()) {
		repoSrcPath = MakePath(srcPath)
	}
	repoDstPath := dstPath
	if !strings.HasPrefix(dstPath, GetRepoRoot()) {
		repoDstPath = MakePath(dstPath)
	}

	entries, err := os.ReadDir(repoSrcPath)
	if err != nil {
		return fmt.Errorf("error reading directory at %s: %s", repoSrcPath, err)
	}
	for _, entry := range entries {
		if !entry.Type().IsDir() {
			if err := copyFile(
				filepath.Join(repoSrcPath, entry.Name()),
				filepath.Join(repoDstPath, entry.Name()),
			); err != nil {
				return fmt.Errorf("error copying file from %s to %s: %s", MakePath(filepath.Join(repoSrcPath, entry.Name())), MakePath(filepath.Join(repoDstPath, entry.Name())), err)
			}
		} else {
			if err := Mkdir(filepath.Join(repoDstPath, entry.Name())); err != nil {
				return fmt.Errorf("error making directory at %s: %s", MakePath(filepath.Join(repoDstPath, entry.Name())), err)
			}
			if err := CopyFiles(
				filepath.Join(repoSrcPath, entry.Name()),
				filepath.Join(repoDstPath, entry.Name()),
			); err != nil {
				return fmt.Errorf("error copying files from %s to %s: %s", MakePath(filepath.Join(repoSrcPath, entry.Name())), MakePath(filepath.Join(repoDstPath, entry.Name())), err)
			}
		}
	}
	return nil
}

func copyFile(srcPath, dstPath string) error {
	sourceFileStat, err := os.Stat(MakePath(srcPath))
	if err != nil {
		return fmt.Errorf("error getting description for file at %s: %s", MakePath(srcPath), err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", MakePath(srcPath))
	}

	source, err := os.Open(MakePath(srcPath))
	if err != nil {
		return fmt.Errorf("error opening file at %s: %s", MakePath(srcPath), err)
	}
	defer source.Close()

	os.MkdirAll(filepath.Dir(MakePath(dstPath)), os.ModePerm)
	destination, err := os.Create(MakePath(dstPath))
	if err != nil {
		return fmt.Errorf("error creating file at %s: %s", MakePath(dstPath), err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("error copying file from %s to %s: %s", MakePath(srcPath), MakePath(dstPath), err)
	}
	return nil
}

// Invokes HTTP GET on the URL and returns the body as a string
func HttpGet(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP GET return code: %d", resp.StatusCode)
	}
	return resp, nil
}

func ReadHttpRespBody(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return fmt.Errorf("error in reading HTTP response body: %s", err)
	}
	return nil
}

// Generic filtering
func Filter[S ~[]E, E any](s S, f func(E) bool) []E {
	result := []E{}

	for i := range s {
		if f(s[i]) {
			result = append(result, s[i])
		}
	}

	return result
}
