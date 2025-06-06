# Scheduler Component - Player Rankings

## 📖 Vue d'ensemble

Le **Scheduler** gère l'exécution automatique quotidienne du workflow Player Rankings via les **Temporal Schedules**. Il configure les horaires, les retry policies et fournit des méthodes de contrôle manuel.

## 📁 Structure

```
scheduler/
├── init.go       # Initialisation du schedule
├── options.go    # Configuration et validation
└── scheduler.go  # Gestionnaire de schedules Temporal
```

## 🎯 Responsabilités

### ✅ Ce que fait ce composant

- ⏰ **Planification quotidienne** du workflow (12h UTC)
- 🔧 **Configuration** des retry policies et timeouts
- 🎮 **Contrôle manuel** (trigger, pause, resume)
- 🧹 **Nettoyage** des workflows en cours
- ✅ **Validation** des paramètres de configuration

### ❌ Ce qu'il ne fait PAS

- Logique métier du workflow
- Gestion des données
- Appels vers des APIs externes
- Calculs de métriques

## ⏰ Configuration du schedule

### Schedule quotidien

```go
const playerRankingsScheduleID = "player-rankings-daily"

scheduleOptions := client.ScheduleOptions{
    ID: playerRankingsScheduleID,
    Spec: client.ScheduleSpec{
        CronExpressions: []string{"0 12 * * *"}, // 12h UTC quotidien
    },
    Action: &client.ScheduleWorkflowAction{
        ID:        workflowID,                    // "player-rankings-2024-01-15"
        Workflow:  "PlayerRankingsWorkflow",      // Nom du workflow
        TaskQueue: "warcraft-logs-sync",          // Queue dédiée
        Args:      []interface{}{params},         // Paramètres du workflow
    },
}
```

### Pourquoi 12h UTC ?

- **Après le reset quotidien** WoW (11h UTC)
- **Avant le pic d'activité** européen (17h-22h UTC)
- **Données fraîches** : Les meilleurs joueurs ont déjà joué

### WorkflowID unique par jour

```go
workflowID := fmt.Sprintf("player-rankings-%s", time.Now().UTC().Format("2006-01-02"))
// Résultat: "player-rankings-2024-01-15"
```

**Avantages :**

- ✅ **Un workflow par jour** : Évite les exécutions multiples
- ✅ **Idempotent** : Re-trigger le même jour = même workflow
- ✅ **Historique clair** : Facile de retrouver l'exécution d'une date

## 🔧 Configuration et options

### PlayerRankingsScheduleConfig

```go
type PlayerRankingsScheduleConfig struct {
    Hour      int    // Heure au format 24h UTC (12)
    Minute    int    // Minute (0)
    TaskQueue string // Nom de la queue ("warcraft-logs-sync")
}

var DefaultPlayerRankingsScheduleConfig = PlayerRankingsScheduleConfig{
    Hour:      12,
    Minute:    0,
    TaskQueue: "warcraft-logs-sync",
}
```

### ScheduleOptions

```go
type ScheduleOptions struct {
    Retry   RetryPolicy   // Politique de retry
    Timeout time.Duration // Timeout d'exécution (6h)
    Paused  bool         // Démarre en pause ou actif
}

func DefaultScheduleOptions() *ScheduleOptions {
    return &ScheduleOptions{
        Retry: RetryPolicy{
            InitialInterval:    time.Minute,     // 1min
            BackoffCoefficient: 2.0,             // Exponentiel
            MaximumInterval:    time.Hour,       // Cap à 1h
            MaximumAttempts:    5,               // 5 tentatives
        },
        Timeout: 6 * time.Hour,  // Temps suffisant
        Paused:  false,          // Actif par défaut
    }
}
```

### Validation de configuration

```go
func ValidateScheduleConfig(config *PlayerRankingsScheduleConfig) error {
    if config.Hour < 0 || config.Hour > 23 {
        return fmt.Errorf("invalid hour: %d", config.Hour)
    }
    if config.Minute < 0 || config.Minute > 59 {
        return fmt.Errorf("invalid minute: %d", config.Minute)
    }
    if config.TaskQueue == "" {
        return fmt.Errorf("task queue cannot be empty")
    }
    return nil
}
```

## 🎮 Contrôles manuels

### PlayerRankingsScheduleManager

```go
type PlayerRankingsScheduleManager struct {
    client                 client.Client
    logger                 *log.Logger
    playerRankingsSchedule client.ScheduleHandle
}
```

### Opérations de contrôle

#### 1. Déclenchement manuel

```go
func (sm *PlayerRankingsScheduleManager) TriggerPlayerRankingsNow(ctx context.Context) error {
    if sm.playerRankingsSchedule == nil {
        sm.playerRankingsSchedule = sm.client.ScheduleClient().GetHandle(ctx, playerRankingsScheduleID)
    }

    return sm.playerRankingsSchedule.Trigger(ctx, client.ScheduleTriggerOptions{})
}
```

#### 2. Pause/Resume

```go
func (sm *PlayerRankingsScheduleManager) PausePlayerRankingsSchedule(ctx context.Context) error {
    return sm.playerRankingsSchedule.Pause(ctx, client.SchedulePauseOptions{})
}

func (sm *PlayerRankingsScheduleManager) UnpausePlayerRankingsSchedule(ctx context.Context) error {
    return sm.playerRankingsSchedule.Unpause(ctx, client.ScheduleUnpauseOptions{})
}
```

#### 3. Suppression

