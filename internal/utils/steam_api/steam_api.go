package steamapi

import (
	"cmp"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/krmckone/ksite/internal/utils"
)

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
	resp, err := utils.HttpGet(baseUrl.String())
	if err != nil {
		return []SteamOwnedGame{}, err
	}
	target := SteamOwnedGamesResponse{}
	utils.ReadHttpRespBody(resp, &target)
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
