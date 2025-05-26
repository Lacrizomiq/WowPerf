# Guide Temporal - Player Rankings

## 📖 Introduction à Temporal

Temporal est un moteur d'orchestration de workflows qui garantit l'exécution fiable de processus métier complexes. Dans Player Rankings, il orchestre la récupération et le traitement quotidien des données WarcraftLogs.

## 🎭 Concepts fondamentaux

### Workflow vs Activity

| Aspect           | Workflow                        | Activity                       |
| ---------------- | ------------------------------- | ------------------------------ |
| **Rôle**         | Orchestrateur, chef d'orchestre | Exécutant, musicien            |
| **État**         | Persistant, survit aux crashes  | Éphémère, peut être restarté   |
| **Déterminisme** | DOIT être déterministe          | Peut être non-déterministe     |
| **Durée**        | Long-running (heures/jours)     | Court-running (minutes/heures) |
| **Code**         | Logique de flux                 | Logique métier                 |

### Exemple concret Player Rankings

```go
// ✅ WORKFLOW - Orchestration déterministe
func (w *PlayerRankingsWorkflow) Execute(ctx workflow.Context, params WorkflowParams) {
    // 1. Orchestrer la récupération
    err := workflow.ExecuteActivity(ctx, "FetchAllDungeonRankings", dungeonIDs)
    if err != nil {
        return err // Temporal gère les retry automatiquement
    }

    // 2. Orchestrer le calcul
    err = workflow.ExecuteActivity(ctx, "CalculateDailyMetrics")
    return err
}

// ✅ ACTIVITY - Logique métier non-déterministe
func (a *PlayerRankingsActivity) FetchAllDungeonRankings(ctx context.Context, dungeonIDs []int) {
    // Appels API, calculs, I/O - peut varier entre exécutions
    for _, dungeonID := range dungeonIDs {
        data, err := api.GetDungeonData(dungeonID) // Non-déterministe
        // ...
    }
}
```

## ⚙️ Configuration appliquée

### Retry Policies

```go
RetryPolicy: &temporal.RetryPolicy{
    InitialInterval:    5 * time.Second,     // Premier retry après 5s
    BackoffCoefficient: 2.0,                 // Exponentiel: 5s, 10s, 20s...
    MaximumInterval:    5 * time.Minute,     // Plafonné à 5min
    MaximumAttempts:    3,                   // 3 tentatives maximum
}
```

**Justification :**

- **5s initial** : Assez rapide pour les erreurs temporaires API
- **Backoff x2** : Évite de surcharger WarcraftLogs en cas de problème
- **5min max** : Évite d'attendre trop longtemps
- **3 attempts** : Balance entre persistance et échec rapide

### Timeouts configurés

```go
ActivityOptions: workflow.ActivityOptions{
    StartToCloseTimeout: 24 * time.Hour,     // Temps max pour l'activity complète
    HeartbeatTimeout:    20 * time.Minute,   // Détection de crash activity
    ScheduleToCloseTimeout: 25 * time.Hour,  // Temps max en queue + exécution
}
```

**Justification :**

- **24h total** : Récupération de 45k rankings peut prendre du temps
- **20min heartbeat** : Détection rapide si le worker crash
- **25h schedule** : Inclut le temps d'attente en queue

## 🔄 Schedules et exécution

### Configuration du schedule quotidien

```go
scheduleOptions := client.ScheduleOptions{
    ID: "player-rankings-daily",
    Spec: client.ScheduleSpec{
        CronExpressions: []string{"0 12 * * *"}, // 12h UTC quotidien
    },
    Action: &client.ScheduleWorkflowAction{
        Workflow:  "PlayerRankingsWorkflow",
        TaskQueue: "warcraft-logs-sync",
        Args:      []interface{}{params},
    },
}
```

**Pourquoi 12h UTC ?**

- **Après le reset quotidien** WoW (11h UTC)
- **Avant le pic d'activité** européen (17h-22h UTC)
- **Données fraîches** : Les meilleurs joueurs ont déjà joué

### Gestion des WorkflowID uniques

```go
workflowID := fmt.Sprintf("player-rankings-%s", time.Now().UTC().Format("2006-01-02"))
// Résultat: "player-rankings-2024-01-15"
```

**Avantages :**

- ✅ **Un workflow par jour** : Évite les exécutions multiples
- ✅ **Idempotent** : Re-trigger le même jour = même workflow
- ✅ **Historique clair** : Facile de retrouver l'exécution d'une date

## 🎯 Activities détaillées

### 1. FetchAllDungeonRankings Activity

```go
func (a *PlayerRankingsActivity) FetchAllDungeonRankings(
    ctx context.Context,
    dungeonIDs []int,
    pagesPerDungeon int,
    maxConcurrency int,
) (*models.RankingsStats, error)
```

**Caractéristiques :**

- **Long-running** : 15-30 minutes d'exécution
- **CPU/Network intensive** : Nombreux appels API
- **Heartbeat activé** : Signale la progression

**Pattern heartbeat :**

```go
for _, dungeonID := range dungeonIDs {
    // Signaler la progression à Temporal
    activity.RecordHeartbeat(ctx, fmt.Sprintf("Processing dungeon %d", dungeonID))

    // Traitement du donjon...
}
```

**Gestion de la concurrence :**

```go
sem := make(chan struct{}, maxConcurrency) // Semaphore à 3
for _, dungeonID := range dungeonIDs {
    go func(dID int) {
        sem <- struct{}{}        // Acquérir
        defer func() { <-sem }() // Libérer
        // Traitement...
    }(dungeonID)
}
```

### 2. CalculateDailyMetrics Activity

```go
func (a *PlayerRankingsActivity) CalculateDailyMetrics(ctx context.Context) error
```

