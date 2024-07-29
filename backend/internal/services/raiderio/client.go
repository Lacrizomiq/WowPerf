package raiderio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const baseURL = "https://raider.io/api/v1"

type Client struct {
	httpClient *http.Client
}

func NewCLient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

type CharacterProfile struct {
	Name                                     string                     `json:"name"`
	Race                                     string                     `json:"race"`
	Class                                    string                     `json:"class"`
	ActiveSpecName                           string                     `json:"active_spec_name"`
	ActiveSpecRole                           string                     `json:"active_spec_role"`
	Gender                                   string                     `json:"gender"`
	Faction                                  string                     `json:"faction"`
	AchievementPoints                        int                        `json:"achievement_points"`
	HonorableKills                           int                        `json:"honorable_kills"`
	ThumbnailURL                             string                     `json:"thumbnail_url"`
	Region                                   string                     `json:"region"`
	Realm                                    string                     `json:"realm"`
	ProfileURL                               string                     `json:"profile_url"`
	Gear                                     *Gear                      `json:"gear,omitempty"`
	Guild                                    *Guild                     `json:"guild,omitempty"`
	RaidProgression                          map[string]RaidProgression `json:"raid_progression,omitempty"`
	MythicPlusScoresBySeason                 []MythicPlusScoreSeason    `json:"mythic_plus_scores_by_season,omitempty"`
	MythicPlusRanks                          *MythicPlusRanks           `json:"mythic_plus_ranks,omitempty"`
	MythicPlusRecentRuns                     []MythicPlusRun            `json:"mythic_plus_recent_runs,omitempty"`
	MythicPlusBestRuns                       []MythicPlusRun            `json:"mythic_plus_best_runs,omitempty"`
	MythicPlusAlternateRuns                  []MythicPlusRun            `json:"mythic_plus_alternate_runs,omitempty"`
	MythicPlusHighestLevelRuns               []MythicPlusRun            `json:"mythic_plus_highest_level_runs,omitempty"`
	MythicPlusWeeklyHighestLevelRuns         []MythicPlusRun            `json:"mythic_plus_weekly_highest_level_runs,omitempty"`
	MythicPlusPreviousWeeklyHighestLevelRuns []MythicPlusRun            `json:"mythic_plus_previous_weekly_highest_level_runs,omitempty"`
	PreviousMythicPlusRanks                  *MythicPlusRanks           `json:"previous_mythic_plus_ranks,omitempty"`
}

type Gear struct {
	ItemLevelEquipped int `json:"item_level_equipped"`
	ItemLevelTotal    int `json:"item_level_total"`
}

type Guild struct {
	Name  string `json:"name"`
	Realm string `json:"realm"`
}

type RaidProgression struct {
	Summary            string `json:"summary"`
	TotalBosses        int    `json:"total_bosses"`
	NormalBossesKilled int    `json:"normal_bosses_killed"`
	HeroicBossesKilled int    `json:"heroic_bosses_killed"`
	MythicBossesKilled int    `json:"mythic_bosses_killed"`
}

type MythicPlusScoreSeason struct {
	Season   string   `json:"season"`
	Scores   Scores   `json:"scores"`
	Segments Segments `json:"segments"`
}

type Scores struct {
	All    float64 `json:"all"`
	DPS    float64 `json:"dps"`
	Healer float64 `json:"healer"`
	Tank   float64 `json:"tank"`
	Spec0  float64 `json:"spec_0"`
	Spec1  float64 `json:"spec_1"`
	Spec2  float64 `json:"spec_2"`
	Spec3  float64 `json:"spec_3"`
}

type Segments struct {
	All    Segment `json:"all"`
	DPS    Segment `json:"dps"`
	Healer Segment `json:"healer"`
	Tank   Segment `json:"tank"`
	Spec0  Segment `json:"spec_0"`
	Spec1  Segment `json:"spec_1"`
	Spec2  Segment `json:"spec_2"`
	Spec3  Segment `json:"spec_3"`
}

type Segment struct {
	Score float64 `json:"score"`
	Color string  `json:"color"`
}

type MythicPlusRanks struct {
	Overall     Rank `json:"overall"`
	Tank        Rank `json:"tank"`
	Healer      Rank `json:"healer"`
	DPS         Rank `json:"dps"`
	Class       Rank `json:"class"`
	ClassTank   Rank `json:"class_tank"`
	ClassHealer Rank `json:"class_healer"`
	ClassDPS    Rank `json:"class_dps"`
}

type Rank struct {
	World  int `json:"world"`
	Region int `json:"region"`
	Realm  int `json:"realm"`
}

type MythicPlusRun struct {
	Dungeon             string  `json:"dungeon"`
	ShortName           string  `json:"short_name"`
	MythicLevel         int     `json:"mythic_level"`
	CompletedAt         string  `json:"completed_at"`
	ClearTimeMS         int     `json:"clear_time_ms"`
	NumKeystoneUpgrades int     `json:"num_keystone_upgrades"`
	Score               float64 `json:"score"`
	URL                 string  `json:"url"`
}

func (c *Client) GetCharacterProfile(region, realm, name string, fields []string) (*CharacterProfile, error) {
	endpoint := fmt.Sprintf("%s/characters/profile", baseURL)

	params := url.Values{}
	params.Add("region", region)
	params.Add("realm", realm)
	params.Add("name", name)
	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}

	resp, err := c.httpClient.Get(fmt.Sprintf("%s?%s", endpoint, params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var profile CharacterProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &profile, nil
}
