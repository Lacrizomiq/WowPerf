package wrapper

import (
	mythicrundetails "wowperf/internal/models/raiderio/mythicrundetails"
)

// TransformMythicPlusRun transforms the raw data from the Raider.io Mythic Plus Runs Detail API into a MythicPlusRun struct
// It helps to clean up the data and remove any unnecessary fields or data that is not needed for the application
func TransformMythicPlusRun(rawData map[string]interface{}) (*mythicrundetails.MythicPlusRun, error) {
	run := &mythicrundetails.MythicPlusRun{}

	// Basic fields
	if v, ok := rawData["clear_time_ms"].(float64); ok {
		run.ClearTimeMS = int(v)
	}
	run.CompletedAt, _ = rawData["completed_at"].(string)
	if v, ok := rawData["deleted_at"].(string); ok {
		run.DeletedAt = &v
	}
	run.Faction, _ = rawData["faction"].(string)
	if v, ok := rawData["keystone_run_id"].(float64); ok {
		run.KeystoneRunID = int(v)
	}
	if v, ok := rawData["keystone_team_id"].(float64); ok {
		run.KeystoneTeamID = int(v)
	}
	if v, ok := rawData["keystone_time_ms"].(float64); ok {
		run.KeystoneTimeMS = int(v)
	}
	if v, ok := rawData["logged_details"].(string); ok {
		run.LoggedDetails = &v
	}
	if v, ok := rawData["logged_run_id"].(string); ok {
		run.LoggedRunID = &v
	}
	if v, ok := rawData["mythic_level"].(float64); ok {
		run.MythicLevel = int(v)
	}
	if v, ok := rawData["num_chests"].(float64); ok {
		run.NumChests = int(v)
	}
	if v, ok := rawData["num_modifiers_active"].(float64); ok {
		run.NumModifiersActive = int(v)
	}
	if v, ok := rawData["score"].(float64); ok {
		run.Score = v
	}
	run.Season, _ = rawData["season"].(string)
	run.Status, _ = rawData["status"].(string)
	if v, ok := rawData["time_remaining_ms"].(float64); ok {
		run.TimeRemainingMS = int(v)
	}

	// Dungeon
	if dungeonData, ok := rawData["dungeon"].(map[string]interface{}); ok {
		run.Dungeon = transformDungeon(dungeonData)
	}

	// Roster
	if rosterData, ok := rawData["roster"].([]interface{}); ok {
		run.Roster = transformRoster(rosterData)
	}

	// Weekly Modifiers
	if modifiersData, ok := rawData["weekly_modifiers"].([]interface{}); ok {
		run.WeeklyModifiers = transformAffixes(modifiersData)
	}

	// Logged Sources
	if sourcesData, ok := rawData["loggedSources"].([]interface{}); ok {
		run.LoggedSources = make([]string, len(sourcesData))
		for i, source := range sourcesData {
			run.LoggedSources[i], _ = source.(string)
		}
	}

	return run, nil
}

func transformDungeon(data map[string]interface{}) mythicrundetails.Dungeon {
	dungeon := mythicrundetails.Dungeon{}

	if v, ok := data["expansion_id"].(float64); ok {
		dungeon.ExpansionID = int(v)
	}
	dungeon.IconURL, _ = data["icon_url"].(string)
	if v, ok := data["id"].(float64); ok {
		dungeon.ID = int(v)
	}
	if v, ok := data["keystone_timer_ms"].(float64); ok {
		dungeon.KeystoneTimerMS = int(v)
	}
	if v, ok := data["map_challenge_mode_id"].(float64); ok {
		dungeon.MapChallengeModeID = int(v)
	}
	dungeon.Name, _ = data["name"].(string)
	if v, ok := data["num_bosses"].(float64); ok {
		dungeon.NumBosses = int(v)
	}
	dungeon.Patch, _ = data["patch"].(string)
	dungeon.ShortName, _ = data["short_name"].(string)
	dungeon.Slug, _ = data["slug"].(string)
	dungeon.Type, _ = data["type"].(string)
	if v, ok := data["wowInstanceId"].(float64); ok {
		dungeon.WowInstanceID = int(v)
	}

	if activityIDs, ok := data["group_finder_activity_ids"].([]interface{}); ok {
		dungeon.GroupFinderActivityIDs = make([]int, len(activityIDs))
		for i, id := range activityIDs {
			if v, ok := id.(float64); ok {
				dungeon.GroupFinderActivityIDs[i] = int(v)
			}
		}
	}

	return dungeon
}

func transformRoster(data []interface{}) []mythicrundetails.Roster {
	roster := make([]mythicrundetails.Roster, len(data))
	for i, memberData := range data {
		if member, ok := memberData.(map[string]interface{}); ok {
			roster[i] = transformRosterMember(member)
		}
	}
	return roster
}