**Caractéristiques :**

- **CPU intensive** : Calculs SQL complexes
- **Short-running** : 2-5 minutes typiquement
- **Transactionnel** : Tout ou rien

**Pattern transactionnel :**

```go
return r.db.Transaction(func(tx *gorm.DB) error {
    // 1. Supprimer les métriques existantes
    if err := tx.Exec("DELETE FROM daily_spec_metrics...").Error; err != nil {
        return err // Rollback automatique
    }

    // 2. Calculer les nouvelles métriques
    if err := tx.CreateInBatches(metrics, 100).Error; err != nil {
        return err // Rollback automatique
    }

    return nil // Commit si tout OK
})
```

## 🛡️ Bonnes pratiques appliquées

### ✅ DO - Ce qu'on fait bien

#### 1. Déterminisme des workflows

```go
// ✅ Utilise workflow.Now() au lieu de time.Now()
result.StartTime = workflow.Now(ctx)

// ✅ Génération d'UUID en paramètre, pas dans le workflow
params.BatchID = fmt.Sprintf("player-rankings-%s", uuid.New().String())
```

#### 2. Gestion d'erreurs structurée

```go
// ✅ Erreurs typées pour Temporal
return temporal.NewApplicationError(
    fmt.Sprintf("Failed to fetch rankings: %v", err),
    "FETCH_ERROR", // Type d'erreur pour les retry policies
)
```

#### 3. Logs structurés

```go
// ✅ Logs avec contexte
logger.Info("Starting rankings fetch",
    "dungeonCount", len(dungeonIDs),
    "maxConcurrency", maxConcurrency,
    "batchID", params.BatchID)
```

#### 4. Activities idempotentes

```go
// ✅ Suppression puis insertion = idempotent
func (a *PlayerRankingsActivity) FetchAllDungeonRankings(...) {
    // Toujours supprimer les données existantes d'abord
    if err := a.repository.DeleteExistingRankings(ctx); err != nil {
        return err
    }
    // Puis insérer les nouvelles données
    return a.repository.StoreRankingsByBatches(ctx, rankings)
}
```

### ❌ DON'T - Ce qu'il faut éviter

#### 1. Variables globales dans workflows

```go
// ❌ État global mutable = non-déterministe
var globalCounter int

func (w *Workflow) Execute(ctx workflow.Context) {
    globalCounter++ // Différent à chaque replay !
}
```

#### 2. Appels directs depuis workflows

```go
// ❌ I/O direct dans workflow = non-déterministe
func (w *Workflow) Execute(ctx workflow.Context) {
    data, err := http.Get("https://api.example.com") // NON !
    // Utiliser une activity à la place
}
```

#### 3. Sleep non-déterministe

```go
// ❌ time.Sleep() = non-déterministe
time.Sleep(5 * time.Second)

// ✅ workflow.Sleep() = déterministe
workflow.Sleep(ctx, 5 * time.Second)
```

#### 4. Logs sans contexte

```go
// ❌ Log peu utile
log.Println("Error occurred")

// ✅ Log avec contexte
logger.Error("Failed to fetch dungeon rankings",
    "dungeonID", dungeonID,
    "error", err,
    "retryAttempt", attempt)
```

## 🔍 Monitoring et observabilité

### Temporal UI - Points de contrôle

#### Status des workflows

- ✅ **Completed** : Exécution réussie
- ❌ **Failed** : Échec après tous les retry
- 🔄 **Running** : En cours d'exécution
- ⏸️ **Terminated** : Arrêt manuel

#### Métriques à surveiller

```
Workflow Duration: 25m (normal: 20-40m)
Activities Started: 2 (FetchAll + CalculateMetrics)
Activities Completed: 2
Activities Failed: 0
Retry Attempts: 0 (idéal)
```

### Logs applicatifs

```go
// Pattern de logs appliqué
logger.Info("Activity completed",
    "activity", "FetchAllDungeonRankings",
    "duration", time.Since(startTime),
    "rankingsCount", len(rankings),
    "tanksCount", tankCount,
    "healersCount", healerCount,
    "dpsCount", dpsCount)
```

### Alerting recommandé

```yaml
alerts:
  - name: "PlayerRankings Workflow Failed"
    condition: "workflow_status == 'FAILED'"
    severity: "HIGH"

  - name: "PlayerRankings Duration Exceeded"
    condition: "workflow_duration > 2h"
    severity: "MEDIUM"

  - name: "Zero Rankings Retrieved"
    condition: "rankings_count == 0"
    severity: "HIGH"
```

## 🔧 Configuration avancée

### Task Queues

```go
// Task queue dédié pour isolation
TaskQueue: "warcraft-logs-sync"
```

**Avantages :**

- **Isolation** : Séparé des autres workflows
- **Scaling** : Workers dédiés ajustables
- **Priorité** : Contrôle de la charge

### Worker Configuration

```go
worker.Options{
    MaxConcurrentActivityExecutions: 10,  // Max 10 activities en parallèle
    MaxConcurrentWorkflowExecutions: 5,   // Max 5 workflows en parallèle
    EnableLoggingInReplay: false,         // Évite les logs dupliqués
}
```

### Schedule Policies

```go
SchedulePolicy: &client.SchedulePolicy{
    Overlap: enums.SCHEDULE_OVERLAP_POLICY_SKIP, // Skip si déjà en cours
    CatchupWindow: time.Hour,                     // Rattrapage jusqu'à 1h
}
```

**SKIP Policy :** Si un workflow est déjà en cours, le suivant est ignoré. Évite les exécutions multiples.

---

Cette configuration Temporal garantit une **exécution fiable** et **observable** du traitement quotidien des données Player Rankings.
