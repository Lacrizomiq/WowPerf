package wrapper

import (
	"fmt"
	"strings"
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

var treeIDs = map[string]int{
	"Hunter":       774,
	"Shaman":       786,
	"Druid":        793,
	"Evoker":       872,
	"Warrior":      850,
	"Death Knight": 750,
	"Paladin":      790,
	"Priest":       795,
	"Monk":         1000,
	"Mage":         658,
	"Rogue":        852,
	"Demon Hunter": 854,
	"Warlock":      720,
}

var specIDs = map[string]map[string]int{
	"Hunter": {
		"Beast Mastery": 253,
		"Marksmanship":  254,
		"Survival":      255,
	},
	"Shaman": {
		"Elemental":   262,
		"Enhancement": 263,
		"Restoration": 264,
	},
	"Druid": {
		"Balance":     102,
		"Feral":       103,
		"Guardian":    104,
		"Restoration": 105,
	},
	"Evoker": {
		"Devastation":  1467,
		"Preservation": 1468,
		"Augmentation": 1473,
	},
	"Warrior": {
		"Arms":       71,
		"Fury":       72,
		"Protection": 73,
	},
	"Death Knight": {
		"Blood":  250,
		"Frost":  251,
		"Unholy": 252,
	},
	"Paladin": {
		"Holy":        65,
		"Protection":  66,
		"Retribution": 70,
	},
	"Priest": {
		"Discipline": 256,
		"Holy":       257,
		"Shadow":     258,
	},
	"Monk": {
		"Brewmaster": 268,
		"Mistweaver": 270,
		"Windwalker": 269,
	},
	"Mage": {
		"Arcane": 62,
		"Fire":   63,
		"Frost":  64,
	},
	"Rogue": {
		"Assassination": 259,
		"Outlaw":        260,
		"Subtlety":      261,
	},
	"Demon Hunter": {
		"Havoc":     577,
		"Vengeance": 581,
	},
	"Warlock": {
		"Affliction":  265,
		"Demonology":  266,
		"Destruction": 267,
	},
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

	// region
	if links, ok := characterData["_links"].(map[string]interface{}); ok {
		if self, ok := links["self"].(map[string]interface{}); ok {
			if href, ok := self["href"].(string); ok {
				// Extraire la région de l'URL
				parts := strings.Split(href, ".")
				if len(parts) > 1 {
					profile.Region = strings.ToLower(parts[0][8:])
				}
			}
		}
	}

	// media
	if assets, ok := mediaData["assets"].([]interface{}); ok {
		for _, asset := range assets {
			assetMap := asset.(map[string]interface{})
			key := assetMap["key"].(string)
			value := assetMap["value"].(string)

			switch key {
			case "avatar":
				profile.AvatarURL = value
			case "inset":
				profile.InsetAvatarURL = value
			case "main-raw":
				profile.MainRawUrl = value
			}
		}
	}

	if profile.Region != "" && profile.Realm != "" && profile.Name != "" {
		profile.ProfileURL = fmt.Sprintf("https://worldofwarcraft.com/en-gb/character/%s/%s/%s",
			strings.ToLower(profile.Region),
			strings.ToLower(profile.Realm),
			strings.ToLower(profile.Name))
	}

	profile.TreeID = treeIDs[profile.Class]

	if specMap, ok := specIDs[profile.Class]; ok {
		if specID, ok := specMap[profile.ActiveSpecName]; ok {
			profile.SpecID = specID
		} else {
			return nil, fmt.Errorf("unknown specialization %s for class %s", profile.ActiveSpecName, profile.Class)
		}
	} else {
		return nil, fmt.Errorf("unknown class %s", profile.Class)
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
