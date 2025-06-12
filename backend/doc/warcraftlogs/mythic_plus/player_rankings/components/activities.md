# Activities Component - Player Rankings

## 📖 Vue d'ensemble

Les **Activities** représentent la couche de logique métier dans l'architecture Temporal. Elles coordonnent les appels entre les couches Queries et Repository, gèrent la déduplication des données et implémentent la logique de traitement parallèle.

## 📁 Structure

```
activities/
├── activities.go                    # Factory et coordination
└── player_rankings_activity.go     # Activities principales
```

## 🎯 Responsabilités

### ✅ Ce que fait ce composant

- 🎭 **Logique métier** de coordination
- 🔄 **Orchestration** entre Queries et Repository
- 🔀 **Déduplication avancée** des données
- ⚡ **Traitement parallèle** avec contrôle de concurrence
- 📊 **Calcul des statistiques** métier
- 💫 **Gestion des erreurs** avec retry Temporal

### ❌ Ce qu'il ne fait PAS

- Gestion d'état de workflow (rôle du Workflow)
- Requêtes SQL directes (rôle du Repository)
- Parsing GraphQL (rôle des Queries)
- Configuration des schedules (rôle du Scheduler)

## 🏗️ Structure des Activities

### 1. PlayerRankingsActivity - Interface

```go
type PlayerRankingsActivity struct {
    client     *warcraftlogs.WarcraftLogsClientService
    repository *playerRankingsRepository.PlayerRankingsRepository
}

// Injection de dépendances
func NewPlayerRankingsActivity(
    client *warcraftlogs.WarcraftLogsClientService,
    repository *playerRankingsRepository.PlayerRankingsRepository,
) *PlayerRankingsActivity
```

### 2. Activities principales

| Activity                    | Durée     | Rôle                          |
| --------------------------- | --------- | ----------------------------- |
| **FetchAllDungeonRankings** | 15-30 min | Récupération et déduplication |
| **CalculateDailyMetrics**   | 2-5 min   | Calcul des métriques          |
| **StoreRankings**           | 2-5 min   | Stockage en base (legacy)     |

## 🚀 FetchAllDungeonRankings Activity

### Signature et paramètres

```go
func (a *PlayerRankingsActivity) FetchAllDungeonRankings(
    ctx context.Context,
    dungeonIDs []int,        // [1200, 1201, ...] - 8 donjons
    pagesPerDungeon int,     // 2 pages par donjon
    maxConcurrency int,      // 3 requêtes parallèles max
) (*models.RankingsStats, error)
```

### Logique de déduplication avancée

#### Structure de clé unique

```go
type playerDungeonKey struct {
    name      string // "PlayerName-ServerName"
    dungeonID int    // ID du donjon
}

bestScores := make(map[playerDungeonKey]*playerRankingModels.PlayerRanking)
```

**Pourquoi cette clé ?**

- **Name + Server** : Évite les confusions entre homonymes
- **DungeonID** : Un joueur peut avoir un score différent par donjon
- **Résultat** : Un score maximum par joueur par donjon

#### Algorithme de déduplication

```go
for _, ranking := range dungeonData.Rankings {
    key := playerDungeonKey{
        name:      fmt.Sprintf("%s-%s", ranking.Name, ranking.Server.Name),
        dungeonID: dID,
    }

    // Garder seulement le MEILLEUR score pour cette clé
    if existing, exists := bestScores[key]; exists {
        if ranking.Score > existing.Score {
            bestScores[key] = playerRanking // Remplace par le meilleur
        }
    } else {
        bestScores[key] = playerRanking // Premier score pour cette clé
    }
}
```

### Gestion de la concurrence

#### Semaphore pattern

```go
sem := make(chan struct{}, maxConcurrency) // Canal buffered à 3

for _, dungeonID := range dungeonIDs {
    wg.Add(1)
    go func(dID int) {
        defer wg.Done()

        sem <- struct{}{}        // Acquiert le slot (bloque si 3 déjà pris)
        defer func() { <-sem }() // Libère le slot

        // Traitement du donjon...
    }(dungeonID)
}
```

**Avantages :**

- ⚡ **Parallélisme contrôlé** : Max 3 requêtes simultanées vers WarcraftLogs
- 🛡️ **Protection API** : Évite le rate limiting
- 🔄 **Heartbeat Temporal** : Signale la progression

