# WowPerf API Documentation - Mythic+ Leaderboards

## Global Leaderboards

### Get Global Leaderboard

```http
GET /warcraftlogs/mythicplus/global/leaderboard
```

Retrieves the overall global leaderboard across all roles, classes and specs.

**Query Parameters:**

- `limit` (optional): Number of entries to return (default: 100)
  - Example: `?limit=10`

**Response Example:**

```json
[
  {
    "player_id": 85877948,
    "name": "Kirawarrior",
    "class": "Warrior",
    "spec": "Protection",
    "role": "tank",
    "total_score": 430.26,
    "rank": 1
  }
  // ...
]
```

### Get Role Leaderboard

```http
GET /warcraftlogs/mythicplus/global/leaderboard/role
```

Retrieves the leaderboard for a specific role.

**Query Parameters:**

- `role` (required): Role to filter by ("tank", "healer", "dps")
- `limit` (optional): Number of entries to return (default: 100)
  - Example: `?role=tank&limit=5`

### Get Class Leaderboard

```http
GET /warcraftlogs/mythicplus/global/leaderboard/class
```

Retrieves the leaderboard for a specific class.

**Query Parameters:**

- `class` (required): Class to filter by (e.g., "warrior", "mage", "druid")
- `limit` (optional): Number of entries to return (default: 100)
  - Example: `?class=warrior&limit=10`

### Get Spec Leaderboard

```http
GET /warcraftlogs/mythicplus/global/leaderboard/spec
```

Retrieves the leaderboard for a specific specialization.

**Query Parameters:**

- `class` (required): Class name (e.g., "warrior", "mage")
- `spec` (required): Specialization name (e.g., "protection", "fire")
- `limit` (optional): Number of entries to return (default: 100)
  - Example: `?class=warrior&spec=protection&limit=5`

## Dungeon Leaderboards

### Get Dungeon Leaderboard

```http
GET /warcraftlogs/mythicplus/rankings/dungeon
```

Retrieves the leaderboard for a specific dungeon.

**Query Parameters:**

- `encounterID` (required): ID of the dungeon
- `page` (optional): Page number for pagination (default: 1)
  - Example: `?encounterID=12660&page=1`

**Response Example:**

```json
{
  "rankings": [
    {
      "score": 430.26,
      "duration": 1234567,
      "bracket_data": 20,
      "medal": "gold",
      "team": [
        {
          "id": 85877948,
          "name": "Kirawarrior",
          "class": "Warrior",
          "spec": "Protection",
          "role": "tank"
        }
        // ... other team members
      ]
    }
    // ... other rankings
  ]
}
```

## Common Response Status Codes

- `200 OK`: Request successful
- `400 Bad Request`: Invalid parameters provided
- `500 Internal Server Error`: Server-side error occurred

## Notes

- All endpoints return JSON responses
- Times are returned in milliseconds
- Scores are returned as floating-point numbers
- Role values are lowercase ("tank", "healer", "dps")
- Class and spec names match WoW's official naming
