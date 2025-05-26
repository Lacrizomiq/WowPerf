# Architecture - Player Rankings

## 🏗️ Vue d'ensemble architecturale

La feature Player Rankings suit une architecture en couches avec séparation claire des responsabilités, orchestrée par Temporal pour la fiabilité et la scalabilité.

## 🔄 Architecture en couches

```
┌─────────────────────────────────────────────────────────┐
│                 Temporal Scheduler                     │
│            (Cron quotidien 12h UTC)                    │
└─────────────────────┬───────────────────────────────────┘
                      │ Trigger
┌─────────────────────▼───────────────────────────────────┐
│              Temporal Workflow                         │
│        (Orchestration, Retry, State Management)        │
└─────┬─────────────────────────────────┬─────────────────┘
      │ Execute                         │ Execute
┌─────▼──────────────┐         ┌───────▼──────────────────┐
│  Activities Layer  │         │    Activities Layer     │
│ FetchAllDungeons   │         │ CalculateDailyMetrics   │
└─────┬──────────────┘         └───────┬──────────────────┘
      │ Use                            │ Use
┌─────▼──────────────┐         ┌───────▼──────────────────┐
│   Queries Layer    │         │   Repository Layer      │
│ (WarcraftLogs API) │         │   (Database SQL)        │
└────────────────────┘         └──────────────────────────┘
```

## 🧩 Responsabilités par couche

### 1. Temporal Scheduler

**Rôle** : Déclenchement automatique quotidien

**Responsabilités :**

- ⏰ Exécution quotidienne à 12h UTC
- 🔧 Configuration des retry policies
- 📊 Gestion du cycle de vie des workflows
- 🚨 Alerting en cas d'échec

**Ne fait PAS :**

- Logique métier
- Appels directs aux APIs
- Manipulations de données

### 2. Temporal Workflow

**Rôle** : Orchestration et coordination

**Responsabilités :**

- 🎭 Orchestration des étapes de traitement
- ⏱️ Gestion des timeouts et retry
- 📝 Logging et tracking de la progression
- 🔄 Coordination entre activities

**Ne fait PAS :**

- Appels directs aux APIs externes
- Requêtes SQL directes
- Calculs métier complexes

### 3. Activities Layer

**Rôle** : Logique métier et coordination

**Responsabilités :**

- 🎯 Logique métier de haut niveau
- 🔄 Coordination entre queries et repository
- 📊 Déduplication et validation des données
- 📈 Calcul des statistiques

**Ne fait PAS :**

- Gestion d'état de workflow
- Parsing bas niveau des APIs
- Optimisations SQL complexes

### 4. Queries Layer

**Rôle** : Interface avec WarcraftLogs API

**Responsabilités :**

- 📞 Requêtes GraphQL vers WarcraftLogs
- 🔄 Parsing des réponses JSON
- ⚠️ Gestion des erreurs API
- 🎯 Déduplication par joueur

**Ne fait PAS :**

- Stockage en base de données
- Logique métier complexe
- Gestion de la concurrence

### 5. Repository Layer

**Rôle** : Persistance et requêtes de données

**Responsabilités :**

- 💾 Opérations CRUD sur les tables
- 📊 Calculs de métriques en SQL
- 🚀 Optimisations de performance
- 🔄 Gestion des transactions

**Ne fait PAS :**

- Appels vers des APIs externes
- Logique de déduplication métier
- Orchestration de workflows

## 🛠️ Choix techniques justifiés

### Pourquoi Temporal ?

| Besoin            | Solution Temporal                    | Alternative écartée |
| ----------------- | ------------------------------------ | ------------------- |
| **Fiabilité**     | Retry automatique, state persistence | Cron jobs simples   |
| **Observabilité** | UI intégrée, logs structurés         | Logs custom         |
| **Scalabilité**   | Worker pools, rate limiting          | Threading manuel    |
| **Durabilité**    | Workflows long-running               | Jobs éphémères      |

### Pourquoi cette séparation en couches ?

| Avantage            | Implémentation                        |
| ------------------- | ------------------------------------- |
| **Testabilité**     | Chaque couche mockable indépendamment |
| **Maintenabilité**  | Changements isolés par responsabilité |
| **Réutilisabilité** | Queries/Repository réutilisables      |
| **Évolutivité**     | Nouvelles features sans refactor      |

### Pourquoi GORM + SQL brut ?

| Aspect                 | GORM | SQL Brut | Usage                        |
| ---------------------- | ---- | -------- | ---------------------------- |
| **CRUD simple**        | ✅   | ❌       | Insertions, updates basiques |
| **Requêtes complexes** | ❌   | ✅       | Métriques, agrégations       |
| **Performance**        | ❌   | ✅       | Requêtes critiques           |
| **Lisibilité**         | ✅   | ❌       | Code métier simple           |

