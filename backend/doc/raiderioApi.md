# Raider.io API

## Endpoints

### Character Profile

Retrieves a character's profile information, including name, realm, and class.

**Endpoint**: `/characters/profile`

**Parameters**:

- `region`: The region of the character's realm. Possible values: `us`, `eu`, `kr`, `tw`, `cn`.
- `realm`: The name of the character's realm.
- `name`: The name of the character.

**Response**:

```json
{
  "name": "Example Character",
  "race": "Human",
  "class": "Warrior",
  "active_spec_name": "Assassination",
  "active_spec_role": "Tank",
  "gender": "Male",
  "faction": "Alliance",
  "achievement_points": 100,
  "honorable_kills": 10,
  "thumbnail_url": "https://example.com/thumbnail.png",
  "region": "us",
  "realm": "Example Realm",
  "profile_url": "https://example.com/profile",
  "gear": {
    "item_level_equipped": 1,
    "item_level_total": 1
  },
  "guild": {
    "name": "Example Guild",
    "realm": "Example Realm"
  },
  "raid_progression": {
    "summary": "100% Complete",
    "total_bosses": 10,
    "normal_bosses_killed": 10,
    "heroic_bosses_killed": 0,
    "mythic_bosses_killed": 0
  },
  "mythic_plus_scores_by_season": [
    {
      "season": "1",
      "scores": {
        "all": 100,
        "dps": 100,
        "healer": 100,
        "tank": 100,
        "spec_0": 100,
        "spec_1": 100,
        "spec_2": 100,
        "spec_3": 100
      },
      "segments": {
        "all": {
          "score": 100,
          "color": "#00FF00"
        },
        "dps": {
          "score": 100,
          "color": "#00FF00"
        },
        "healer": {
          "score": 100,
          "color": "#00FF00"
        }
      }
    }
  ],
  "mythic_plus_ranks": {
    "overall": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "tank": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "healer": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "dps": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "class": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "class_tank": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "class_healer": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "class_dps": {
      "world": 1,
      "region": 1,
      "realm": 1
    }
  },
  "mythic_plus_recent_runs": [
    {
      "dungeon": "Example Dungeon",
      "short_name": "Example Short Name",
      "mythic_level": 1,
      "completed_at": "2023-01-01T00:00:00Z",
      "clear_time_ms": 1000,
      "num_keystone_upgrades": 1,
      "score": 100,
      "url": "https://example.com/run"
    }
  ],
  "mythic_plus_best_runs": [
    {
      "dungeon": "Example Dungeon",
      "short_name": "Example Short Name",
      "mythic_level": 1,
      "completed_at": "2023-01-01T00:00:00Z",
      "clear_time_ms": 1000,
      "num_keystone_upgrades": 1,
      "score": 100,
      "url": "https://example.com/run"
    }
  ],
  "mythic_plus_alternate_runs": [
    {
      "dungeon": "Example Dungeon",
      "short_name": "Example Short Name",
      "mythic_level": 1,
      "completed_at": "2023-01-01T00:00:00Z",
      "clear_time_ms": 1000,
      "num_keystone_upgrades": 1,
      "score": 100,
      "url": "https://example.com/run"
    }
  ],
  "mythic_plus_highest_level_runs": [
    {
      "dungeon": "Example Dungeon",
      "short_name": "Example Short Name",
      "mythic_level": 1,
      "completed_at": "2023-01-01T00:00:00Z",
      "clear_time_ms": 1000,
      "num_keystone_upgrades": 1,
      "score": 100,
      "url": "https://example.com/run"
    }
  ],
  "mythic_plus_weekly_highest_level_runs": [
    {
      "dungeon": "Example Dungeon",
      "short_name": "Example Short Name",
      "mythic_level": 1,
      "completed_at": "2023-01-01T00:00:00Z",
      "clear_time_ms": 1000,
      "num_keystone_upgrades": 1,
      "score": 100,
      "url": "https://example.com/run"
    }
  ],
  "mythic_plus_previous_weekly_highest_level_runs": [
    {
      "dungeon": "Example Dungeon",
      "short_name": "Example Short Name",
      "mythic_level": 1,
      "completed_at": "2023-01-01T00:00:00Z",
      "clear_time_ms": 1000,
      "num_keystone_upgrades": 1,
      "score": 100,
      "url": "https://example.com/run"
    }
  ],
  "previous_mythic_plus_ranks": {
    "overall": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "tank": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "healer": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "dps": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "class": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "class_tank": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "class_healer": {
      "world": 1,
      "region": 1,
      "realm": 1
    },
    "class_dps": {
      "world": 1,
      "region": 1,
      "realm": 1
    }
  }
}
```

### Mythic Plus Scores

Retrieves a character's mythic plus scores, including seasons, tiers, and keystone upgrades.

**Endpoint**: `/characters/mythic-plus-scores`

**Parameters**:

- `region`: The region of the character's realm. Possible values: `us`, `eu`, `kr`, `tw`, `cn`.
- `realm`: The name of the character's realm.
- `name`: The name of the character.

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
  },
  "tier2": {
    "season": 1,
    "tier": 1,
    "keystoneUpgrades": 1
  }
}
```

### Raid Progression

Retrieves a character's raid progression, including summary, total bosses, normal bosses killed, heroic bosses killed, and mythic bosses killed.

**Endpoint**: `/characters/raid-progression`

**Parameters**:

- `region`: The region of the character's realm. Possible values: `us`, `eu`, `kr`, `tw`, `cn`.
- `realm`: The name of the character's realm.
- `name`: The name of the character.

**Response**:

```json
{
  "summary": "100% Complete",
  "total_bosses": 10,
  "normal_bosses_killed": 10,
  "heroic_bosses_killed": 0,
  "mythic_bosses_killed": 0
}
```
