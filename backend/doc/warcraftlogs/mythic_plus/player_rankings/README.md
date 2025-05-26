# Player Rankings - Documentation

## 📋 Vue d'ensemble

La feature **Player Rankings** collecte, traite et analyse automatiquement les données de performance des joueurs World of Warcraft dans les donjons Mythic+. Elle s'appuie sur l'API WarcraftLogs et utilise Temporal pour orchestrer le traitement de manière fiable et scalable.

## 🎯 Objectifs

### Problématiques résolues

- **Volume de données** : Traitement de dizaines de milliers de classements en parallèle
- **Fiabilité** : Retry automatique en cas d'échec, gestion des timeouts API
- **Cohérence des données** : Déduplication pour éviter les scores gonflés
- **Performance** : Requêtes optimisées et calculs sur les TOP 10 joueurs
- **Observabilité** : Logs détaillés et monitoring via Temporal UI

### Données collectées

- 📊 **Classements de joueurs** par donjon et spécialisation
- 📈 **Métriques de performance** (score, durée, niveau de clé)
- 🏆 **Rankings** par rôle (Tank, Healer, DPS)
- 📅 **Métriques quotidiennes** par spécialisation

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Temporal Workflow                   │
│              (Orchestration & Retry Logic)             │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                 Activities Layer                       │
│         (Business Logic & External Calls)              │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│              Repository Layer                          │
│           (Database Operations & Queries)              │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                 Query Layer                            │
│         (External API Calls - WarcraftLogs)            │
└─────────────────────────────────────────────────────────┘
```

## 🛠️ Technologies

| Technologie    | Version | Usage                   |
| -------------- | ------- | ----------------------- |
| **Go**         | 1.21+   | Langage principal       |
| **Temporal**   | 1.20+   | Orchestration workflows |
| **GORM**       | 2.0+    | ORM base de données     |
| **PostgreSQL** | 14+     | Base de données         |
| **GraphQL**    | -       | API WarcraftLogs        |

## 📁 Structure du code

```
internal/services/warcraftlogs/mythicplus/player_rankings/
├── queries/                          # 🔍 Appels API WarcraftLogs
├── repository/                       # 💾 Opérations base de données
├── temporal/                         # ⏰ Orchestration Temporal
│   ├── activities/                   # 🎯 Logique métier
│   ├── scheduler/                    # 📅 Gestion des schedules
│   └── workflows/                    # 🔄 Orchestration
```

## 🔄 Fonctionnement

### Exécution quotidienne automatique

1. **12h00 UTC** : Déclenchement du schedule Temporal
2. **Récupération** : Collecte des classements depuis WarcraftLogs (8 donjons)
3. **Traitement** : Déduplication et stockage en base de données
4. **Calcul** : Génération des métriques quotidiennes par spécialisation
5. **Finalisation** : Mise à jour des vues et invalidation des caches

### Données traitées

- **8 donjons Mythic+** de la saison actuelle
- **2 pages par donjon** (~100 joueurs top par donjon)
- **39 spécialisations** WoW différentes
- **~45 000 classements** traités quotidiennement

## 📊 Résultats

### Tables principales

- `player_rankings` : Classements bruts des joueurs
- `daily_spec_metrics_mythic_plus` : Métriques calculées par spécialisation

### Vues disponibles

- `spec_global_score_averages` : Scores moyens par spécialisation
- Rankings par rôle (Tank, Healer, DPS)

## 📖 Documentation détaillée

| Document                                   | Description                  |
| ------------------------------------------ | ---------------------------- |
| **[Architecture](architecture.md)**        | Choix techniques et patterns |
| **[Temporal Guide](temporal-guide.md)**    | Concepts et bonnes pratiques |
| **[Queries](components/queries.md)**       | API WarcraftLogs et GraphQL  |
| **[Repository](components/repository.md)** | Opérations base de données   |
| **[Activities](components/activities.md)** | Logique métier Temporal      |
| **[Workflows](components/workflows.md)**   | Orchestration Temporal       |
| **[Scheduler](components/scheduler.md)**   | Gestion des schedules        |
| **[Data Flow](data-flow.md)**              | Flux de données détaillé     |
| **[Debugging](debugging.md)**              | Guide de debugging           |
| **[Glossaire](glossary.md)**               | Termes techniques et métier  |

## 🚀 Quick Start Développeur

### 1. Comprendre le flux

```bash
# 1. Lire la vue d'ensemble
cat architecture.md

# 2. Comprendre Temporal
cat temporal-guide.md

# 3. Examiner un composant
cat components/activities.md
```

### 2. Examiner les données

```sql
-- Vérifier les dernières données
SELECT COUNT(*) FROM player_rankings;

-- Voir les métriques récentes
SELECT * FROM daily_spec_metrics_mythic_plus
WHERE capture_date = CURRENT_DATE
LIMIT 10;
```

### 3. Observer un workflow

- Aller sur **Temporal UI** (http://localhost:8080)
- Chercher `PlayerRankingsWorkflow`
- Examiner les executions récentes

## 🐛 Support

### Logs principaux

```bash
# Logs du worker Temporal
tail -f logs/temporal-worker.log

# Logs de l'application
tail -f logs/app.log | grep "player_rankings"
```

### Monitoring

- **Temporal UI** : État des workflows et activities
- **Base de données** : Requêtes de diagnostic dans [debugging.md](debugging.md)
- **Métriques** : Vérification de cohérence des scores

## ⚠️ Points d'attention

### Données sensibles

- **Déduplication critique** : Un joueur = un score max par donjon
- **TOP 10 only** : Métriques calculées sur les 10 meilleurs joueurs
- **Région CN exclue** : Serveurs chinois exclus des calculs

### Performance

- **Concurrence limitée** : Max 3 requêtes parallèles vers WarcraftLogs
- **Batch SQL** : Insertions par lots de 100
- **Timeout généreux** : 24h pour le workflow complet

---

**Pour toute question, consulter la documentation détaillée ou contacter l'équipe de développement.**
