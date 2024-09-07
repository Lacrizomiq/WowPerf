package wrapper

import (
	"fmt"
	raidsEncounter "wowperf/internal/models/raids"
)

// TransformRaidData transforms the raid encouter data from the Blizzard API into an easier to use struct
func TransformRaidData(rawData interface{}) (raidsEncounter.ExpansionRaids, error) {
	result := raidsEncounter.ExpansionRaids{
		Expansions: make([]raidsEncounter.ExpansionWithRaids, 0),
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		return result, fmt.Errorf("invalid data format")
	}

	expansions, ok := data["expansions"].([]interface{})
	if !ok {
		return result, fmt.Errorf("expansions data not found")
	}

	for _, exp := range expansions {
		expData, ok := exp.(map[string]interface{})
		if !ok {
			continue
		}

		expansion, ok := expData["expansion"].(map[string]interface{})
		if !ok {
			continue
		}

		expID := int(expansion["id"].(float64))
		expName := expansion["name"].(string)

		if expID != 503 && expID != 505 {
			continue
		}

		expWithRaids := raidsEncounter.ExpansionWithRaids{
			ID:    expID,
			Name:  expName,
			Raids: make([]raidsEncounter.Raids, 0),
		}

		instances, ok := expData["instances"].([]interface{})
		if !ok {
			continue
		}

		for _, inst := range instances {
			instance, ok := inst.(map[string]interface{})
			if !ok {
				continue
			}
			raid := transformRaid(instance)
			expWithRaids.Raids = append(expWithRaids.Raids, raid)
		}

		result.Expansions = append(result.Expansions, expWithRaids)
	}

	return result, nil
}

// transformRaid transforms the raid data from the Blizzard API into an easier to use struct
func transformRaid(instanceData map[string]interface{}) raidsEncounter.Raids {
	instance := instanceData["instance"].(map[string]interface{})
	raid := raidsEncounter.Raids{
		ID:   int(instance["id"].(float64)),
		Name: instance["name"].(string),
	}

	modes := instanceData["modes"].([]interface{})
	for _, m := range modes {
		mode := m.(map[string]interface{})
		raidMode := transformMode(mode)
		raid.Modes = append(raid.Modes, raidMode)
	}

	return raid
}

// transformMode transforms the mode data from the Blizzard API into an easier to use struct
func transformMode(modeData map[string]interface{}) raidsEncounter.Mode {
	mode := raidsEncounter.Mode{}

	if difficulty, ok := modeData["difficulty"].(map[string]interface{}); ok {
		if name, ok := difficulty["name"].(string); ok {
			mode.Difficulty = name
		}
	}

	if progress, ok := modeData["progress"].(map[string]interface{}); ok {
		if completedCount, ok := progress["completed_count"].(float64); ok {
			mode.Progress.CompletedCount = int(completedCount)
		}
		if totalCount, ok := progress["total_count"].(float64); ok {
			mode.Progress.TotalCount = int(totalCount)
		}
		if encounters, ok := progress["encounters"].([]interface{}); ok {
			for _, e := range encounters {
				if enc, ok := e.(map[string]interface{}); ok {
					encounter := transformEncounter(enc)
					mode.Progress.Encounters = append(mode.Progress.Encounters, encounter)
				}
			}
		}
	}

	if status, ok := modeData["status"].(map[string]interface{}); ok {
		if statusType, ok := status["type"].(string); ok {
			mode.Status = statusType
		}
	}

	return mode
}

// transformEncounter transforms the encounter data from the Blizzard API into an easier to use struct
func transformEncounter(encData map[string]interface{}) raidsEncounter.EncounterProgress {
	encounter := raidsEncounter.EncounterProgress{}

	if enc, ok := encData["encounter"].(map[string]interface{}); ok {
		if id, ok := enc["id"].(float64); ok {
			encounter.ID = int(id)
		}
		if name, ok := enc["name"].(string); ok {
			encounter.Name = name
		}
	}

	if completedCount, ok := encData["completed_count"].(float64); ok {
		encounter.CompletedCount = int(completedCount)
	}

	if lastKillTimestamp, ok := encData["last_kill_timestamp"].(float64); ok {
		encounter.LastKillTimestamp = int64(lastKillTimestamp)
	}

	return encounter
}

// GetRaidsByExpansionID gets the raids by expansion ID
func GetRaidsByExpansionID(data raidsEncounter.ExpansionRaids, expansionID int) []raidsEncounter.Raids {
	for _, exp := range data.Expansions {
		if exp.ID == expansionID {
			return exp.Raids
		}
	}
	return nil
}
