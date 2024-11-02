# WarcraftLogs Mythic+ API Documentation

## Dungeon Player Rankings Endpoint

`GET /warcraftlogs/mythicplus/rankings/dungeon/player`

This endpoint retrieves player rankings for specific dungeons in Mythic+.

### Required Parameters

| Parameter   | Type    | Description                 | Example |
| ----------- | ------- | --------------------------- | ------- |
| encounterID | integer | The ID of the dungeon       | 2579    |
| page        | integer | Page number (defaults to 1) | 1       |

### Optional Parameters

| Parameter    | Type   | Description                       | Example    |
| ------------ | ------ | --------------------------------- | ---------- |
| serverRegion | string | Filter by specific region         | EU, US, KR |
| serverSlug   | string | Filter by specific server         | kael-thas  |
| className    | string | Filter by specific class          | warrior    |
| specName     | string | Filter by specific specialization | arms       |

### Examples

#### Basic Usage (Required Parameters Only)

```
/warcraftlogs/mythicplus/rankings/dungeon/player?encounterID=2579&page=1
```

Returns the first page of rankings for dungeon ID 2579.

#### Filter by Region

```
/warcraftlogs/mythicplus/rankings/dungeon/player?encounterID=2579&page=1&serverRegion=EU
```

Returns rankings for EU region only.

#### Filter by Server

```
/warcraftlogs/mythicplus/rankings/dungeon/player?encounterID=2579&page=1&serverRegion=EU&serverSlug=kael-thas
```

Returns rankings for the Kael'thas server in EU region.

#### Filter by Class

```
/warcraftlogs/mythicplus/rankings/dungeon/player?encounterID=2579&page=1&className=warrior
```

Returns rankings for warrior players only.

#### Filter by Specialization

```
/warcraftlogs/mythicplus/rankings/dungeon/player?encounterID=2579&page=1&className=warrior&specName=arms
```

Returns rankings for arms warrior players only.

#### Complete Filter Example

```
/warcraftlogs/mythicplus/rankings/dungeon/player?encounterID=2579&page=1&serverRegion=EU&serverSlug=kael-thas&className=warrior&specName=arms
```

Returns rankings for arms warriors on the Kael'thas EU server.

### Response Format

```json
{
  "page": 1,
  "hasMorePages": true,
  "count": 50,
  "rankings": [
    {
      "name": "PlayerName",
      "class": "Warrior",
      "spec": "Arms",
      "score": 125.4,
      "server": {
        "name": "Kael'thas",
        "region": "EU"
      }
      // ... other player data
    }
    // ... more rankings
  ]
}
```

### Notes

- All optional parameters can be combined in any way
- If `page` is not specified, it defaults to 1
- Server filtering requires both `serverRegion` and `serverSlug` to be specific to a server
- Class and spec filtering can be used independently or together
- Results are ordered by score in descending order
