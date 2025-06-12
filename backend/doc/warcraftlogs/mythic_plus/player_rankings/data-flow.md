# Data Flow - Player Rankings

## 📖 Vue d'ensemble

Ce document détaille le flux complet des données dans le système Player Rankings, depuis l'API WarcraftLogs jusqu'aux vues finales utilisées par l'application.

## 🔄 Flux général

```
WarcraftLogs API ──► Queries ──► Activities ──► Repository ──► Database ──► Views
      │                 │           │             │             │           │
   GraphQL          Parsing    Déduplication   Calculs      Tables      APIs
   Responses                                    SQL         Finales   Application
```

## 📊 Volume de données

### Données d'entrée (quotidiennes)

- **8 donjons** Mythic+ de la saison
- **2 pages par donjon** (~100 top players/page)
- **39 spécialisations** WoW
- **624 requêtes GraphQL** total (8 × 2 × 39)
- **~60 000 rankings bruts** récupérés

### Données traitées (quotidiennes)

- **~45 000 rankings uniques** après déduplication (~25% de réduction)
- **~2 500 métriques** calculées (39 specs × 8 donjons × 2 types)
- **~600 MB** de données insérées en base

## 🚀 Étape 1 : Récupération depuis WarcraftLogs

### Point de départ : Temporal Schedule

```
12h00 UTC quotidien
       │
       ▼
PlayerRankingsWorkflow démarre
       │
       ▼
FetchAllDungeonRankings Activity
```

### Boucle de récupération par donjon

```
Pour chaque donjon (8 au total):
├── Pour chaque page (1-2):
│   ├── Pour chaque spécialisation (39 au total):
│   │   ├── Requête GraphQL vers WarcraftLogs
│   │   ├── Parsing JSON → Ranking objects
│   │   └── Collecte temporaire en mémoire
│   └── Déduplication par joueur (nom + serveur)
└── Stockage en base par batch
```

### Exemple de requête GraphQL

```graphql
query getDungeonLeaderboard(
  $encounterId: Int!
  $page: Int!
  $className: String!
  $specName: String!
) {
  worldData {
    encounter(id: $encounterId) {
      characterRankings(
        leaderboard: Any
        page: $page
        className: $className
        specName: $specName
      ) {
        rankings {
          name
          class
          spec
          score
          amount
          hardModeLevel
          duration
          # ... autres champs
        }
      }
    }
  }
}
```

### Transformation des données

```
WarcraftLogs Response → Parsing → PlayerRanking Model

{                         │         {
  "name": "Playerx",     │           "Name": "Playerx",
  "class": "Warrior",    │    ──►    "Class": "Warrior",
  "spec": "Protection",  │           "Spec": "Protection",
  "score": 3840.5       │           "Role": "Tank",        // ← Calculé
}                        │           "Score": 3840.5
                         │         }
```

## 🔀 Étape 2 : Déduplication avancée

### Problème résolu

Un même joueur peut apparaître plusieurs fois car :

1. **Pagination** : Présent sur plusieurs pages d'un donjon
2. **Spécialisations multiples** : API peut retourner le même joueur pour différentes specs
3. **Runs multiples** : Plusieurs tentatives du même donjon

### Solution : Map de déduplication

```go
type playerDungeonKey struct {
    name      string // "Playerx-Stormrage"
    dungeonID int    // 1200
}

bestScores := make(map[playerDungeonKey]*PlayerRanking)

// Pour chaque ranking récupéré
key := playerDungeonKey{
    name:      fmt.Sprintf("%s-%s", ranking.Name, ranking.Server.Name),
    dungeonID: dungeonID,
}

// Garde seulement le MEILLEUR score
if existing, exists := bestScores[key]; exists {
    if ranking.Score > existing.Score {
        bestScores[key] = ranking // Remplace par le meilleur
    }
} else {
    bestScores[key] = ranking // Premier score pour cette clé
}
```

### Résultat de la déduplication

```
Avant: ~60 000 rankings (avec doublons)
Après: ~45 000 rankings (uniques)
Réduction: ~25% de doublons éliminés
```

## 💾 Étape 3 : Stockage en base

### Stratégie : Delete + Insert

