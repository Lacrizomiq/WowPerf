package models

type CharacterProfile struct {
	Name                                     string                     `json:"name"`
	Race                                     string                     `json:"race"`
	Class                                    string                     `json:"class"`
	ActiveSpecName                           string                     `json:"active_spec_name"`
	ActiveSpecRole                           string                     `json:"active_spec_role"`
	TreeID                                   int                        `json:"tree_id"`
	SpecID                                   int                        `json:"spec_id"`
	Gender                                   string                     `json:"gender"`
	Faction                                  string                     `json:"faction"`
	AchievementPoints                        int                        `json:"achievement_points"`
	HonorableKills                           int                        `json:"honorable_kills"`
	AvatarURL                                string                     `json:"avatar_url"`
	InsetAvatarURL                           string                     `json:"inset_avatar_url"`
	MainRawUrl                               string                     `json:"main_raw_url"`
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
	TalentLoadout                            *TalentLoadout             `json:"talentLoadout,omitempty"`
}

type TalentLoadout struct {
	LoadoutSpecID      int           `json:"loadout_spec_id"`
	TreeID             int           `json:"tree_id"`
	LoadoutText        string        `json:"loadout_text"`
	EncodedLoadoutText string        `json:"encoded_loadout_text"`
	ClassIcon          string        `json:"class_icon"`
	SpecIcon           string        `json:"spec_icon"`
	ClassTalents       []TalentNode  `json:"class_talents"`
	SpecTalents        []TalentNode  `json:"spec_talents"`
	SubTreeNodes       []SubTreeNode `json:"sub_tree_nodes"`
	HeroTalents        []HeroTalent  `json:"hero_talents"`
}

type TalentNode struct {
	NodeID    int           `json:"id"`
	NodeType  string        `json:"nodeType"`
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	PosX      int           `json:"posX"`
	PosY      int           `json:"posY"`
	MaxRanks  int           `json:"maxRanks"`
	EntryNode bool          `json:"entryNode"`
	ReqPoints int           `json:"reqPoints,omitempty"`
	FreeNode  bool          `json:"freeNode,omitempty"`
	Next      []int         `json:"next"`
	Prev      []int         `json:"prev"`
	Entries   []TalentEntry `json:"entries"`
	Rank      int           `json:"rank"`
}

type TalentEntry struct {
	EntryID      int    `json:"id"`
	DefinitionID int    `json:"definitionId"`
	MaxRanks     int    `json:"maxRanks"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	SpellID      int    `json:"spellId"`
	Icon         string `json:"icon"`
	Index        int    `json:"index"`
}

type HeroTalent struct {
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	PosX    int         `json:"posX"`
	PosY    int         `json:"posY"`
	Rank    int         `json:"rank"`
	Entries []HeroEntry `json:"entries"`
}

type HeroEntry struct {
	EntryID      int    `json:"id"`
	DefinitionID int    `json:"definitionId"`
	MaxRanks     int    `json:"maxRanks"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	SpellID      int    `json:"spellId"`
	Icon         string `json:"icon"`
	Index        int    `json:"index"`
}

type HeroNode struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	PosX      int         `json:"posX"`
	PosY      int         `json:"posY"`
	MaxRanks  int         `json:"maxRanks"`
	EntryNode bool        `json:"entryNode"`
	SubTreeID int         `json:"subTreeId"`
	Next      []int       `json:"next"`
	Prev      []int       `json:"prev"`
	Entries   []HeroEntry `json:"entries"`
	FreeNode  bool        `json:"freeNode"`
}

type SubTreeNode struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	Type    string         `json:"type"`
	Entries []SubTreeEntry `json:"entries"`
}

type SubTreeEntry struct {
	ID              int    `json:"id"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	TraitSubTreeID  int    `json:"traitSubTreeId"`
	AtlasMemberName string `json:"atlasMemberName"`
	Nodes           []int  `json:"nodes"`
	Rank            int    `json:"rank"`
}

type Gear struct {
	ItemLevelEquipped float64         `json:"item_level_equipped"`
	ItemLevelTotal    float64         `json:"item_level_total"`
	Items             map[string]Item `json:"items"`
}

type Item struct {
	ItemID      int     `json:"item_id"`
	ItemLevel   float64 `json:"item_level"`
	ItemQuality int     `json:"item_quality"`
	IsTwoHand   bool    `json:"is_two_hand"`
	IconName    string  `json:"icon_name"`
	IconURL     string  `json:"icon_url"`
	Name        string  `json:"name"`
	Enchant     *int    `json:"enchant,omitempty"`
	Gems        []int   `json:"gems,omitempty"`
	Bonuses     []int   `json:"bonuses,omitempty"`
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
}

type Segments struct {
	All    Segment `json:"all"`
	DPS    Segment `json:"dps"`
	Healer Segment `json:"healer"`
	Tank   Segment `json:"tank"`
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
	Affixes             []Affix `json:"affixes"`
	URL                 string  `json:"url"`
}

type Affix struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
