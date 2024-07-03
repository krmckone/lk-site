package utils

// TODO: Update these to return errors rather than log
import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
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

type SteamOwnedGame struct {
	AppId                   int     `json:"appid"`
	Name                    string  `json:"name"`
	PlaytimeForever         float64 `json:"playtime_forever"`
	ImgIconUrl              string  `json:"img_icon_url"`
	PlaytimeWindowsForever  float64 `json:"playtime_windows_forever"`
	PlaytimeMacForever      float64 `json:"playtime_mac_forever"`
	PlaytimeLinuxForever    float64 `json:"playtime_linux_forever"`
	PlaytimeDeckForever     float64 `json:"playtime_deck_forever"`
	RTimeLastPlayed         int64   `json:"rtime_last_played"`
	FormattedTimeLastPlayed string
}

type SteamOwnedGamesResponse struct {
	Response struct {
		GamesCount int
		Games      []SteamOwnedGame
	}
}

func GetSteamOwnedGames() ([]SteamOwnedGame, error) {
	steam_api_key, present := os.LookupEnv("STEAM_API_KEY")
	if !present || steam_api_key == "" {
		return []SteamOwnedGame{}, fmt.Errorf("STEAM_API_KEY variable not present in env")
	}
	baseUrl, err := url.Parse("https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/")
	if err != nil {
		return []SteamOwnedGame{}, err
	}
	params := url.Values{}
	params.Add("key", steam_api_key)
	params.Add("steamid", "76561197988460908") // me
	params.Add("include_appinfo", "true")
	params.Add("include_extended_appinfo", "true")
	params.Add("include_played_free_games", "true")
	params.Add("include_free_sub", "true")
	params.Add("skip_unvetted_apps", "true")
	baseUrl.RawQuery = params.Encode()
	resp, err := httpGet(baseUrl.String())
	if err != nil {
		return []SteamOwnedGame{}, err
	}
	target := SteamOwnedGamesResponse{}
	readHttpRespBody(resp, &target)
	return target.Response.Games, nil
}

func GetTopFiftySteamDeckGames() ([]SteamOwnedGame, error) {
	games, err := GetSteamOwnedGames()
	if err != nil {
		return []SteamOwnedGame{}, err
	}
	if err := ProcessOwnedGames(games); err != nil {
		return []SteamOwnedGame{}, err
	}
	slices.SortFunc(games, func(a, b SteamOwnedGame) int {
		return cmp.Compare(b.PlaytimeDeckForever, a.PlaytimeDeckForever)
	})
	return games[:50], nil
}

// For adding any additional processing/formatting to the owned games data
func ProcessOwnedGames(games []SteamOwnedGame) error {
	for i := range games {
		// Steam deck played time from minutes to hours
		hours := games[i].PlaytimeDeckForever / 60.0
		truncate := fmt.Sprintf("%.2f", hours)
		t, err := strconv.ParseFloat(truncate, 64)
		if err != nil {
			return err
		}
		games[i].PlaytimeDeckForever = t

		lastPlayed := time.Unix(games[i].RTimeLastPlayed, 0)
		games[i].FormattedTimeLastPlayed = fmt.Sprintf("%s %d %d", lastPlayed.Month(), lastPlayed.Day(), lastPlayed.Year())
	}
	return nil
}

// Invokes HTTP GET on the URL and returns the body as a string
func httpGet(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP GET return code: %d", resp.StatusCode)
	}
	return resp, nil
}

func readHttpRespBody(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return fmt.Errorf("error in reading HTTP response body: %s", err)
	}
	return nil
}
