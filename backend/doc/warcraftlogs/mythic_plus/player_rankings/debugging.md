# Debugging Guide - Player Rankings

## 📖 Vue d'ensemble

Guide pratique pour diagnostiquer et résoudre les problèmes du système Player Rankings. Couvre les outils, requêtes SQL de diagnostic, et procédures de dépannage.

## 🔍 Points de contrôle principaux

### 1. État du Workflow Temporal

### 2. Qualité des données en base

### 3. Performance des requêtes

### 4. Logs applicatifs

### 5. Métriques business

## 🎭 Debugging Temporal

### Temporal UI - Points de contrôle

**URL locale :** http://localhost:8080

#### Status des workflows

| Status            | Signification     | Action                      |
| ----------------- | ----------------- | --------------------------- |
| ✅ **Completed**  | Exécution réussie | ✅ Normal                   |
| ❌ **Failed**     | Échec après retry | 🔍 Voir les erreurs         |
| 🔄 **Running**    | En cours          | ⏱️ Vérifier progression     |
| ⏸️ **Terminated** | Arrêt manuel      | 🔧 Redémarrer si nécessaire |

#### Workflow à surveiller

```
Workflow ID: player-rankings-2024-01-15
Workflow Type: PlayerRankingsWorkflow
Task Queue: warcraft-logs-sync
```

#### Métriques normales

```
Duration totale: 20-40 minutes
├── FetchAllDungeonRankings: 15-30 minutes
└── CalculateDailyMetrics: 2-5 minutes

Activities Started: 2
Activities Completed: 2
Activities Failed: 0
Retry Attempts: 0 (idéal)
```

### Commandes Temporal CLI

```bash
# Lister les executions récentes
temporal workflow list --query "WorkflowType='PlayerRankingsWorkflow'" --limit 10

# Détail d'une exécution
temporal workflow show --workflow-id player-rankings-2024-01-15

# Vérifier le schedule
temporal schedule describe --schedule-id player-rankings-daily

# Trigger manuel
temporal schedule trigger --schedule-id player-rankings-daily

# Terminer un workflow en cours
temporal workflow terminate --workflow-id player-rankings-2024-01-15 --reason "Manual stop"
```

## 🗄️ Requêtes de diagnostic SQL

### 1. Vérification des doublons

```sql
-- ⚠️ CRITIQUE: Doit retourner 0 ligne
SELECT
    name,
    server_name,
    dungeon_id,
    COUNT(*) as count_entries
FROM player_rankings
GROUP BY name, server_name, dungeon_id
HAVING COUNT(*) > 1
ORDER BY count_entries DESC, name
LIMIT 10;
```

### 2. Sanity check des scores

```sql
-- Scores réalistes: MAX < 4000 généralement
SELECT
    'Global Max Score' as metric,
    MAX(total_score) as value
FROM (
    SELECT SUM(best_score) as total_score
    FROM (
        SELECT name, server_name, dungeon_id, MAX(score) as best_score
        FROM player_rankings
        GROUP BY name, server_name, dungeon_id
    ) grouped
    GROUP BY name, server_name
) totals

UNION ALL

SELECT
    'Total Rankings' as metric,
    COUNT(*)::text as value
FROM player_rankings

UNION ALL

SELECT
    'Unique Players' as metric,
    COUNT(DISTINCT CONCAT(name, '-', server_name))::text as value
FROM player_rankings;
```

### 3. Vérification des métriques TOP 10

```sql
-- Chaque spé/donjon doit avoir MAX 10 entrées dans le calcul
SELECT
    spec,
    class,
    role,
    dungeon_id,
    COUNT(*) as player_count
FROM (
    SELECT
        spec, class, role, dungeon_id,
        ROW_NUMBER() OVER (
            PARTITION BY spec, class, role, dungeon_id
            ORDER BY score DESC
        ) as rank
    FROM player_rankings
    WHERE server_region != 'CN'
) ranked
WHERE rank <= 10
GROUP BY spec, class, role, dungeon_id
HAVING COUNT(*) > 10  -- ⚠️ Ne doit pas retourner de lignes
ORDER BY player_count DESC;
```

### 4. Distribution par rôle

