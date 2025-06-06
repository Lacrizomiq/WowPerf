# Repository Component - Player Rankings

## 📖 Vue d'ensemble

Le composant **Repository** gère toutes les opérations de persistance et les calculs SQL complexes. Il implémente le pattern Repository pour abstraire la couche de données et centraliser les requêtes de performance.

## 📁 Structure

```
repository/
└── player_rankings_repository.go    # Opérations CRUD et calculs SQL
```

## 🎯 Responsabilités

### ✅ Ce que fait ce composant

- 💾 **Opérations CRUD** sur les tables principales
- 📊 **Calculs de métriques** en SQL optimisé
- 🚀 **Insertions par batch** pour la performance
- 🔄 **Gestion des transactions** (atomicité)
- 📈 **Requêtes d'agrégation** complexes

### ❌ Ce qu'il ne fait PAS

- Appels vers des APIs externes
- Logique de déduplication métier
- Orchestration de workflows
- Parsing de données externes

## 🏗️ Interface du Repository

```go
type PlayerRankingsRepository struct {
    db *gorm.DB
}

// Méthodes principales
func (r *PlayerRankingsRepository) DeleteExistingRankings(ctx context.Context) error
func (r *PlayerRankingsRepository) StoreRankingsByBatches(ctx context.Context, rankings []PlayerRanking) error
func (r *PlayerRankingsRepository) CalculateDailySpecMetrics(ctx context.Context) error
func (r *PlayerRankingsRepository) GetGlobalRankings(ctx context.Context) (*GlobalRankings, error)
```

## 💾 Opérations de stockage

### 1. DeleteExistingRankings

```go
func (r *PlayerRankingsRepository) DeleteExistingRankings(ctx context.Context) error {
    log.Println("Deleting existing rankings")
    return r.db.WithContext(ctx).Exec("DELETE FROM player_rankings").Error
}
```

**Usage :** Nettoie la table avant chaque import quotidien pour garantir la fraîcheur des données.

**⚠️ Attention :** Opération destructive ! Assure-toi que les nouvelles données arrivent après.

### 2. StoreRankingsByBatches

```go
func (r *PlayerRankingsRepository) StoreRankingsByBatches(
    ctx context.Context,
    rankings []PlayerRanking
) error
```

#### Logique de batch processing

```go
const batchSize = 100 // Insertions par groupes de 100

totalRankings := len(rankings)
batches := int(math.Ceil(float64(totalRankings) / float64(batchSize)))

for i := 0; i < batches; i++ {
    start := i * batchSize
    end := start + batchSize
    if end > totalRankings {
        end = totalRankings
    }

    batch := rankings[start:end]
    // Insertion SQL du batch
}
```

#### Requête SQL optimisée

```sql
INSERT INTO player_rankings (
    created_at, updated_at, dungeon_id, name, class, spec, role,
    amount, hard_mode_level, duration, start_time, report_code,
    report_fight_id, report_start_time, guild_id, guild_name,
    guild_faction, server_id, server_name, server_region,
    bracket_data, faction, affixes, medal, score, leaderboard
) VALUES
    ($1, $2, $3, ..., $26),
    ($27, $28, $29, ..., $52),
    -- ... jusqu'à 100 lignes par batch
```

**Optimisations :**

- **Paramètres préparés** : Évite les injections SQL
- **Batch de 100** : Balance entre mémoire et performance
- **Transaction par batch** : Évite les locks longs

## 📊 Calculs de métriques

### CalculateDailySpecMetrics - Vue d'ensemble

Cette fonction est le **cœur analytique** du système. Elle calcule deux types de métriques :

1. **Métriques par donjon** (IsGlobal = false) - TOP 10 par spé/donjon
2. **Métriques globales** (IsGlobal = true) - TOP 10 joueurs ayant complété les 8 donjons

### 1. Métriques par donjon (TOP 10)

#### Requête SQL optimisée

