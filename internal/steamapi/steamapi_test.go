package steamapi

import (
	"fmt"
	"testing"
	"time"
)

func TestProcessOwnedGames(t *testing.T) {
	now := time.Now()
	truncated, err := truncateFloat(5 / 60.0)
	if err != nil {
		t.Fatalf("Unable to truncate %f: %s", 5/60.0, err)
	}
	cases := []struct {
		games  []SteamOwnedGame
		expect []SteamOwnedGame
	}{
		{
			[]SteamOwnedGame{
				{
					1,
					"TestGame1",
					100.0,
					"img_icon_url",
					0,
					0,
					5,
					5,
					now.Unix(),
					"",
				},
			},
			[]SteamOwnedGame{
				{
					1,
					"TestGame1",
					100.0,
					"img_icon_url",
					0,
					0,
					5,
					truncated,
					time.Now().Unix(),
					fmt.Sprintf("%s %d %d", now.Month(), now.Day(), now.Year()),
				},
			},
		},
	}

	for i, c := range cases {
		tName := fmt.Sprintf("%v: %v,%v", i, c.games, c.expect)
		t.Run(tName, func(t *testing.T) {
			ProcessOwnedGames(c.games)
			for j := range c.games {
				if c.games[j] != c.expect[j] {
					t.Errorf("Expected: %v, actual: %v", c.games[j], c.expect[j])
				}
			}
		})
	}
}