func transformRosterMember(data map[string]interface{}) mythicrundetails.Roster {
	member := mythicrundetails.Roster{}

	if characterData, ok := data["character"].(map[string]interface{}); ok {
		member.Character = transformCharacter(characterData)
	}
	if guildData, ok := data["guild"].(map[string]interface{}); ok {
		guild := transformGuild(guildData)
		member.Guild = &guild
	}
	member.IsTransfer, _ = data["isTransfer"].(bool)
	if itemsData, ok := data["items"].(map[string]interface{}); ok {
		member.Items = transformItems(itemsData)
	}
	if v, ok := data["oldCharacter"].(string); ok {
		member.OldCharacter = &v
	}
	if ranksData, ok := data["ranks"].(map[string]interface{}); ok {
		member.Ranks = transformRanks(ranksData)
	}
	member.Role, _ = data["role"].(string)

	return member
}

func transformCharacter(data map[string]interface{}) mythicrundetails.Character {
	character := mythicrundetails.Character{}

	if classData, ok := data["class"].(map[string]interface{}); ok {
		character.Class = transformClass(classData)
	}
	character.Faction, _ = data["faction"].(string)
	if v, ok := data["id"].(float64); ok {
		character.ID = int(v)
	}
	if v, ok := data["level"].(float64); ok {
		character.Level = int(v)
	}
	character.Name, _ = data["name"].(string)
	character.Path, _ = data["path"].(string)
	if v, ok := data["persona_id"].(float64); ok {
		character.PersonaID = int(v)
	}
	if raceData, ok := data["race"].(map[string]interface{}); ok {
		character.Race = transformRace(raceData)
	}
	if realmData, ok := data["realm"].(map[string]interface{}); ok {
		character.Realm = transformRealm(realmData)
	}
	if regionData, ok := data["region"].(map[string]interface{}); ok {
		character.Region = transformRegion(regionData)
	}
	if specData, ok := data["spec"].(map[string]interface{}); ok {
		character.Spec = transformSpec(specData)
	}
	if talentData, ok := data["talentLoadout"].(map[string]interface{}); ok {
		character.TalentLoadout = transformTalentLoadout(talentData)
	}

	// Handle recruitmentProfiles if needed

	return character
}

func transformClass(data map[string]interface{}) mythicrundetails.Class {
	class := mythicrundetails.Class{}
	if v, ok := data["id"].(float64); ok {
		class.ID = int(v)
	}
	class.Name, _ = data["name"].(string)
	class.Slug, _ = data["slug"].(string)
	return class
}

func transformRace(data map[string]interface{}) mythicrundetails.Race {
	race := mythicrundetails.Race{}
	race.Faction, _ = data["faction"].(string)
	if v, ok := data["id"].(float64); ok {
		race.ID = int(v)
	}
	race.Name, _ = data["name"].(string)
	race.Slug, _ = data["slug"].(string)
	return race
}

func transformRealm(data map[string]interface{}) mythicrundetails.Realm {
	realm := mythicrundetails.Realm{}
	if v, ok := data["altName"].(string); ok {
		realm.AltName = &v
	}
	realm.AltSlug, _ = data["altSlug"].(string)
	if v, ok := data["connectedRealmId"].(float64); ok {
		realm.ConnectedRealmID = int(v)
	}
	if v, ok := data["id"].(float64); ok {
		realm.ID = int(v)
	}
	realm.IsConnected, _ = data["isConnected"].(bool)
	realm.Locale, _ = data["locale"].(string)
	realm.Name, _ = data["name"].(string)
	realm.RealmType, _ = data["realmType"].(string)
	realm.Slug, _ = data["slug"].(string)
	if v, ok := data["wowConnectedRealmId"].(float64); ok {
		realm.WowConnectedRealmID = int(v)
	}
	if v, ok := data["wowRealmId"].(float64); ok {
		realm.WowRealmID = int(v)
	}
	return realm
}

func transformRegion(data map[string]interface{}) mythicrundetails.Region {
	region := mythicrundetails.Region{}
	region.Name, _ = data["name"].(string)
	region.ShortName, _ = data["short_name"].(string)
	region.Slug, _ = data["slug"].(string)
	return region
}

func transformSpec(data map[string]interface{}) mythicrundetails.Spec {
	spec := mythicrundetails.Spec{}
	if v, ok := data["class_id"].(float64); ok {
		spec.ClassID = int(v)
	}
	if v, ok := data["id"].(float64); ok {
		spec.ID = int(v)
	}
	spec.IsMelee, _ = data["is_melee"].(bool)
	spec.Name, _ = data["name"].(string)
	spec.Patch, _ = data["patch"].(string)
	spec.Role, _ = data["role"].(string)
	spec.Slug, _ = data["slug"].(string)
	return spec
}

