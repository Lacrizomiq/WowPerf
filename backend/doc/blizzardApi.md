# Blizzard API

## Authentication

Blizzard API uses OAuth 2.0 to authenticate requests. The client ID and client secret are required to obtain an access token.

The access token is used to authenticate requests to the Blizzard API. The token is obtained using the OAuth 2.0 authorization code flow.

## Endpoints

### Character Profile

Retrieves a character's profile information, including name, realm, and class.

**Endpoint**: `/blizzard/characters/{realmSlug}/{characterName}`

**Parameters**:

- `region`: The region of the character's realm. Possible values: `us`, `eu`, `kr`, `tw`, `cn`.
- `realmSlug`: The slug of the character's realm.
- `characterName`: The name of the character.
- `namespace`: The namespace of the character's profile.
- `locale`: The locale of the character's profile. Possible values: `en_US`, `enGB`, `deDE`, `frFR`, `koKR`, `esES`, `zhCN`, `zhTW`, `ptBR`, `ruRU`, `itIT`, `jaJP`, `plPL`.

**Response**:

```json
{
  "name": "Example Character",
  "realm": {
    "name": "Example Realm",
    "slug": "example-realm"
  }
}
```

### Mythic Keystone Profile

Retrieves a character's mythic keystone profile information, including seasons, tiers, and keystone upgrades.

**Endpoint**: `/blizzard/characters/{realmSlug}/{characterName}/mythic-keystone-profile`

**Parameters**:

- `region`: The region of the character's realm. Possible values: `us`, `eu`, `kr`, `tw`, `cn`.
- `realmSlug`: The slug of the character's realm.
- `characterName`: The name of the character.
- `namespace`: The namespace of the character's profile.
- `locale`: The locale of the character's profile. Possible values: `enUS`, `enGB`, `deDE`, `frFR`, `koKR`, `esES`, `zhCN`, `zhTW`, `ptBR`, `ruRU`, `itIT`, `jaJP`, `plPL`.

**Response**:

```json
{
  "overall": {
    "season": 1,
    "tier": 1,
    "keystoneUpgrades": 0
  },
  "tier1": {
    "season": 1,
    "tier": 1,
    "keystoneUpgrades": 0
  }
}
```

### Equipment

Retrieves a character's equipment information, including items, gems, and mounts.

**Endpoint**: `/blizzard/characters/{realmSlug}/{characterName}/equipment`

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
