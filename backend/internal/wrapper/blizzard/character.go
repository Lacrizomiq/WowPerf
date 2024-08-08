package wrapper

import (
	"wowperf/internal/models"
)

const (
	RoleTank   = "Tank"
	RoleHealer = "Healer"
	RoleDPS    = "DPS"
)

var specRoles = map[string]string{
	"Blood":         RoleTank,
	"Frost":         RoleDPS,
	"Unholy":        RoleDPS,
	"Havoc":         RoleDPS,
	"Vengeance":     RoleTank,
	"Balance":       RoleDPS,
	"Feral":         RoleDPS,
	"Guardian":      RoleTank,
	"Restoration":   RoleHealer,
	"Devastation":   RoleDPS,
	"Preservation":  RoleHealer,
	"BeastMastery":  RoleDPS,
	"Marksmanship":  RoleDPS,
	"Survival":      RoleDPS,
	"Arcane":        RoleDPS,
	"Fire":          RoleDPS,
	"Brewmaster":    RoleTank,
	"Mistweaver":    RoleHealer,
	"Windwalker":    RoleDPS,
	"HolyPaladin":   RoleHealer,
	"Protection":    RoleTank,
	"Retribution":   RoleDPS,
	"Discipline":    RoleHealer,
	"HolyPriest":    RoleHealer,
	"Shadow":        RoleDPS,
	"Assassination": RoleDPS,
	"Outlaw":        RoleDPS,
	"Subtlety":      RoleDPS,
	"Elemental":     RoleDPS,
	"Enhancement":   RoleDPS,
	"Affliction":    RoleDPS,
	"Demonology":    RoleDPS,
	"Destruction":   RoleDPS,
	"Arms":          RoleDPS,
	"Fury":          RoleDPS,
	"ProtWarrior":   RoleTank,
}

// TransformCharacterInfo transforms the character data from the Blizzard API into an easier to use CharacterProfile struct
func TransformCharacterInfo(characterData map[string]interface{}, mediaData map[string]interface{}) (*models.CharacterProfile, error) {
	profile := &models.CharacterProfile{}

	// basic profile info
	profile.Name = characterData["name"].(string)
	profile.Race = characterData["race"].(map[string]interface{})["name"].(string)
	profile.Class = characterData["character_class"].(map[string]interface{})["name"].(string)
	profile.ActiveSpecName = characterData["active_spec"].(map[string]interface{})["name"].(string)
	profile.ActiveSpecRole = getRoleFromSpec(profile.ActiveSpecName)
	profile.Gender = characterData["gender"].(map[string]interface{})["name"].(string)
	profile.Faction = characterData["faction"].(map[string]interface{})["name"].(string)
	profile.AchievementPoints = int(characterData["achievement_points"].(float64))
	profile.Realm = characterData["realm"].(map[string]interface{})["name"].(string)

	// media
	if assets, ok := mediaData["assets"].([]interface{}); ok {
		for _, asset := range assets {
			assetMap := asset.(map[string]interface{})
			key := assetMap["key"].(string)
			value := assetMap["value"].(string)

			switch key {
			case "avatar":
				profile.AvatarURL = value
			case "inset_avatar":
				profile.InsetAvatarURL = value
			case "main":
				profile.MainRawUrl = value
			}
		}
	}
	return profile, nil
}

func getRoleFromSpec(specName string) string {
	role, ok := specRoles[specName]
	if !ok {
		return ""
	}
	return role
}
