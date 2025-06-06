# Workflows Component - Player Rankings

## 📖 Vue d'ensemble

Le **Workflow** est le chef d'orchestre de l'ensemble du processus Player Rankings. Il définit la séquence d'exécution, gère les retry policies, les timeouts, et maintient l'état de progression de manière durable.

## 📁 Structure

```
workflows/
├── definitions/          # Interfaces et constantes
│   ├── activities.go     # Définitions des activities
│   ├── config.go         # Chargement de configuration
│   └── workflow.go       # Interface du workflow
├── models/              # Structures de données
│   ├── config.go        # Paramètres de workflow
│   └── result.go        # Résultats et statistiques
└── player_rankings/     # Implémentation
    └── player_rankings_workflow.go
```

## 🎯 Responsabilités

### ✅ Ce que fait ce composant

- 🎭 **Orchestration** de la séquence d'activities
- ⏱️ **Gestion des timeouts** et retry policies
- 📊 **Tracking de progression** et métriques
- 🔄 **Gestion d'état** persistant (survit aux crashes)
- 📝 **Logging structuré** de l'exécution

### ❌ Ce qu'il ne fait PAS

- Logique métier complexe (rôle des Activities)
- Appels directs aux APIs (non-déterministe)
- Requêtes SQL directes (non-déterministe)
- Calculs intensifs (rôle des Activities)

## 🔄 Séquence d'exécution

### Vue d'ensemble du flux

```go
func (w *PlayerRankingsWorkflow) Execute(ctx workflow.Context, params WorkflowParams) (*WorkflowResult, error) {
    // 1. 🎬 Initialisation
    result := &WorkflowResult{StartTime: workflow.Now(ctx), BatchID: params.BatchID}

    // 2. ⚙️ Configuration des activities
    ctx = workflow.WithActivityOptions(ctx, activityOptions)

    // 3. 📥 Étape 1: Récupération et stockage (15-30min)
    err := workflow.ExecuteActivity(ctx, "FetchAllDungeonRankings", ...).Get(ctx, &stats)

    // 4. 📊 Étape 2: Calcul des métriques (2-5min)
    err = workflow.ExecuteActivity(ctx, "CalculateDailyMetrics").Get(ctx, nil)

    // 5. ✅ Finalisation
    result.EndTime = workflow.Now(ctx)
    return result, nil
}
```

## ⚙️ Configuration des Activities

### Activity Options

```go
activityOptions := workflow.ActivityOptions{
    StartToCloseTimeout: 24 * time.Hour,     // Temps max pour une activity
    RetryPolicy: &temporal.RetryPolicy{
        InitialInterval:    5 * time.Second,      // Délai initial
        BackoffCoefficient: 2.0,                  // Exponentiel
        MaximumInterval:    5 * time.Minute,      // Cap à 5min
        MaximumAttempts:    3,                    // 3 attempts
    },
    HeartbeatTimeout: 20 * time.Minute,      // Détection de crash
}
```

## 🎯 Déterminisme et bonnes pratiques

### ✅ Pratiques respectées

#### 1. Temps déterministe

```go
// ✅ Utilise workflow.Now() au lieu de time.Now()
result.StartTime = workflow.Now(ctx)
result.EndTime = workflow.Now(ctx)
```

#### 2. Gestion d'erreurs structurée

```go
if err != nil {
    logger.Error("Failed to fetch rankings", "error", err)
    result.Error = err.Error()
    return result, err // Retourne résultat partiel + erreur
}
```

### ❌ Anti-patterns évités

```go
// ❌ Interdit - non-déterministe
data, err := http.Get("https://api.example.com")
time.Sleep(5 * time.Second)
workflowID := uuid.New().String()

// ✅ Correct - via activities et paramètres
err := workflow.ExecuteActivity(ctx, "FetchDataActivity").Get(ctx, &data)
workflow.Sleep(ctx, 5 * time.Second)
workflowID := params.BatchID // Passé en paramètre
```

## 📊 Métriques collectées

```go
type WorkflowResult struct {
    TotalDuration  time.Duration // ~25-40 minutes
    FetchDuration  time.Duration // ~20-30 minutes
    MetricDuration time.Duration // ~2-5 minutes
    RankingsCount  int          // ~45 000 rankings
    TankCount      int          // ~3 000 (7%)
    HealerCount    int          // ~5 000 (12%)
    DPSCount       int          // ~37 000 (81%)
}
```

---