func transformTalentLoadout(data map[string]interface{}) mythicrundetails.TalentLoadout {
	loadout := mythicrundetails.TalentLoadout{}
	if v, ok := data["heroSubTreeId"].(float64); ok {
		loadout.HeroSubTreeID = int(v)
	}
	loadout.LoadoutText, _ = data["loadoutText"].(string)
	if v, ok := data["specId"].(float64); ok {
		loadout.SpecID = int(v)
	}
	return loadout
}

func transformGuild(data map[string]interface{}) mythicrundetails.Guild {
	guild := mythicrundetails.Guild{}
	guild.Faction, _ = data["faction"].(string)
	if v, ok := data["id"].(float64); ok {
		guild.ID = int(v)
	}
	guild.Name, _ = data["name"].(string)
	guild.Path, _ = data["path"].(string)
	if realmData, ok := data["realm"].(map[string]interface{}); ok {
		guild.Realm = transformRealm(realmData)
	}
	if regionData, ok := data["region"].(map[string]interface{}); ok {
		guild.Region = transformRegion(regionData)
	}
	return guild
}

func transformItems(data map[string]interface{}) mythicrundetails.Items {
	items := mythicrundetails.Items{}
	if v, ok := data["item_level_equipped"].(float64); ok {
		items.ItemLevelEquipped = v
	}
	if v, ok := data["item_level_total"].(float64); ok {
		items.ItemLevelTotal = int(v)
	}
	items.UpdatedAt, _ = data["updated_at"].(string)

	if itemsData, ok := data["items"].(map[string]interface{}); ok {
		items.Items = make(map[string]mythicrundetails.EquippedItem)
		for slot, itemData := range itemsData {
			if item, ok := itemData.(map[string]interface{}); ok {
				items.Items[slot] = transformEquippedItem(item)
			}
		}
	}

	return items
}

func transformEquippedItem(data map[string]interface{}) mythicrundetails.EquippedItem {
	item := mythicrundetails.EquippedItem{}

	if bonusesData, ok := data["bonuses"].([]interface{}); ok {
		item.Bonuses = make([]int, len(bonusesData))
		for i, bonus := range bonusesData {
			if v, ok := bonus.(float64); ok {
				item.Bonuses[i] = int(v)
			}
		}
	}

	if v, ok := data["enchant"].(float64); ok {
		enchant := int(v)
		item.Enchant = &enchant
	}

	if gemsData, ok := data["gems"].([]interface{}); ok {
		item.Gems = make([]int, len(gemsData))
		for i, gem := range gemsData {
			if v, ok := gem.(float64); ok {
				item.Gems[i] = int(v)
			}
		}
	}

	item.Icon, _ = data["icon"].(string)
	item.IsLegendary, _ = data["is_legendary"].(bool)
	if v, ok := data["item_id"].(float64); ok {
		item.ItemID = int(v)
	}
	if v, ok := data["item_level"].(float64); ok {
		item.ItemLevel = int(v)
	}
	if v, ok := data["item_quality"].(float64); ok {
		item.ItemQuality = int(v)
	}
	item.Name, _ = data["name"].(string)
	item.Tier, _ = data["tier"].(string)

	return item
}

func transformRanks(data map[string]interface{}) mythicrundetails.Ranks {
	ranks := mythicrundetails.Ranks{}
	if v, ok := data["realm"].(float64); ok {
		ranks.Realm = int(v)
	}
	if v, ok := data["region"].(float64); ok {
		ranks.Region = int(v)
	}
	if v, ok := data["score"].(float64); ok {
		ranks.Score = v
	}
	if v, ok := data["world"].(float64); ok {
		ranks.World = int(v)
	}
	return ranks
}

func transformAffixes(data []interface{}) []mythicrundetails.Affix {
	affixes := make([]mythicrundetails.Affix, len(data))
	for i, affixData := range data {
		if affix, ok := affixData.(map[string]interface{}); ok {
			affixes[i] = transformAffix(affix)
		}
	}
	return affixes
}

func transformAffix(data map[string]interface{}) mythicrundetails.Affix {
	affix := mythicrundetails.Affix{}
	affix.Description, _ = data["description"].(string)
	affix.Icon, _ = data["icon"].(string)
	if v, ok := data["id"].(float64); ok {
		affix.ID = int(v)
	}
	affix.Name, _ = data["name"].(string)
	affix.Slug, _ = data["slug"].(string)
	return affix
}