```go
func (sm *PlayerRankingsScheduleManager) DeletePlayerRankingsSchedule(ctx context.Context) error {
    err := sm.playerRankingsSchedule.Delete(ctx)
    if err == nil {
        sm.playerRankingsSchedule = nil
        sm.logger.Printf("[INFO] Deleted schedule: %s", playerRankingsScheduleID)
    }
    return err
}
```

## 🧹 Nettoyage et maintenance

### Nettoyage des workflows en cours

```go
func (sm *PlayerRankingsScheduleManager) CleanupPlayerRankingsWorkflows(ctx context.Context) error {
    query := fmt.Sprintf("WorkflowType='%s'", "PlayerRankingsWorkflow")

    resp, err := sm.client.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
        Namespace: "default",
        Query:     query,
    })

    var terminatedCount int
    for _, execution := range resp.Executions {
        if execution.Status != enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
            continue // Skip non-running workflows
        }

        workflowID := execution.Execution.WorkflowId
        runID := execution.Execution.RunId

        err := sm.client.TerminateWorkflow(ctx, workflowID, runID,
            "Cleanup during service restart")
        if err != nil {
            sm.logger.Printf("[WARN] Failed to terminate workflow %s: %v", workflowID, err)
            continue
        }

        terminatedCount++
        sm.logger.Printf("[INFO] Terminated workflow: %s", workflowID)
    }

    sm.logger.Printf("[INFO] Cleanup completed - terminated %d workflows", terminatedCount)
    return nil
}
```

### Nettoyage complet

```go
func (sm *PlayerRankingsScheduleManager) CleanupAll(ctx context.Context) error {
    // 1. Supprimer le schedule
    if err := sm.DeletePlayerRankingsSchedule(ctx); err != nil {
        sm.logger.Printf("[WARN] Failed to delete schedule: %v", err)
    }

    // 2. Nettoyer les workflows
    if err := sm.CleanupPlayerRankingsWorkflows(ctx); err != nil {
        return fmt.Errorf("failed to cleanup workflows: %w", err)
    }

    return nil
}
```

## 🚀 Initialisation

### InitPlayerRankingsSchedule

```go
func InitPlayerRankingsSchedule(
    ctx context.Context,
    scheduleManager *PlayerRankingsScheduleManager,
    configPath string,
    opts *ScheduleOptions,
    logger *log.Logger
) error {
    // 1. Charger les paramètres du workflow
    params, err := LoadPlayerRankingsParams(configPath)
    if err != nil {
        logger.Printf("[ERROR] Failed to load params: %v", err)
        return err
    }

    // 2. Créer le schedule
    if err := scheduleManager.CreatePlayerRankingsSchedule(ctx, params, nil, opts); err != nil {
        logger.Printf("[ERROR] Failed to create schedule: %v", err)
        return err
    }

    logger.Printf("[INFO] Successfully created schedule with batch ID: %s", params.BatchID)
    return nil
}
```

### Instructions pour déclenchement manuel

```go
func LogPlayerRankingsManualTriggerInstructions(logger *log.Logger) {
    logger.Printf("[INFO] To trigger manually via code:")
    logger.Printf("[INFO] - playerRankingsScheduleManager.TriggerPlayerRankingsNow(ctx)")
}
```

## 📊 Configuration avancée

### Schedule Policies

```go
SchedulePolicy: &client.SchedulePolicy{
    Overlap: enums.SCHEDULE_OVERLAP_POLICY_SKIP, // Skip si déjà en cours
    CatchupWindow: time.Hour,                     // Rattrapage jusqu'à 1h
}
```

**SKIP Policy :** Si un workflow est déjà en cours, le suivant est ignoré. Évite les exécutions multiples accidentelles.

### Task Queue isolation

```go
TaskQueue: "warcraft-logs-sync"
```

**Avantages :**

- **Isolation** : Séparé des autres workflows de l'application
- **Scaling** : Workers dédiés ajustables indépendamment
- **Priorité** : Contrôle fin de la charge de travail

## 🐛 Gestion des erreurs

### Erreurs de création de schedule

```go
handle, err := sm.client.ScheduleClient().Create(ctx, scheduleOptions)
if err != nil {
    return fmt.Errorf("failed to create player rankings schedule: %w", err)
}

sm.playerRankingsSchedule = handle
sm.logger.Printf("[INFO] Created schedule: %s (cron: %s)", scheduleID, cronExpression)
```

### Logs de diagnostic

```go
sm.logger.Printf("[INFO] Listing workflows of type: %s", workflowType)
sm.logger.Printf("[INFO] Workflow cleanup completed - terminated %d workflows", terminatedCount)
sm.logger.Printf("[WARN] Failed to terminate workflow %s: %v", workflowID, err)
```

## 📈 Monitoring du scheduler

### Métriques à surveiller

- **Schedule Status** : Active/Paused
- **Last Trigger Time** : Dernière exécution
- **Next Trigger Time** : Prochaine exécution
- **Failed Triggers** : Échecs de déclenchement

### Commandes utiles

```bash
# Vérifier le statut du schedule
temporal schedule describe --schedule-id player-rankings-daily

# Lister les executions récentes
temporal workflow list --query "WorkflowType='PlayerRankingsWorkflow'"

# Trigger manuel
temporal schedule trigger --schedule-id player-rankings-daily
```

---

Le Scheduler garantit l'**exécution automatique** et **contrôlable** du workflow tout en fournissant les outils de **maintenance** et **debugging** nécessaires.