#### Heartbeat pour long-running activity

```go
// Signal de vie pour Temporal toutes les X secondes
activity.RecordHeartbeat(ctx, fmt.Sprintf("Processing dungeon %d", dID))
```

### Mapping vers les modèles de données

#### Détermination du rôle

```go
func determineRole(class, spec string) string {
    tanks := map[string][]string{
        "Warrior":     {"Protection"},
        "Paladin":     {"Protection"},
        "DeathKnight": {"Blood"},
        "DemonHunter": {"Vengeance"},
        "Druid":       {"Guardian"},
        "Monk":        {"Brewmaster"},
    }

    healers := map[string][]string{
        "Priest":  {"Holy", "Discipline"},
        "Paladin": {"Holy"},
        "Druid":   {"Restoration"},
        "Shaman":  {"Restoration"},
        "Monk":    {"Mistweaver"},
        "Evoker":  {"Preservation"},
    }

    // Logique de classification...
    return "Tank" | "Healer" | "DPS"
}
```

#### Transformation des données

```go
playerRanking := &playerRankingModels.PlayerRanking{
    DungeonID:       dID,
    Name:            ranking.Name,
    Class:           ranking.Class,
    Spec:            ranking.Spec,
    Role:            determineRole(ranking.Class, ranking.Spec),
    Amount:          ranking.Amount,
    HardModeLevel:   ranking.HardModeLevel,
    Duration:        ranking.Duration,
    StartTime:       ranking.StartTime,
    ReportCode:      ranking.Report.Code,
    ReportFightID:   ranking.Report.FightID,
    ReportStartTime: ranking.Report.StartTime,
    // Données de guilde et serveur...
    Score:           ranking.Score,
}
```

### Statistiques calculées

```go
// Comptage par rôle pendant la déduplication
for _, ranking := range bestScores {
    switch ranking.Role {
    case "Tank":
        tankCount++
    case "Healer":
        healerCount++
    case "DPS":
        dpsCount++
    }
}

// Retour des statistiques
return &models.RankingsStats{
    TotalCount:        len(allRankings),    // ~45 000 rankings
    DungeonsProcessed: len(dungeonIDs),     // 8 donjons
    ProcessingTime:    time.Since(startTime), // ~20-30 minutes
    TankCount:         tankCount,           // ~3 000 tanks
    HealerCount:       healerCount,         // ~5 000 healers
    DPSCount:          dpsCount,           // ~37 000 DPS
}
```

## 📊 CalculateDailyMetrics Activity

### Signature simple

```go
func (a *PlayerRankingsActivity) CalculateDailyMetrics(ctx context.Context) error
```

### Rôle de coordination

```go
func (a *PlayerRankingsActivity) CalculateDailyMetrics(ctx context.Context) error {
    logger := activity.GetLogger(ctx)
    logger.Info("Starting daily metrics calculation")

    startTime := time.Now()

    // Délégation au Repository pour les calculs SQL complexes
    if err := a.repository.CalculateDailySpecMetrics(ctx); err != nil {
        logger.Error("Failed to calculate daily metrics", "error", err)
        return temporal.NewApplicationError(
            fmt.Sprintf("Failed to calculate daily metrics: %v", err),
            "METRICS_ERROR",
        )
    }

    duration := time.Since(startTime)
    logger.Info("Successfully calculated daily metrics", "duration", duration)
    return nil
}
```

**Simplification volontaire :** Cette activity est un simple wrapper qui :

- ✅ Ajoute les logs Temporal
- ✅ Gère les erreurs dans le format Temporal
- ✅ Mesure les performances
- ✅ Délègue la complexité au Repository

## 💾 StoreRankings Activity (Legacy)

### Usage historique

```go
func (a *PlayerRankingsActivity) StoreRankings(
    ctx context.Context,
    rankings []playerRankingModels.PlayerRanking,
) error
```

**Note :** Cette activity n'est plus utilisée dans le workflow actuel. `FetchAllDungeonRankings` stocke directement les données pour optimiser les performances.

**Garde pour :**

- Compatibilité avec anciens workflows
- Tests unitaires
- Usage manuel éventuel

## 🔧 Patterns appliqués

### 1. Coordinator Pattern

