package models

// MythicPlusRun represents a mythic run from the Raider.io API.
// https://raider.io/api/v1/mythic-plus/run-details?season=season-tww-1&id=2608796 for example
// It allows us to get the full details of a run, including the dungeon, the roster, the items, talents, the affixes, etc.
type MythicPlusRun struct {
	ClearTimeMS        int      `json:"clear_time_ms"`
	CompletedAt        string   `json:"completed_at"`
	DeletedAt          *string  `json:"deleted_at"`
	Dungeon            Dungeon  `json:"dungeon"`
	Faction            string   `json:"faction"`
	KeystoneRunID      int      `json:"keystone_run_id"`
	KeystoneTeamID     int      `json:"keystone_team_id"`
	KeystoneTimeMS     int      `json:"keystone_time_ms"`
	LoggedSources      []string `json:"loggedSources"`
	LoggedDetails      *string  `json:"logged_details"`
	LoggedRunID        *string  `json:"logged_run_id"`
	MythicLevel        int      `json:"mythic_level"`
	NumChests          int      `json:"num_chests"`
	NumModifiersActive int      `json:"num_modifiers_active"`
	Roster             []Roster `json:"roster"`
	Score              float64  `json:"score"`
	Season             string   `json:"season"`
	Status             string   `json:"status"`
	TimeRemainingMS    int      `json:"time_remaining_ms"`
	WeeklyModifiers    []Affix  `json:"weekly_modifiers"`
}

type Dungeon struct {
	ExpansionID            int    `json:"expansion_id"`
	GroupFinderActivityIDs []int  `json:"group_finder_activity_ids"`
	IconURL                string `json:"icon_url"`
	ID                     int    `json:"id"`
	KeystoneTimerMS        int    `json:"keystone_timer_ms"`
	MapChallengeModeID     int    `json:"map_challenge_mode_id"`
	Name                   string `json:"name"`
	NumBosses              int    `json:"num_bosses"`
	Patch                  string `json:"patch"`
	ShortName              string `json:"short_name"`
	Slug                   string `json:"slug"`
	Type                   string `json:"type"`
	WowInstanceID          int    `json:"wowInstanceId"`
}

type Roster struct {
	Character    Character `json:"character"`
	Guild        *Guild    `json:"guild"`
	IsTransfer   bool      `json:"isTransfer"`
	Items        Items     `json:"items"`
	OldCharacter *string   `json:"oldCharacter"`
	Ranks        Ranks     `json:"ranks"`
	Role         string    `json:"role"`
}

type Character struct {
	Class               Class         `json:"class"`
	Faction             string        `json:"faction"`
	ID                  int           `json:"id"`
	Level               int           `json:"level"`
	Name                string        `json:"name"`
	Path                string        `json:"path"`
	PersonaID           int           `json:"persona_id"`
	Race                Race          `json:"race"`
	Realm               Realm         `json:"realm"`
	RecruitmentProfiles []string      `json:"recruitmentProfiles"`
	Region              Region        `json:"region"`
	Spec                Spec          `json:"spec"`
	TalentLoadout       TalentLoadout `json:"talentLoadout"`
}

type Class struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Race struct {
	Faction string `json:"faction"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
}

type Realm struct {
	AltName             *string `json:"altName"`
	AltSlug             string  `json:"altSlug"`
	ConnectedRealmID    int     `json:"connectedRealmId"`
	ID                  int     `json:"id"`
	IsConnected         bool    `json:"isConnected"`
	Locale              string  `json:"locale"`
	Name                string  `json:"name"`
	RealmType           string  `json:"realmType"`
	Slug                string  `json:"slug"`
	WowConnectedRealmID int     `json:"wowConnectedRealmId"`
	WowRealmID          int     `json:"wowRealmId"`
}

type Region struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Slug      string `json:"slug"`
}

type Spec struct {
	ClassID int    `json:"class_id"`
	ID      int    `json:"id"`
	IsMelee bool   `json:"is_melee"`
	Name    string `json:"name"`
	Patch   string `json:"patch"`
	Role    string `json:"role"`
	Slug    string `json:"slug"`
}

type TalentLoadout struct {
	HeroSubTreeID int    `json:"heroSubTreeId"`
	LoadoutText   string `json:"loadoutText"`
	SpecID        int    `json:"specId"`
}

/*
type Talent struct {
	EntryIndex int  `json:"entryIndex"`
	Node       Node `json:"node"`
	Rank       int  `json:"rank"`
}

type Node struct {
	Col       int     `json:"col"`
	Entries   []Entry `json:"entries"`
	ID        int     `json:"id"`
	Important bool    `json:"important"`
	PosX      int     `json:"posX"`
	PosY      int     `json:"posY"`
	Row       int     `json:"row"`
	SubTreeID int     `json:"subTreeId"`
	TreeID    int     `json:"treeId"`
	Type      int     `json:"type"`
}

type Entry struct {
	ID                int   `json:"id"`
	MaxRanks          int   `json:"maxRanks"`
	Spell             Spell `json:"spell"`
	TraitDefinitionID int   `json:"traitDefinitionId"`
	TraitSubTreeID    int   `json:"traitSubTreeId"`
	Type              int   `json:"type"`
}

type Spell struct {
	HasCooldown bool    `json:"hasCooldown"`
	Icon        string  `json:"icon"`
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Rank        *string `json:"rank"`
	School      int     `json:"school"`
}
*/

type Guild struct {
	Faction string `json:"faction"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Realm   Realm  `json:"realm"`
	Region  Region `json:"region"`
}

type Items struct {
	ItemLevelEquipped float64                 `json:"item_level_equipped"`
	ItemLevelTotal    int                     `json:"item_level_total"`
	Items             map[string]EquippedItem `json:"items"`
	UpdatedAt         string                  `json:"updated_at"`
}

type EquippedItem struct {
	Bonuses     []int  `json:"bonuses"`
	Enchant     *int   `json:"enchant,omitempty"`
	Gems        []int  `json:"gems,omitempty"`
	Icon        string `json:"icon"`
	IsLegendary bool   `json:"is_legendary"`
	ItemID      int    `json:"item_id"`
	ItemLevel   int    `json:"item_level"`
	ItemQuality int    `json:"item_quality"`
	Name        string `json:"name"`
	Tier        string `json:"tier,omitempty"`
}

type Ranks struct {
	Realm  int     `json:"realm"`
	Region int     `json:"region"`
	Score  float64 `json:"score"`
	World  int     `json:"world"`
}

type Affix struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
}
