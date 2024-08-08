# Blizzard API

## Authentication

Blizzard API uses OAuth 2.0 to authenticate requests. The client ID and client secret are required to obtain an access token.

The access token is used to authenticate requests to the Blizzard API. The token is obtained using the OAuth 2.0 authorization code flow.

## Endpoints

### Character Profile

Retrieves a character's profile information, including name, realm, and class.

**Endpoint**: `/blizzard/characters/{realmSlug}/{characterName}`
`/blizzard/characters/{realmSlug}/{characterName}?region={region}&namespace={namespace}&locale={locale}`

**Parameters**:

- `region`: The region of the character's realm. Possible values: `us`, `eu`, `kr`, `tw`, `cn`.
- `realmSlug`: The slug of the character's realm.
- `characterName`: The name of the character.
- `namespace`: The namespace of the character's profile.
- `locale`: The locale of the character's profile. Possible values: `en_US`, `enGB`, `deDE`, `frFR`, `koKR`, `esES`, `zhCN`, `zhTW`, `ptBR`, `ruRU`, `itIT`, `jaJP`, `plPL`.

**Response**:

```json
{
  "_links": {
    "self": {
      "href": "string"
    }
  },
  "achievement_points": "number",
  "achievements": {
    "href": "string"
  },
  "achievements_statistics": {
    "href": "string"
  },
  "active_spec": {
    "id": "number",
    "key": {
      "href": "string"
    },
    "name": "string"
  },
  "active_title": {
    "display_string": "string",
    "id": "number",
    "key": {
      "href": "string"
    },
    "name": "string"
  },
  "appearance": {
    "href": "string"
  },
  "average_item_level": "number",
  "character_class": {
    "id": "number",
    "key": {
      "href": "string"
    },
    "name": "string"
  },
  "collections": {
    "href": "string"
  },
  "covenant_progress": {
    "chosen_covenant": {
      "id": "number",
      "key": {
        "href": "string"
      },
      "name": "string"
    },
    "renown_level": "number",
    "soulbinds": {
      "href": "string"
    }
  },
  "encounters": {
    "href": "string"
  },
  "equipment": {
    "href": "string"
  },
  "equipped_item_level": "number",
  "experience": "number",
  "faction": {
    "name": "string",
    "type": "string"
  },
  "gender": {
    "name": "string",
    "type": "string"
  },
  "id": "number",
  "last_login_timestamp": "number",
  "level": "number",
  "media": {
    "href": "string"
  },
  "mythic_keystone_profile": {
    "href": "string"
  },
  "name": "string",
  "name_search": "string",
  "professions": {
    "href": "string"
  },
  "pvp_summary": {
    "href": "string"
  },
  "quests": {
    "href": "string"
  },
  "race": {
    "id": "number",
    "key": {
      "href": "string"
    },
    "name": "string"
  },
  "realm": {
    "id": "number",
    "key": {
      "href": "string"
    },
    "name": "string",
    "slug": "string"
  },
  "reputations": {
    "href": "string"
  },
  "specializations": {
    "href": "string"
  },
  "statistics": {
    "href": "string"
  },
  "titles": {
    "href": "string"
  }
}
```

### Mythic Keystone Profile

Retrieves a character's mythic keystone profile information, including seasons, tiers, and keystone upgrades.

**Endpoint**: `/blizzard/characters/{realmSlug}/{characterName}/mythic-keystone-profile`
`/blizzard/characters/{realmSlug}/{characterName}/mythic-keystone-profile?region={region}&namespace={namespace}&locale={locale}`

**Parameters**:

- `region`: The region of the character's realm. Possible values: `us`, `eu`, `kr`, `tw`, `cn`.
- `realmSlug`: The slug of the character's realm.
- `characterName`: The name of the character.
- `namespace`: The namespace of the character's profile.
- `locale`: The locale of the character's profile. Possible values: `enUS`, `enGB`, `deDE`, `frFR`, `koKR`, `esES`, `zhCN`, `zhTW`, `ptBR`, `ruRU`, `itIT`, `jaJP`, `plPL`.

**Response**:

```json
{
  "_links": {
    "self": {
      "href": "string"
    }
  },
  "character": {
    "id": "number",
    "key": {
      "href": "string"
    },
    "name": "string",
    "realm": {
      "id": "number",
      "key": {
        "href": "string"
      },
      "name": "string",
      "slug": "string"
    }
  },
  "current_mythic_rating": {
    "color": {
      "a": "number",
      "b": "number",
      "g": "number",
      "r": "number"
    },
    "rating": "number"
  },
  "current_period": {
    "period": {
      "id": "number",
      "key": {
        "href": "string"
      }
    }
  },
  "seasons": [
    {
      "id": "number",
      "key": {
        "href": "string"
      }
    }
  ]
}
```

### Mythic Keystone Season Details

Retrieves a character's mythic keystone season details, including seasons, tiers, and keystone upgrades.

**Endpoint**: `/blizzard/characters/{realmSlug}/{characterName}/mythic-keystone-profile/season/{seasonId}`
`/blizzard/characters/{realmSlug}/{characterName}/mythic-keystone-profile/season/{seasonId}?region={region}&namespace={namespace}&locale={locale}`

**Parameters**:

