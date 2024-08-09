package wrapper

import (
	"fmt"
	"strings"
	"wowperf/internal/models"
)

func TransformCharacterGear(data map[string]interface{}) (*models.Gear, error) {
	gear := &models.Gear{
		Items: make(map[string]models.Item),
	}

	equippedItems, ok := data["equipped_items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("equipped_items not found or not a slice")
	}

	for _, item := range equippedItems {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		slotInfo, ok := itemMap["slot"].(map[string]interface{})
		if !ok {
			continue
		}

		slotType, ok := slotInfo["type"].(string)
		if !ok {
			continue
		}

		itemInfo, ok := itemMap["item"].(map[string]interface{})
		if !ok {
			continue
		}

		itemID, ok := itemInfo["id"].(float64)
		if !ok {
			continue
		}

		level, ok := itemMap["level"].(map[string]interface{})
		if !ok {
			continue
		}

		itemLevel, ok := level["value"].(float64)
		if !ok {
			continue
		}

		quality, ok := itemMap["quality"].(map[string]interface{})
		if !ok {
			continue
		}

		itemQuality, ok := quality["type"].(string)
		if !ok {
			continue
		}

		name, ok := itemMap["name"].(string)
		if !ok {
			continue
		}

		media, ok := itemMap["media"].(map[string]interface{})
		if !ok {
			continue
		}

		icon := ""
		if mediaKey, ok := media["key"].(map[string]interface{}); ok {
			if href, ok := mediaKey["href"].(string); ok {
				parts := strings.Split(href, "/")
				if len(parts) > 0 {
					icon = parts[len(parts)-1]
				}
			}
		}

		var enchant *int
		if enchantments, ok := itemMap["enchantments"].([]interface{}); ok && len(enchantments) > 0 {
			if enchantment, ok := enchantments[0].(map[string]interface{}); ok {
				if enchantID, ok := enchantment["enchantment_id"].(float64); ok {
					enchantValue := int(enchantID)
					enchant = &enchantValue
				}
			}
		}

		bonusList := []int{}
		if bonuses, ok := itemMap["bonus_list"].([]interface{}); ok {
			for _, bonus := range bonuses {
				if bonusInt, ok := bonus.(float64); ok {
					bonusList = append(bonusList, int(bonusInt))
				}
			}
		}

		gems := []int{}
		if sockets, ok := itemMap["sockets"].([]interface{}); ok {
			for _, socket := range sockets {
				if socketMap, ok := socket.(map[string]interface{}); ok {
					if item, ok := socketMap["item"].(map[string]interface{}); ok {
						if gemID, ok := item["id"].(float64); ok {
							gems = append(gems, int(gemID))
						}
					}
				}
			}
		}

		gear.Items[strings.ToLower(slotType)] = models.Item{
			ItemID:      int(itemID),
			ItemLevel:   itemLevel,
			ItemQuality: getItemQualityInt(itemQuality),
			Icon:        icon,
			Name:        name,
			Enchant:     enchant,
			Gems:        gems,
			Bonuses:     bonusList,
		}
	}

	return gear, nil
}

func getItemQualityInt(quality string) int {
	switch quality {
	case "POOR":
		return 0
	case "COMMON":
		return 1
	case "UNCOMMON":
		return 2
	case "RARE":
		return 3
	case "EPIC":
		return 4
	case "LEGENDARY":
		return 5
	case "ARTIFACT":
		return 6
	case "HEIRLOOM":
		return 7
	default:
		return 1
	}
}