```sql
WITH ranked_players AS (
    SELECT
        spec,
        class,
        role,
        dungeon_id,
        score,
        hard_mode_level,
        ROW_NUMBER() OVER (
            PARTITION BY spec, class, role, dungeon_id
            ORDER BY score DESC
        ) as player_rank
    FROM player_rankings
    WHERE server_region != 'CN'  -- Exclusion des serveurs chinois
)
SELECT
    spec,
    class,
    role,
    dungeon_id,
    COALESCE(AVG(score), 0) AS avg_score,
    COALESCE(MAX(score), 0) AS max_score,
    COALESCE(MIN(score), 0) AS min_score,
    COALESCE(AVG(hard_mode_level), 0) AS avg_key_level,
    COALESCE(MAX(hard_mode_level), 0) AS max_key_level,
    COALESCE(MIN(hard_mode_level), 0) AS min_key_level,
    COUNT(*) AS count
FROM ranked_players
WHERE player_rank <= 10  -- 🎯 TOP 10 seulement
GROUP BY spec, class, role, dungeon_id
```

**Résultat :** Une métrique par combinaison (spé, classe, rôle, donjon) basée sur les 10 meilleurs scores.

#### Exemple de données générées

```go
type DailySpecMetricMythicPlus struct {
    CaptureDate time.Time // 2024-01-15
    Spec        string    // "Protection"
    Class       string    // "Warrior"
    Role        string    // "Tank"
    EncounterID int       // 1200 (Cinderbrewery)
    IsGlobal    bool      // false
    AvgScore    float64   // 3750.5 (moyenne des TOP 10)
    MaxScore    float64   // 3890.2 (meilleur score)
    MinScore    float64   // 3654.1 (10ème meilleur score)
    AvgKeyLevel float64   // 18.7 (niveau de clé moyen)
    MaxKeyLevel int       // 22 (plus haut niveau)
    MinKeyLevel int       // 16 (plus bas niveau dans le TOP 10)
    RoleRank    int       // 2 (2ème meilleur tank pour ce donjon)
    OverallRank int       // 15 (15ème toutes spés confondues)
}
```

### 2. Métriques globales (TOP 10 players)

#### Calcul des scores totaux par joueur

```sql
SELECT
    spec,
    class,
    role,
    name,
    server_name,
    SUM(best_score) AS total_score,
    AVG(avg_key_level) AS avg_key_level,
    MAX(max_key_level) AS max_key_level,
    MIN(min_key_level) AS min_key_level,
    COUNT(DISTINCT dungeon_id) AS dungeon_count
FROM (
    SELECT
        spec, class, role, name, server_name, dungeon_id,
        MAX(score) as best_score,
        AVG(hard_mode_level) as avg_key_level,
        MAX(hard_mode_level) as max_key_level,
        MIN(hard_mode_level) as min_key_level
    FROM player_rankings
    WHERE server_region != 'CN'
    GROUP BY spec, class, role, name, server_name, dungeon_id
) AS best_scores
GROUP BY spec, class, role, name, server_name
HAVING COUNT(DISTINCT dungeon_id) = 8  -- Seulement les joueurs avec les 8 donjons
```

#### Sélection des TOP 10 par spécialisation

```go
// Tri par score total décroissant
sort.Slice(players, func(i, j int) bool {
    return players[i].TotalScore > players[j].TotalScore
})

// Sélection des TOP 10 (ou moins si pas assez de joueurs)
topCount := topPlayersCount // 10
if len(players) < topCount {
    topCount = len(players)
}
topPlayers := players[:topCount]

// Calcul des métriques sur ces TOP 10
avgScore := totalScores / float64(topCount)
```

**Résultat :** Une métrique globale par (spé, classe, rôle) basée sur les 10 meilleurs joueurs mondiaux.

## 🏆 Calcul des rankings

### Rankings par donjon

```go
func calculateDungeonMetricsRankings(metrics *[]DailySpecMetricMythicPlus) {
    // Grouper par donjon
    dungeonGroups := make(map[int][]int) // map[dungeonID][]metricIndex

    for _, indices := range dungeonGroups {
        // Tri par avgScore décroissant
        sort.Slice(indices, func(i, j int) bool {
            return (*metrics)[indices[i]].AvgScore > (*metrics)[indices[j]].AvgScore
        })

        // Attribution des ranks
        for rank, idx := range indices {
            (*metrics)[idx].OverallRank = rank + 1
        }

        // Rankings par rôle dans ce donjon
        // ...
    }
}
```

### Rankings globaux

```go
func calculateGlobalMetricsRankings(metrics *[]DailySpecMetricMythicPlus) {
    // Tri global par avgScore
    sort.Slice(*metrics, func(i, j int) bool {
        return (*metrics)[i].AvgScore > (*metrics)[j].AvgScore
    })

    // Attribution des ranks globaux
    for i := range *metrics {
        (*metrics)[i].OverallRank = i + 1
    }
}
```

