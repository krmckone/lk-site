package utils

// TODO: Update these to return errors rather than log
import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/krmckone/lk-site/internal/steamapi"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Parameterizes specific values needed to load assets and configuration
// at runtime
type RuntimeConfig struct {
	AssetsPath    string
	ConfigsPath   string
	BuildPath     string
	TemplateFuncs template.FuncMap
}

func NewRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		AssetsPath:    "assets",
		ConfigsPath:   "configs",
		BuildPath:     "build",
		TemplateFuncs: getTemplateFuncs(),
	}
}

func getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"makeHrefs":         makeHrefs,
		"makeNavTitle":      makeNavTitleFromHref,
		"getSteamDeckTop50": steamapi.GetSteamDeckTop50Wrapper,
	}
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
		dirs = append(dirs, filepath.Join(MakePath(runtime.BuildPath), dir))
	}
	if err := Clean(MakePath(runtime.BuildPath)); err != nil {
		return fmt.Errorf("error cleaning directory %s: %s", MakePath(runtime.BuildPath), err)
	}
	for _, dir := range dirs { // Maybe we could combine these loops
		if err := Mkdir(dir); err != nil {
			return fmt.Errorf("error making directory %s: %s", dir, err)
		}
	}
	for _, dir := range assetDirs {
		if err := CopyAssetToBuild(runtime, dir); err != nil {
			return fmt.Errorf("error copying %s to %s: %s", dir, runtime.BuildPath, err)
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
	return os.CopyFS(repoDstPath, os.DirFS(repoSrcPath))
}

func makeHrefs(path string) ([]string, error) {
	var hrefs []string

	assets, err := getAssets(path)
	if err != nil {
		return hrefs, err
	}

	sort.Strings(assets)
	for _, v := range assets {
		hrefs = append(hrefs, makeHref(v, path))
	}

	return hrefs, nil
}

func makeHref(assetName, originalPath string) string {
	_, file := path.Split(originalPath)
	return path.Join(string(os.PathSeparator), file, assetName)
}

func getAssets(path string) ([]string, error) {
	var assets []string

	dir, err := os.Open(MakePath(path))
	if err != nil {
		return assets, err
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return assets, err
	}

	for i, v := range files {
		// We only want to treat files as assets here. If there's nested
		// directories containg more assets, then getAssets needs to get
		// called with that nested path to handle that case separately
		if !v.IsDir() {
			assets = append(assets, strings.Split(files[i].Name(), ".")[0])
		}
	}

	return assets, nil
}

func makeNavTitleFromHref(assetHref string) string {
	_, file := path.Split(assetHref)
	caser := cases.Title(language.AmericanEnglish)
	return caser.String(
		strings.Join(strings.Split(file, "_"), " "),
	)
}
