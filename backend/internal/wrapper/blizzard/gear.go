package wrapper

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"wowperf/internal/models"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/gamedata"
)

func isItemEmpty(item models.Item) bool {
	return item.ItemID == 0 && item.ItemLevel == 0 && item.Name == ""
}

// TransformCharacterGear transforms the gear data from the Blizzard API into an easier to use Gear struct.
// Using a channel to handle the concurrency of the requests.
func TransformCharacterGear(data map[string]interface{}, gameDataService *blizzard.GameDataService, region, namespace, locale string) (*models.Gear, error) {
	gear := &models.Gear{
		Items: make(map[string]models.Item),
	}

	equippedItems, ok := data["equipped_items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("equipped_items not found or not a slice")
	}

	var wg sync.WaitGroup
	itemChan := make(chan struct {
		slotType string
		item     models.Item
	}, len(equippedItems))
	errorChan := make(chan error, len(equippedItems))

	for _, item := range equippedItems {
		wg.Add(1)
		go func(item interface{}) {
			defer wg.Done()
			slotType, transformedItem, err := transformSingleItem(item, gameDataService, region, namespace, locale)
			if err != nil {
				errorChan <- err
				return
			}
			itemChan <- struct {
				slotType string
				item     models.Item
			}{
				slotType: slotType,
				item:     transformedItem,
			}
		}(item)
	}

	go func() {
		wg.Wait()
		close(itemChan)
		close(errorChan)
	}()

	var totalItemLevel float64
	var totalWeight float64
	slotWeights := map[string]float64{
		"head": 1, "neck": 1, "shoulder": 1, "back": 1, "chest": 1, "wrist": 1,
		"hands": 1, "waist": 1, "legs": 1, "feet": 1, "finger_1": 1, "finger_2": 1,
		"trinket_1": 1, "trinket_2": 1, "main_hand": 1, "off_hand": 1,
	}

	hasTwoHandWeapon := false

	for item := range itemChan {
		gear.Items[strings.ToLower(item.slotType)] = item.item

		weight := slotWeights[strings.ToLower(item.slotType)]
		if strings.ToLower(item.slotType) == "main_hand" {
			if item.item.IsTwoHand {
				hasTwoHandWeapon = true
				weight = 2
			}
		}

		totalItemLevel += item.item.ItemLevel * weight
		totalWeight += weight
	}

	if len(errorChan) > 0 {
		return nil, <-errorChan
	}

	if !hasTwoHandWeapon {
		offHandItem, exists := gear.Items["off_hand"]
		if !exists || isItemEmpty(offHandItem) {
			totalWeight--
		}
	}

	if totalWeight > 0 {
		gear.ItemLevelEquipped = totalItemLevel / totalWeight
	}

	return gear, nil
}

// transformSingleItem transforms a single item from the Blizzard API into a struct.
func transformSingleItem(item interface{}, gameDataService *blizzard.GameDataService, region, namespace, locale string) (string, models.Item, error) {
	itemMap, ok := item.(map[string]interface{})
	if !ok {
		return "", models.Item{}, fmt.Errorf("item is not a map")
	}

	slotInfo, ok := itemMap["slot"].(map[string]interface{})
	if !ok {
		return "", models.Item{}, fmt.Errorf("slot info not found")
	}

	slotType, ok := slotInfo["type"].(string)
	if !ok {
		return "", models.Item{}, fmt.Errorf("slot type not found")
	}

	itemInfo, ok := itemMap["item"].(map[string]interface{})
	if !ok {
		return "", models.Item{}, fmt.Errorf("item info not found")
	}

	itemID, ok := itemInfo["id"].(float64)
	if !ok {
		return "", models.Item{}, fmt.Errorf("item ID not found")
	}

	level, ok := itemMap["level"].(map[string]interface{})
	if !ok {
		return "", models.Item{}, fmt.Errorf("level info not found")
	}

	itemLevel, ok := level["value"].(float64)
	if !ok {
		return "", models.Item{}, fmt.Errorf("item level not found")
	}

	quality, ok := itemMap["quality"].(map[string]interface{})
	if !ok {
		return "", models.Item{}, fmt.Errorf("quality info not found")
	}

	itemQuality, ok := quality["type"].(string)
	if !ok {
		return "", models.Item{}, fmt.Errorf("item quality not found")
	}

	name, ok := itemMap["name"].(string)
	if !ok {
		return "", models.Item{}, fmt.Errorf("item name not found")
	}

	iconName, iconURL, err := getItemMedia(itemInfo, gameDataService, region, namespace, locale)
	if err != nil {
		return "", models.Item{}, err
	}

	enchant := getEnchant(itemMap)
	enchantName := getEnchantName(itemMap)
	stats := getItemStats(itemMap)
	bonusList := getBonusList(itemMap)
	gems := getGems(itemMap)

	isTwoHand := false
	if inventory, ok := itemMap["inventory_type"].(map[string]interface{}); ok {
		if inventoryType, ok := inventory["type"].(string); ok {
			isTwoHand = inventoryType == "TWOHWEAPON" || inventoryType == "RANGEDRIGHT"
		}
	}

	transformedItem := models.Item{
		ItemID:      int(itemID),
		ItemLevel:   itemLevel,
		ItemQuality: getItemQualityInt(itemQuality),
		IconName:    iconName,
		IconURL:     iconURL,
		Name:        name,
		Enchant:     enchant,
		EnchantName: enchantName,
		Stats:       stats,
		Gems:        gems,
		Bonuses:     bonusList,
		IsTwoHand:   isTwoHand,
	}

	return slotType, transformedItem, nil
}