## 📈 Requêtes de lecture

### GetGlobalRankings

```go
func (r *PlayerRankingsRepository) GetGlobalRankings(
    ctx context.Context
) (*GlobalRankings, error)
```

#### Calcul des scores totaux optimisé

```sql
SELECT
    name, class, spec, role,
    SUM(best_score) as total_score,
    guild_id, guild_name, guild_faction,
    server_id, server_name, server_region
FROM (
    SELECT
        name, class, spec, role, dungeon_id,
        MAX(score) as best_score,  -- Meilleur score par donjon
        guild_id, guild_name, guild_faction,
        server_id, server_name, server_region
    FROM player_rankings
    GROUP BY name, class, spec, role, dungeon_id,
             guild_id, guild_name, guild_faction,
             server_id, server_name, server_region
) as best_runs
GROUP BY name, class, spec, role,
         guild_id, guild_name, guild_faction,
         server_id, server_name, server_region
ORDER BY total_score DESC
```

#### Structure de retour

```go
type GlobalRankings struct {
    Tanks   RoleRankings `json:"tanks"`
    Healers RoleRankings `json:"healers"`
    DPS     RoleRankings `json:"dps"`
}

type RoleRankings struct {
    Count   int           `json:"count"`
    Players []PlayerScore `json:"players"`
}

type PlayerScore struct {
    Name       string  `json:"name"`
    Class      string  `json:"class"`
    Spec       string  `json:"spec"`
    Role       string  `json:"role"`
    TotalScore float64 `json:"total_score"`
    Guild      Guild   `json:"guild"`
    Server     Server  `json:"server"`
    Faction    int     `json:"faction"`
    Runs       []Run   `json:"runs"`
}
```

## 🚀 Optimisations de performance

### 1. Indexation recommandée

```sql
-- Index composé pour les métriques par donjon
CREATE INDEX idx_player_rankings_metrics
ON player_rankings (spec, class, role, dungeon_id, score DESC);

-- Index pour les calculs globaux
CREATE INDEX idx_player_rankings_global
ON player_rankings (name, server_name, dungeon_id, score DESC);

-- Index pour l'exclusion des serveurs chinois
CREATE INDEX idx_player_rankings_region
ON player_rankings (server_region) WHERE server_region != 'CN';
```

### 2. Transactions optimisées

```go
return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    // 1. DELETE rapide
    if err := tx.Exec("DELETE FROM daily_spec_metrics_mythic_plus WHERE capture_date = ?",
        processingDate).Error; err != nil {
        return err
    }

    // 2. Calculs en une seule requête
    if err := tx.Raw(complexQuery).Scan(&metrics).Error; err != nil {
        return err
    }

    // 3. INSERT par batches
    if err := tx.CreateInBatches(allMetrics, 100).Error; err != nil {
        return err
    }

    return nil // Commit atomique
})
```

### 3. Gestion mémoire

```go
// Streaming des insertions pour éviter l'OOM
for i := 0; i < batches; i++ {
    start := i * batchSize
    end := start + batchSize
    if end > totalRankings {
        end = totalRankings
    }

    batch := rankings[start:end]
    // Process batch immédiatement, pas de buffer global
}
```

## 📊 Métriques de performance

### Temps d'exécution typiques

- **DeleteExistingRankings** : 100-500ms (selon taille table)
- **StoreRankingsByBatches** : 2-5 minutes (45k rankings)
- **CalculateDailySpecMetrics** : 1-3 minutes (calculs complexes)
- **GetGlobalRankings** : 500ms-2s (lecture optimisée)

### Tailles de données

- **player_rankings** : ~45 000 lignes quotidiennes
- **daily_spec_metrics_mythic_plus** : ~2 500 métriques/jour (39 specs × 8 donjons × 2 types)
- **Batch size** : 100 insertions par lot
- **Transaction size** : Opération complète atomique

### Ressources consommées

- **CPU** : Calculs intensifs sur les CTE SQL
- **Mémoire** : ~50-100 MB pour les rankings en mémoire
- **I/O** : Lectures séquentielles optimisées par les index
- **Locks** : Transactions courtes pour minimiser les conflits

---

Ce composant est **crucial pour la performance** : des requêtes SQL optimisées ici impactent directement le temps d'exécution total du workflow.