- `region`: The region of the character's realm. Possible values: `us`, `eu`, `kr`, `tw`, `cn`.
- `realmSlug`: The slug of the character's realm.
- `characterName`: The name of the character.
- `seasonId`: The ID of the season to retrieve details for.
- `namespace`: The namespace of the character's profile.
- `locale`: The locale of the character's profile. Possible values: `enUS`, `enGB`, `deDE`, `frFR`, `koKR`, `esES`, `zhCN`, `zhTW`, `ptBR`, `ruRU`, `itIT`, `jaJP`, `plPL`.

**Response**:

```json
{
  "_links": {
    "self": {
      "href": "string"
    }
  },
  "character": {
    "id": "number",
    "key": {
      "href": "string"
    },
    "name": "string",
    "realm": {
      "id": "number",
      "key": {
        "href": "string"
      },
      "name": "string",
      "slug": "string"
    }
  },
  "current_mythic_rating": {
    "color": {
      "a": "number",
      "b": "number",
      "g": "number",
      "r": "number"
    },
    "rating": "number"
  },
  "current_period": {
    "period": {
      "id": "number",
      "key": {
        "href": "string"
      }
    }
  },
  "seasons": [
    {
      "id": "number",
      "key": {
        "href": "string"
      }
    }
  ]
}
```

### Equipment

Retrieves a character's equipment information, including items, gems, and mounts.

**Endpoint**: `/blizzard/characters/{realmSlug}/{characterName}/equipment`
`blizzard/characters/{realmSlug}/{characterName}/equipment?region={region}&namespace={namespace}&locale={locale}`

**Parameters**:

- `region`: The region of the character's realm. Possible values: `us`, `eu`, `kr`, `tw`, `cn`.
- `realmSlug`: The slug of the character's realm.
- `characterName`: The name of the character.
- `namespace`: The namespace of the character's profile.
- `locale`: The locale of the character's profile. Possible values: `en_US`, `enGB`, `deDE`, `frFR`, `koKR`, `esES`, `zhCN`, `zhTW`, `ptBR`, `ruRU`, `itIT`, `jaJP`, `plPL`.

**Response**:

```json
{
  "_links": { ... },
  "character": { ... },
  "equipped_item_sets": [ ... ],
  "equipped_items": [ ... ]
}
```

Example response:

````json
 "equipped_item_sets": [
    {
      "display_string": "Waycrest Legacy (1/2)",
      "effects": [
        {
          "display_string": "(2) Set: When you trigger a Cacaphonous Chord or Harmonious Chord there is a 40% chance to trigger the other chord at 60% of the power.",
          "required_count": 2
        }
      ],
      "item_set": {
        "id": 1439,
        "key": {
          "href": "https://eu.api.blizzard.com/data/wow/item-set/1439?namespace=static-11.0.0_55478-eu"
        },
        "name": "Waycrest Legacy"
      },
      "items": [
        {
          "is_equipped": true,
          "item": {
            "id": 158362,
            "key": {
              "href": "https://eu.api.blizzard.com/data/wow/item/158362?namespace=static-11.0.0_55478-eu"
            },
            "name": "Lord Waycrest's Signet"
          }
        },
        {
          "item": {
            "id": 159631,
            "key": {
              "href": "https://eu.api.blizzard.com/data/wow/item/159631?namespace=static-11.0.0_55478-eu"
            },
            "name": "Lady Waycrest's Music Box"
          }
        }
      ]
    },

```
````

### Specializations

Retrieves a character's specializations information, including spec groups, specs, and spec tiers.

**Endpoint**: `/blizzard/characters/{realmSlug}/{characterName}/specializations`
`/blizzard/characters/{realmSlug}/{characterName}/specializations?region={region}&namespace={namespace}&locale={locale}`

**Parameters**:

- `region`: The region of the character's realm. Possible values: `us`, `eu`, `kr`, `tw`, `cn`.
- `realmSlug`: The slug of the character's realm.
- `characterName`: The name of the character.
- `namespace`: The namespace of the character's profile.
- `locale`: The locale of the character's profile. Possible values: `en_US`, `enGB`, `deDE`, `frFR`, `koKR`, `esES`, `zhCN`, `zhTW`, `ptBR`, `ruRU`, `itIT`, `jaJP`, `plPL`.

**Response**:

```json
{
              "id": 82670,
              "rank": 1,
              "tooltip": {
                "spell_tooltip": {
                  "cast_time": "Passive",
                  "description": "Your direct damage spells inflict 30% of their damage on all other targets afflicted by your Vampiric Touch within 46 yards.\r\n\r\nDoes not apply to damage from Shadowy Apparitions, Shadow Word: Pain, and Vampiric Touch.",
                  "spell": {
                    "id": 199484,
                    "key": {
                      "href": "https://eu.api.blizzard.com/data/wow/spell/199484?namespace=static-11.0.0_55478-eu"
                    },
                    "name": "Psychic Link"
                  }
                },
                "talent": {
                  "id": 108819,
                  "key": {
                    "href": "https://eu.api.blizzard.com/data/wow/talent/108819?namespace=static-11.0.0_55478-eu"
                  },
                  "name": "Psychic Link"
                }
              }
            },
```

The API response take less than a second to return and is around 13k lines of json.