```sql
-- Vérification de la répartition Tank/Healer/DPS
SELECT
    role,
    COUNT(*) as count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 1) as percentage
FROM player_rankings
GROUP BY role
ORDER BY count DESC;

-- Résultats attendus:
-- DPS: ~81-83%
-- Healer: ~10-12%
-- Tank: ~6-8%
```

### 5. Couverture des spécialisations

```sql
-- Vérifier que toutes les spés sont représentées
SELECT
    class,
    spec,
    role,
    COUNT(*) as player_count,
    COUNT(DISTINCT dungeon_id) as dungeons_covered
FROM player_rankings
GROUP BY class, spec, role
ORDER BY player_count DESC;

-- Doit montrer ~39 spécialisations avec des données
```

### 6. État des métriques quotidiennes

```sql
-- Vérifier les métriques d'aujourd'hui
SELECT
    'Dungeon Metrics' as type,
    COUNT(*) as count
FROM daily_spec_metrics_mythic_plus
WHERE capture_date = CURRENT_DATE
  AND is_global = false

UNION ALL

SELECT
    'Global Metrics' as type,
    COUNT(*) as count
FROM daily_spec_metrics_mythic_plus
WHERE capture_date = CURRENT_DATE
  AND is_global = true

UNION ALL

SELECT
    'Last Update' as type,
    MAX(capture_date)::text as count
FROM daily_spec_metrics_mythic_plus;
```

## 📊 Debugging des performances

### Requêtes lentes à surveiller

#### 1. Calcul des métriques par donjon

```sql
-- Si cette requête est lente (>30s), vérifier les index
EXPLAIN ANALYZE
WITH ranked_players AS (
    SELECT
        spec, class, role, dungeon_id, score,
        ROW_NUMBER() OVER (PARTITION BY spec, class, role, dungeon_id ORDER BY score DESC) as rank
    FROM player_rankings
    WHERE server_region != 'CN'
)
SELECT spec, class, role, dungeon_id, AVG(score)
FROM ranked_players
WHERE rank <= 10
GROUP BY spec, class, role, dungeon_id;
```

#### 2. Index manquants

```sql
-- Vérifier la présence des index critiques
SELECT
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'player_rankings'
ORDER BY indexname;

-- Index attendus:
-- idx_player_rankings_metrics (spec, class, role, dungeon_id, score DESC)
-- idx_player_rankings_global (name, server_name, dungeon_id, score DESC)
-- idx_player_rankings_region (server_region)
```

#### 3. Taille des tables

```sql
-- Surveiller la croissance des tables
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size,
    pg_stat_get_tuples_inserted(c.oid) as inserts,
    pg_stat_get_tuples_updated(c.oid) as updates,
    pg_stat_get_tuples_deleted(c.oid) as deletes
FROM pg_tables pt
JOIN pg_class c ON c.relname = pt.tablename
WHERE tablename IN ('player_rankings', 'daily_spec_metrics_mythic_plus')
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

## 📝 Logs applicatifs

### Patterns de logs à surveiller

#### ✅ Logs normaux (INFO)

```bash
# Progression normale
grep "Starting rankings fetch" logs/app.log
grep "Fetch and storage completed" logs/app.log
grep "Successfully calculated daily metrics" logs/app.log

# Exemple de logs normaux:
# [INFO] Starting rankings fetch for multiple dungeons dungeonCount=8 maxConcurrency=3
# [INFO] Fetch and storage completed rankingsCount=45120 duration=28m15s
# [INFO] Successfully calculated daily metrics duration=2m30s
```

#### ⚠️ Logs d'attention (WARN)

```bash
# Erreurs récupérables
grep "GraphQL error" logs/app.log
grep "No rankings found" logs/app.log
grep "Failed to fetch dungeon" logs/app.log

# Exemple:
# [WARN] GraphQL error for class Warrior, spec Protection: Rate limit exceeded
# [WARN] No rankings found for class Evoker, spec Augmentation
```

#### 🚨 Logs d'erreur (ERROR)

```bash
# Erreurs critiques
grep "Failed to delete existing rankings" logs/app.log
grep "Failed to store rankings" logs/app.log
grep "Database connection failed" logs/app.log
grep "Activity failed after all retries" logs/app.log