// getItemMedia retrieves the media assets for an item.
func getItemMedia(itemInfo map[string]interface{}, gameDataService *blizzard.GameDataService, region, namespace, locale string) (string, string, error) {
	if itemID, ok := itemInfo["id"].(float64); ok {
		mediaData, err := gamedata.GetItemMedia(gameDataService, int(itemID), region, "static-"+region, locale)
		if err == nil {
			if assets, ok := mediaData["assets"].([]interface{}); ok && len(assets) > 0 {
				if asset, ok := assets[0].(map[string]interface{}); ok {
					if value, ok := asset["value"].(string); ok {
						iconURL := value
						parts := strings.Split(value, "/")
						if len(parts) > 0 {
							iconName := strings.TrimSuffix(parts[len(parts)-1], ".jpg")
							return iconName, iconURL, nil
						}
					}
				}
			}
		}
		return "", "", err
	}
	return "", "", fmt.Errorf("item ID not found for media")
}

// getEnchant returns the enchantment ID for an item, if any.
func getEnchant(itemMap map[string]interface{}) *int {
	if enchantments, ok := itemMap["enchantments"].([]interface{}); ok && len(enchantments) > 0 {
		if enchantment, ok := enchantments[0].(map[string]interface{}); ok {
			if enchantID, ok := enchantment["enchantment_id"].(float64); ok {
				enchantValue := int(enchantID)
				return &enchantValue
			}
		}
	}
	return nil
}

// getBonusList returns a list of bonus IDs for an item, if any.
func getBonusList(itemMap map[string]interface{}) []int {
	bonusList := []int{}
	if bonuses, ok := itemMap["bonus_list"].([]interface{}); ok {
		for _, bonus := range bonuses {
			if bonusInt, ok := bonus.(float64); ok {
				bonusList = append(bonusList, int(bonusInt))
			}
		}
	}
	return bonusList
}

// getGemList returns a list of gem IDs for an item, if any.
func getGems(itemMap map[string]interface{}) []int {
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
	return gems
}

// getItemQualityInt converts a string representation of an item quality to an integer.
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

// getEnchantName retrieves the display string for an enchantment.
func getEnchantName(itemMap map[string]interface{}) string {
	if enchantments, ok := itemMap["enchantments"].([]interface{}); ok && len(enchantments) > 0 {
		if enchantment, ok := enchantments[0].(map[string]interface{}); ok {
			if displayString, ok := enchantment["display_string"].(string); ok {
				re := regexp.MustCompile(`\+(\d+)\s+([A-Za-z\s]+)`)
				matches := re.FindStringSubmatch(displayString)
				if len(matches) == 3 {
					return "+" + matches[1] + " " + strings.TrimSpace(matches[2])
				}
			}
		}
	}
	return ""
}

// getItemStats retrieves the stats for an item.
func getItemStats(itemMap map[string]interface{}) []models.ItemStat {
	var stats []models.ItemStat
	if statsData, ok := itemMap["stats"].([]interface{}); ok {
		for _, statData := range statsData {
			if stat, ok := statData.(map[string]interface{}); ok {
				if isNegated, ok := stat["is_negated"].(bool); ok && isNegated {
					continue
				}
				if statType, ok := stat["type"].(map[string]interface{}); ok {
					if statName, ok := statType["name"].(string); ok {
						if value, ok := stat["value"].(float64); ok {
							stats = append(stats, models.ItemStat{
								Type:  statName,
								Value: int(value),
							})
						}
					}
				}
			}
		}
	}
	return stats
}
