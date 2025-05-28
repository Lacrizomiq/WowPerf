// models/raiderio/mythicplus_runs/mythicplus_runs_raw.go
package raiderioMythicPlusRunsModels

import (
	"fmt"
	"time"
)

// Run contient toutes les informations d'une run mythic+
type Run struct {
	KeystoneTeamID     int64            `json:"keystone_team_id"`
	Score              float64          `json:"score"`
	Season             string           `json:"season"`
	Status             string           `json:"status"`
	Dungeon            DungeonInfo      `json:"dungeon"`
	KeystoneRunID      int64            `json:"keystone_run_id"`
	MythicLevel        int              `json:"mythic_level"`
	ClearTimeMs        int64            `json:"clear_time_ms"`
	KeystoneTimeMs     int64            `json:"keystone_time_ms"`
	CompletedAt        time.Time        `json:"completed_at"`
	NumChests          int              `json:"num_chests"`
	TimeRemainingMs    int64            `json:"time_remaining_ms"`
	LoggedRunID        *int64           `json:"logged_run_id"`
	Videos             []interface{}    `json:"videos"` // Array vide dans l'exemple
	WeeklyModifiers    []WeeklyModifier `json:"weekly_modifiers"`
	NumModifiersActive int              `json:"num_modifiers_active"`
	Faction            string           `json:"faction"`
	DeletedAt          *time.Time       `json:"deleted_at"`
	KeystonePlatoonID  *int64           `json:"keystone_platoon_id"`
	Platoon            *interface{}     `json:"platoon"` // null dans l'exemple
	Roster             []RosterMember   `json:"roster"`
}

// DungeonInfo contient les informations du donjon
type DungeonInfo struct {
	Type                   string `json:"type"`
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	ShortName              string `json:"short_name"`
	Slug                   string `json:"slug"`
	ExpansionID            int    `json:"expansion_id"`
	IconURL                string `json:"icon_url"`
	Patch                  string `json:"patch"`
	WowInstanceID          int    `json:"wowInstanceId"`
	MapChallengeModeID     int    `json:"map_challenge_mode_id"`
	KeystoneTimerMs        int64  `json:"keystone_timer_ms"`
	NumBosses              int    `json:"num_bosses"`
	GroupFinderActivityIDs []int  `json:"group_finder_activity_ids"`
}

// WeeklyModifier représente un modificateur hebdomadaire (affix)
type WeeklyModifier struct {
	ID          int    `json:"id"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// RosterMember représente un membre du groupe
type RosterMember struct {
	Character    Character  `json:"character"`
	OldCharacter *Character `json:"oldCharacter"`
	IsTransfer   bool       `json:"isTransfer"`
	Role         string     `json:"role"` // tank, healer, dps
}

// Character contient les informations d'un personnage
type Character struct {
	ID                  int64         `json:"id"`
	PersonaID           int64         `json:"persona_id"`
	Name                string        `json:"name"`
	Class               ClassInfo     `json:"class"`
	Race                RaceInfo      `json:"race"`
	Faction             string        `json:"faction"`
	Level               int           `json:"level"`
	Spec                SpecInfo      `json:"spec"`
	Path                string        `json:"path"`
	Realm               RealmInfo     `json:"realm"`
	Region              RegionInfo    `json:"region"`
	Stream              *StreamInfo   `json:"stream"`
	RecruitmentProfiles []interface{} `json:"recruitmentProfiles"` // Array vide dans l'exemple
}

// ClassInfo contient les informations de classe
type ClassInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// RaceInfo contient les informations de race
type RaceInfo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Faction string `json:"faction"`
}

// SpecInfo contient les informations de spécialisation
type SpecInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// RealmInfo contient les informations du royaume
type RealmInfo struct {
	ID                  int     `json:"id"`
	ConnectedRealmID    int     `json:"connectedRealmId"`
	WowRealmID          int     `json:"wowRealmId"`
	WowConnectedRealmID int     `json:"wowConnectedRealmId"`
	Name                string  `json:"name"`
	AltName             *string `json:"altName"`
	Slug                string  `json:"slug"`
	AltSlug             string  `json:"altSlug"`
	Locale              string  `json:"locale"`
	IsConnected         bool    `json:"isConnected"`
	RealmType           string  `json:"realmType"`
}

// RegionInfo contient les informations de région
type RegionInfo struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	ShortName string `json:"short_name"`
}

// StreamInfo contient les informations de stream Twitch (optionnel)
type StreamInfo struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	UserID       string        `json:"user_id"`
	GameID       string        `json:"game_id"`
	Type         string        `json:"type"`
	Title        string        `json:"title"`
	CommunityIDs []interface{} `json:"community_ids"`
	ViewerCount  int           `json:"viewer_count"`
	StartedAt    time.Time     `json:"started_at"`
	Language     string        `json:"language"`
	ThumbnailURL string        `json:"thumbnail_url"`
}

// Méthodes utilitaires pour faciliter l'utilisation

// GetClearTimeFormatted retourne le temps de clear formaté en minutes:secondes
func (r *Run) GetClearTimeFormatted() string {
	totalSeconds := r.ClearTimeMs / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// GetTimeRemainingFormatted retourne le temps restant formaté
func (r *Run) GetTimeRemainingFormatted() string {
	if r.TimeRemainingMs <= 0 {
		return "0:00"
	}
	totalSeconds := r.TimeRemainingMs / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// IsKeystoneInTime indique si la clé a été faite dans les temps
func (r *Run) IsKeystoneInTime() bool {
	return r.TimeRemainingMs > 0
}

// GetMainAffixes retourne les affixes principaux (exclut les affixes saisonniers)
func (r *Run) GetMainAffixes() []WeeklyModifier {
	var mainAffixes []WeeklyModifier
	for _, modifier := range r.WeeklyModifiers {
		// Les affixes principaux ont généralement des IDs spécifiques
		// 9 = Tyrannical, 10 = Fortified, etc.
		if modifier.ID == 9 || modifier.ID == 10 {
			mainAffixes = append(mainAffixes, modifier)
		}
	}
	return mainAffixes
}

// GetRoleMembers retourne les membres par rôle
func (r *Run) GetRoleMembers(role string) []RosterMember {
	var members []RosterMember
	for _, member := range r.Roster {
		if member.Role == role {
			members = append(members, member)
		}
	}
	return members
}

// GetTank retourne le tank du groupe
func (r *Run) GetTank() *RosterMember {
	for _, member := range r.Roster {
		if member.Role == "tank" {
			return &member
		}
	}
	return nil
}

// GetHealer retourne le healer du groupe
func (r *Run) GetHealer() *RosterMember {
	for _, member := range r.Roster {
		if member.Role == "healer" {
			return &member
		}
	}
	return nil
}

// GetDPS retourne les DPS du groupe
func (r *Run) GetDPS() []RosterMember {
	return r.GetRoleMembers("dps")
}