# Exemple:
# [ERROR] Failed to store rankings: database connection lost
# [ERROR] Activity FetchAllDungeonRankings failed after 3 retries
```

### Commandes de monitoring des logs

```bash
# Logs en temps réel
tail -f logs/app.log | grep "player_rankings"

# Logs du worker Temporal
tail -f logs/temporal-worker.log

# Erreurs des dernières 24h
grep "ERROR" logs/app.log | grep "$(date +%Y-%m-%d)"

# Performances des activities
grep "duration=" logs/app.log | tail -10
```

## 🚨 Procédures de dépannage

### Problème 1: Workflow Failed

```bash
# 1. Identifier la cause
temporal workflow show --workflow-id player-rankings-$(date +%Y-%m-%d)

# 2. Vérifier les activities en échec
# → Regarder dans Temporal UI les détails de l'erreur

# 3. Solutions courantes:
# - API WarcraftLogs down → Attendre et re-trigger
# - Base de données indisponible → Vérifier connexion DB
# - Timeout → Augmenter les timeouts si nécessaire

# 4. Re-trigger manuel
temporal schedule trigger --schedule-id player-rankings-daily
```

### Problème 2: Données incohérentes

```sql
-- 1. Vérifier les doublons
SELECT COUNT(*) FROM (
    SELECT name, server_name, dungeon_id, COUNT(*)
    FROM player_rankings GROUP BY name, server_name, dungeon_id HAVING COUNT(*) > 1
);

-- 2. Si doublons détectés, nettoyer
DELETE FROM player_rankings
WHERE id NOT IN (
    SELECT MIN(id) FROM player_rankings
    GROUP BY name, server_name, dungeon_id
);

-- 3. Recalculer les métriques
DELETE FROM daily_spec_metrics_mythic_plus WHERE capture_date = CURRENT_DATE;
-- Puis re-trigger le workflow
```

### Problème 3: Performance dégradée

```sql
-- 1. Identifier les requêtes lentes
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
WHERE query LIKE '%player_rankings%'
ORDER BY mean_exec_time DESC LIMIT 5;

-- 2. Vérifier les index manquants
SELECT schemaname, tablename, attname, n_distinct, correlation
FROM pg_stats
WHERE tablename = 'player_rankings'
  AND attname IN ('spec', 'class', 'role', 'dungeon_id');

-- 3. Statistiques de la table
ANALYZE player_rankings;
ANALYZE daily_spec_metrics_mythic_plus;
```

### Problème 4: Zéro données récupérées

```bash
# 1. Vérifier la connectivité WarcraftLogs
curl -s "https://www.warcraftlogs.com/api/v2/client" | jq .

# 2. Vérifier les tokens d'authentification
grep "authentication" logs/app.log | tail -5

# 3. Tester une requête GraphQL manuelle
# → Utiliser Postman ou curl avec une requête simple

# 4. Vérifier les paramètres de configuration
cat config/player_rankings.yaml
```

## 📈 Alerting recommandé

### Métriques critiques à monitorer

#### Temporal

- **Workflow Status** ≠ Completed
- **Workflow Duration** > 2 heures
- **Activity Retry Count** > 2
- **Schedule Last Run** > 25 heures

#### Base de données

- **Rankings Count** < 40 000 ou > 60 000
- **Duplicate Count** > 0
- **Metrics Count** < 2 000
- **Max Score** > 5 000 (scores irréalistes)

#### Performance

- **Query Duration** > 30s pour calculs métriques
- **API Response Time** > 10s vers WarcraftLogs
- **Memory Usage** > 2GB pour le worker

### Configuration d'alertes

```yaml
# Exemple de configuration Prometheus/AlertManager
alerts:
  - alert: PlayerRankingsWorkflowFailed
    expr: temporal_workflow_status{workflow_type="PlayerRankingsWorkflow"} != 1
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "Player Rankings workflow failed"

  - alert: PlayerRankingsNoData
    expr: player_rankings_count < 40000
    for: 10m
    labels:
      severity: high
    annotations:
      summary: "Player Rankings count too low"
```

---

Ce guide fournit tous les outils nécessaires pour **diagnostiquer rapidement** et **résoudre efficacement** les problèmes du système Player Rankings.