```sql
-- 1. Nettoyage des données existantes
DELETE FROM player_rankings;

-- 2. Insertion par batch de 100
INSERT INTO player_rankings (
    created_at, updated_at, dungeon_id, name, class, spec, role,
    amount, hard_mode_level, duration, start_time, report_code,
    -- ... 26 colonnes au total
) VALUES
    ($1, $2, $3, ..., $26),
    ($27, $28, $29, ..., $52),
    -- ... jusqu'à 100 lignes par batch
```

### Structure de la table player_rankings

```sql
CREATE TABLE player_rankings (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,

    -- Données du donjon
    dungeon_id INTEGER NOT NULL,
    hard_mode_level INTEGER NOT NULL,
    duration BIGINT NOT NULL,
    start_time BIGINT NOT NULL,

    -- Données du joueur
    name VARCHAR(255) NOT NULL,
    class VARCHAR(50) NOT NULL,
    spec VARCHAR(50) NOT NULL,
    role VARCHAR(20) NOT NULL,

    -- Métriques de performance
    score DECIMAL(10,3) NOT NULL,
    amount INTEGER NOT NULL,
    medal VARCHAR(20),
    affixes INTEGER[],
    bracket_data INTEGER,
    faction INTEGER,

    -- Données de guilde
    guild_id INTEGER,
    guild_name VARCHAR(255),
    guild_faction INTEGER,

    -- Données de serveur
    server_id INTEGER NOT NULL,
    server_name VARCHAR(255) NOT NULL,
    server_region VARCHAR(10) NOT NULL,

    -- Données de rapport
    report_code VARCHAR(50),
    report_fight_id INTEGER,
    report_start_time BIGINT,

    leaderboard INTEGER DEFAULT 0
);
```

### Index pour performance

```sql
-- Index pour métriques par donjon
CREATE INDEX idx_player_rankings_metrics
ON player_rankings (spec, class, role, dungeon_id, score DESC);

-- Index pour calculs globaux
CREATE INDEX idx_player_rankings_global
ON player_rankings (name, server_name, dungeon_id, score DESC);

-- Index pour exclusion région chinoise
CREATE INDEX idx_player_rankings_region
ON player_rankings (server_region) WHERE server_region != 'CN';
```

## 📊 Étape 4 : Calcul des métriques

### CalculateDailyMetrics Activity

```
CalculateDailyMetrics démarre
       │
       ▼
Suppression métriques existantes pour today
       │
       ▼
Calcul métriques par donjon (TOP 10)
       │
       ▼
Calcul métriques globales (TOP 10 players)
       │
       ▼
Calcul des rankings (par rôle et global)
       │
       ▼
Stockage dans daily_spec_metrics_mythic_plus
```

### Métriques par donjon (IsGlobal = false)

```sql
WITH ranked_players AS (
    SELECT
        spec, class, role, dungeon_id, score, hard_mode_level,
        ROW_NUMBER() OVER (
            PARTITION BY spec, class, role, dungeon_id
            ORDER BY score DESC
        ) as player_rank
    FROM player_rankings
    WHERE server_region != 'CN'
)
SELECT
    spec, class, role, dungeon_id,
    AVG(score) AS avg_score,
    MAX(score) AS max_score,
    MIN(score) AS min_score,
    AVG(hard_mode_level) AS avg_key_level,
    COUNT(*) AS count
FROM ranked_players
WHERE player_rank <= 10  -- 🎯 TOP 10 seulement
GROUP BY spec, class, role, dungeon_id
```

**Résultat :** ~2000 métriques (39 specs × 8 donjons × ~66% coverage)

### Métriques globales (IsGlobal = true)

```sql
-- 1. Calcul des scores totaux par joueur
SELECT
    spec, class, role, name, server_name,
    SUM(best_score) AS total_score,
    AVG(avg_key_level) AS avg_key_level
FROM (
    SELECT
        spec, class, role, name, server_name, dungeon_id,
        MAX(score) as best_score,
        AVG(hard_mode_level) as avg_key_level
    FROM player_rankings
    WHERE server_region != 'CN'
    GROUP BY spec, class, role, name, server_name, dungeon_id
) AS best_scores
GROUP BY spec, class, role, name, server_name
HAVING COUNT(DISTINCT dungeon_id) = 8  -- Seulement joueurs avec 8 donjons
```