## 🔧 Patterns appliqués

### 1. Repository Pattern

```go
type PlayerRankingsRepository interface {
    StoreRankingsByBatches(ctx context.Context, rankings []PlayerRanking) error
    CalculateDailySpecMetrics(ctx context.Context) error
    GetGlobalRankings(ctx context.Context) (*GlobalRankings, error)
}
```

**Avantages :**

- Abstraction de la persistence
- Facilite les tests unitaires
- Centralise les requêtes SQL

### 2. Factory Pattern

```go
func NewPlayerRankingsActivity(
    client *WarcraftLogsClientService,
    repository *PlayerRankingsRepository,
) *PlayerRankingsActivity
```

**Avantages :**

- Injection de dépendances claire
- Configuration centralisée
- Facilite les mocks

### 3. Saga Pattern (via Temporal)

```go
// Workflow = Saga coordinator
func (w *PlayerRankingsWorkflow) Execute(ctx workflow.Context, params WorkflowParams) {
    // Step 1: Fetch data
    err := workflow.ExecuteActivity(ctx, FetchAllDungeonRankings, ...)
    if err != nil {
        return err // Temporal handles compensation
    }

    // Step 2: Calculate metrics
    err = workflow.ExecuteActivity(ctx, CalculateDailyMetrics, ...)
    // ...
}
```

**Avantages :**

- Compensation automatique en cas d'échec
- État persistent entre étapes
- Retry granulaire par étape

## 🏢 Structure de packages

### Principe d'organisation

```
player_rankings/
├── queries/           # Port vers l'extérieur (WarcraftLogs)
├── repository/        # Port vers l'intérieur (Database)
├── temporal/
│   ├── activities/    # Logique métier orchestrée
│   ├── workflows/     # Orchestration pure
│   └── scheduler/     # Configuration temporelle
```

### Règles de dépendances

```
Scheduler ──► Workflow ──► Activities ──► Repository
                    │                      │
                    └──► Queries           │
                                          │
                                          ▼
                                     Database
```

**Règles :**

- ✅ **Couches supérieures** peuvent dépendre des **couches inférieures**
- ❌ **Couches inférieures** ne dépendent JAMAIS des **couches supérieures**
- ✅ **Queries et Repository** sont au même niveau (pas de dépendances croisées)

### Interfaces et contrats

```go
// Activities définissent des contrats clairs
type PlayerRankingsActivity interface {
    FetchAllDungeonRankings(ctx context.Context, dungeonIDs []int, ...) (*RankingsStats, error)
    CalculateDailyMetrics(ctx context.Context) error
}

// Repository abstrait la persistence
type PlayerRankingsRepository interface {
    StoreRankingsByBatches(ctx context.Context, rankings []PlayerRanking) error
    CalculateDailySpecMetrics(ctx context.Context) error
}
```

## 🚀 Scalabilité et performance

### Gestion de la concurrence

```go
// Semaphore pour limiter la charge sur WarcraftLogs
sem := make(chan struct{}, maxConcurrency) // 3 requêtes parallèles max

// Traitement par batch pour éviter les gros INSERT
const batchSize = 100
```

### Optimisations SQL

```sql
-- Utilisation de CTE pour les calculs complexes
WITH ranked_players AS (
    SELECT *, ROW_NUMBER() OVER (...) as player_rank
    FROM player_rankings
)
SELECT spec, AVG(score)
FROM ranked_players
WHERE player_rank <= 10
GROUP BY spec;
```

### Gestion de la mémoire

```go
// Streaming des résultats au lieu de charger tout en mémoire
func (r *Repository) StoreRankingsByBatches(rankings []PlayerRanking) error {
    for i := 0; i < len(rankings); i += batchSize {
        batch := rankings[i:min(i+batchSize, len(rankings))]
        // Process batch
    }
}
```

## 🔒 Considérations de sécurité

### API Rate Limiting

- **Concurrence limitée** : 3 requêtes parallèles max vers WarcraftLogs
- **Timeout configuré** : 30s par requête API
- **Retry avec backoff** : Évite le spam en cas d'erreur

### Gestion des erreurs

```go
// Errors contextualisées pour le debugging
return temporal.NewApplicationError(
    fmt.Sprintf("Failed to fetch dungeon %d: %v", dungeonID, err),
    "FETCH_ERROR",
)
```

### Isolation des données

- **Transactions SQL** : Atomicité des calculs de métriques
- **Validation des inputs** : Vérification des paramètres de workflow
- **Logs sans données sensibles** : Pas de tokens/credentials en logs

---

Cette architecture garantit la **maintenabilité**, **scalabilité** et **fiabilité** du système tout en gardant une complexité maîtrisée.
