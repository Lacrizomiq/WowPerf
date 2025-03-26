package protectedProfile

// CharacterSelection represent the selection of a character by the user
type CharacterSelection struct {
	CharacterID int64  `json:"character_id"`
	Name        string `json:"name"`
	Realm       string `json:"realm"`
	Region      string `json:"region"`
	Class       string `json:"class"`
	Race        string `json:"race"`
	Level       int    `json:"level"`
	Faction     string `json:"faction"`
	IsFavorite  bool   `json:"is_favorite"`
}

// CharacterBasicInfo contains basic informations about a character
type CharacterBasicInfo struct {
	CharacterID int64  `json:"character_id"`
	Name        string `json:"name"`
	Realm       string `json:"realm"`
	Region      string `json:"region"`
	Class       string `json:"class"`
	Race        string `json:"race"`
	Level       int    `json:"level"`
	Faction     string `json:"faction"`
}

// CharacterSyncResult represent the result of a sync
type CharacterSyncResult struct {
	Success     bool   `json:"success"`
	CharacterID int64  `json:"character_id"`
	Message     string `json:"message"`
}

// AccountCharactersResponse is the blizzard API struct response
type AccountCharactersResponse struct {
	WowAccounts []WowAccount `json:"wow_accounts"`
}

type WowAccount struct {
	ID         int                `json:"id"`
	Characters []AccountCharacter `json:"characters"`
}

type AccountCharacter struct {
	ID            int64            `json:"id"`
	Name          string           `json:"name"`
	Realm         Realm            `json:"realm"`
	PlayableClass PlayableClass    `json:"playable_class"`
	PlayableRace  PlayableRace     `json:"playable_race"`
	Level         int              `json:"level"`
	Faction       CharacterFaction `json:"faction"`
}

type Realm struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type PlayableClass struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PlayableRace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CharacterFaction struct {
	Type string `json:"type"`
	Name string `json:"name"`
}