```go
// 2. Sélection TOP 10 par spécialisation en Go
sort.Slice(players, func(i, j int) bool {
    return players[i].TotalScore > players[j].TotalScore
})

topPlayers := players[:min(10, len(players))]

// 3. Calcul des métriques sur ces TOP 10
avgScore := totalScores / float64(len(topPlayers))
```

**Résultat :** ~39 métriques globales (une par spécialisation active)

### Structure daily_spec_metrics_mythic_plus

```sql
CREATE TABLE daily_spec_metrics_mythic_plus (
    id SERIAL PRIMARY KEY,
    capture_date DATE NOT NULL,

    -- Identifiants
    spec VARCHAR(50) NOT NULL,
    class VARCHAR(50) NOT NULL,
    role VARCHAR(20) NOT NULL,
    encounter_id INTEGER NOT NULL,  -- 0 pour global, sinon dungeon_id
    is_global BOOLEAN NOT NULL,

    -- Métriques calculées
    avg_score DECIMAL(10,3) NOT NULL,
    max_score DECIMAL(10,3) NOT NULL,
    min_score DECIMAL(10,3) NOT NULL,
    avg_key_level DECIMAL(5,2) NOT NULL,
    max_key_level INTEGER NOT NULL,
    min_key_level INTEGER NOT NULL,

    -- Rankings
    role_rank INTEGER NOT NULL,
    overall_rank INTEGER NOT NULL,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## 📈 Étape 5 : Génération des vues

### Vue spec_global_score_averages

```sql
CREATE VIEW spec_global_score_averages AS
SELECT
    class,
    spec,
    LOWER(REPLACE(CONCAT(class, spec), ' ', '-')) as slug,
    avg_score as avg_global_score,
    max_score as max_global_score,
    min_score as min_global_score,
    10 as player_count,  -- TOP 10 utilisé
    role,
    overall_rank,
    role_rank
FROM daily_spec_metrics_mythic_plus
WHERE is_global = true
  AND capture_date = CURRENT_DATE
ORDER BY avg_score DESC;
```

## 🔄 Transformations de données clés

### 1. Détermination du rôle

```
Input: Class="Warrior", Spec="Protection"
Logic: Lookup dans map[class][]specs tanks/healers
Output: Role="Tank"
```

### 2. Calcul des scores totaux

```
Input: 8 runs d'un joueur (un par donjon)
Logic: MAX(score) par donjon puis SUM des 8 scores
Output: TotalScore (ex: 3740.5)
```

### 3. Rankings par spécialisation

```
Input: Toutes les métriques d'un donjon
Logic: ORDER BY avg_score DESC puis attribution rank
Output: OverallRank, RoleRank par métrique
```

## 📊 Validation de la qualité des données

### Checks automatiques

```sql
-- 1. Vérification des doublons (doit être 0)
SELECT COUNT(*) FROM (
    SELECT name, server_name, dungeon_id, COUNT(*)
    FROM player_rankings
    GROUP BY name, server_name, dungeon_id
    HAVING COUNT(*) > 1
);

-- 2. Cohérence des scores (doit être < 4000 généralement)
SELECT MAX(total_score) FROM (
    SELECT SUM(best_score) as total_score FROM (
        SELECT name, dungeon_id, MAX(score) as best_score
        FROM player_rankings GROUP BY name, dungeon_id
    ) grouped GROUP BY name
);

-- 3. Couverture des spécialisations (doit être ~39)
SELECT COUNT(DISTINCT CONCAT(class, '-', spec))
FROM daily_spec_metrics_mythic_plus
WHERE capture_date = CURRENT_DATE AND is_global = true;
```

### Métriques de santé

- **Doublons** : 0 après déduplication
- **Score max** : < 4000 points (sanity check)
- **Couverture specs** : 35-39 spécialisations actives
- **Ratio rôles** : ~7% Tanks, ~12% Healers, ~81% DPS

---

Ce flux garantit la **cohérence**, **fraîcheur** et **qualité** des données Player Rankings de bout en bout.