```go
// L'activity coordonne sans implémenter
func (a *Activity) FetchAllDungeonRankings(...) {
    // 1. Utilise Queries pour récupérer
    data, err := queries.GetDungeonLeaderboardByPlayer(...)

    // 2. Applique la logique métier (déduplication)
    deduplicated := deduplicateData(data)

    // 3. Utilise Repository pour stocker
    err = a.repository.StoreRankingsByBatches(deduplicated)
}
```

### 2. Error Wrapping Pattern

```go
// Wrapping des erreurs pour Temporal
if err != nil {
    return temporal.NewApplicationError(
        fmt.Sprintf("Failed to fetch dungeon %d: %v", dungeonID, err),
        "FETCH_ERROR", // Type pour retry policies
    )
}
```

### 3. Progress Reporting Pattern

```go
// Reporting de progression pour long-running
activity.RecordHeartbeat(ctx, fmt.Sprintf("Processing dungeon %d", dID))

// Logs structurés avec métadonnées
logger.Info("Completed dungeon rankings fetch",
    "dungeonID", dID,
    "rankingsCount", len(dungeonRankings))
```

## ⚡ Optimisations de performance

### 1. Déduplication en mémoire

```go
// Map plus efficace que les slices pour la déduplication
bestScores := make(map[playerDungeonKey]*PlayerRanking)
// O(1) lookup vs O(n) avec slice.Find()
```

### 2. Stockage direct

```go
// Évite de transférer 45k rankings entre activity et workflow
// Stockage direct dans l'activity au lieu de retourner les données
if err := a.repository.StoreRankingsByBatches(ctx, allRankings); err != nil {
    return err
}
```

### 3. Parallélisme contrôlé

```go
// Semaphore évite la surcharge tout en maximisant le throughput
sem := make(chan struct{}, 3) // Sweet spot pour WarcraftLogs API
```

### 4. Early termination

```go
// Arrêt dès qu'il n'y a plus de pages
if !dungeonData.HasMorePages {
    break
}
```

## 🐛 Gestion des erreurs

### Stratégies par type d'erreur

#### 1. Erreurs récupérables (Retry)

```go
// Timeouts API, rate limits
errorsChan <- fmt.Errorf("failed to fetch dungeon %d page %d: %w", dID, page, err)
// Temporal va retry automatiquement selon la RetryPolicy
```

#### 2. Erreurs de données (Continue)

```go
// Données manquantes pour une spécialisation
log.Printf("No rankings found for class %s, spec %s", spec.ClassName, spec.SpecName)
continue // Continue avec les autres spécialisations
```

#### 3. Erreurs systémiques (Fail fast)

```go
// Erreurs de base de données, configuration
return temporal.NewApplicationError(
    fmt.Sprintf("Database error: %v", err),
    "DB_ERROR",
)
```

### Collecte d'erreurs

```go
// Collecte toutes les erreurs sans arrêter le traitement
errorsChan := make(chan error, len(dungeonIDs)*pagesPerDungeon)

// Processing...

// Évaluation finale
var errors []error
for err := range errorsChan {
    errors = append(errors, err)
}

if len(errors) > 0 {
    return temporal.NewApplicationError(
        fmt.Sprintf("Encountered %d errors during processing", len(errors)),
        "PARTIAL_FAILURE",
    )
}
```

## 📈 Métriques et observabilité

### Logs structurés

```go
logger.Info("Starting rankings fetch for multiple dungeons",
    "dungeonCount", len(dungeonIDs),
    "pagesPerDungeon", pagesPerDungeon,
    "maxConcurrency", maxConcurrency)

logger.Info("Successfully fetched rankings for all dungeons",
    "totalRankingsCount", len(allRankings),
    "tankCount", tankCount,
    "healerCount", healerCount,
    "dpsCount", dpsCount,
    "processingTime", time.Since(startTime))
```

### Métriques de performance typiques

- **Fetching** : 15-30 minutes pour 8 donjons × 2 pages
- **Deduplication** : ~10-20% de réduction des doublons
- **Storage** : 2-5 minutes pour 45k rankings
- **Metrics calculation** : 1-3 minutes pour calculs SQL

---

Les Activities constituent la **logique métier centrale** qui transforme les données brutes en informations exploitables tout en respectant les contraintes de performance et de fiabilité.
