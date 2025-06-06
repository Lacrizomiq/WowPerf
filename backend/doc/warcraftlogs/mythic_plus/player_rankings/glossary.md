# Glossaire - Player Rankings

## 📖 Termes métier World of Warcraft

### Concepts de base

| Terme         | Définition                                              | Exemple                         |
| ------------- | ------------------------------------------------------- | ------------------------------- |
| **Mythic+**   | Système de donjons avec difficulté croissante et chrono | Donjon niveau 15 avec timer     |
| **Key Level** | Niveau de difficulté du donjon (2-35+)                  | Clé niveau 18 = +18             |
| **Score**     | Points calculés selon performance (vitesse, niveau)     | 3841.5 points                   |
| **Affixes**   | Modificateurs de difficulté hebdomadaires               | Tyrannical, Fortified, Sanguine |
| **Timer**     | Temps limite pour compléter un donjon                   | 30 minutes pour Cinderbrew      |
| **Medal**     | Médaille obtenue selon le temps                         | Gold, Silver, Bronze            |

### Classes et spécialisations

| Terme              | Définition                                | Exemple                          |
| ------------------ | ----------------------------------------- | -------------------------------- |
| **Classe**         | Archétype de personnage WoW               | Warrior, Paladin, Mage           |
| **Spécialisation** | Variant d'une classe avec rôle spécifique | Warrior Protection, Warrior Arms |
| **Role**           | Fonction dans le groupe                   | Tank, Healer, DPS                |
| **Spec**           | Abréviation de Spécialisation             | "Prot" pour Protection           |

### Rôles spécifiques

| Rôle       | Fonction                                 | Spécialisations (exemples)             |
| ---------- | ---------------------------------------- | -------------------------------------- |
| **Tank**   | Absorbe les dégâts, contrôle les ennemis | Warrior Protection, Paladin Protection |
| **Healer** | Soigne l'équipe                          | Priest Holy, Druid Restoration         |
| **DPS**    | Inflige des dégâts                       | Mage Fire, Rogue Assassination         |

## 🔧 Termes techniques

### Architecture

| Terme                  | Définition                                      | Usage dans Player Rankings |
| ---------------------- | ----------------------------------------------- | -------------------------- |
| **Repository Pattern** | Abstraction de la couche de données             | PlayerRankingsRepository   |
| **Factory Pattern**    | Création d'objets avec injection de dépendances | NewPlayerRankingsActivity  |
| **Saga Pattern**       | Orchestration de transactions distribuées       | Via Temporal Workflow      |
| **Semaphore**          | Contrôle de concurrence                         | Limite 3 requêtes // API   |

### Base de données

| Terme               | Définition                                    | Exemple                         |
| ------------------- | --------------------------------------------- | ------------------------------- |
| **CTE**             | Common Table Expression (sous-requête nommée) | WITH ranked_players AS (...)    |
| **Window Function** | Fonction de fenêtrage SQL                     | ROW_NUMBER() OVER (...)         |
| **Batch Insert**    | Insertion en lots pour performance            | 100 rankings par batch          |
| **UPSERT**          | INSERT ou UPDATE selon existence              | ON CONFLICT DO UPDATE           |
| **Index Composé**   | Index sur plusieurs colonnes                  | (spec, class, role, dungeon_id) |

### Performance

| Terme             | Définition                     | Valeur typique              |
| ----------------- | ------------------------------ | --------------------------- |
| **Throughput**    | Débit de traitement            | 45k rankings/jour           |
| **Latency**       | Temps de réponse               | 300ms par requête API       |
| **Concurrency**   | Nombre d'opérations parallèles | 3 requêtes simultanées      |
| **Rate Limiting** | Limitation du taux de requêtes | 5 req/sec vers WarcraftLogs |

## ⏰ Termes Temporal

### Concepts fondamentaux

| Terme          | Définition                        | Rôle dans Player Rankings    |
| -------------- | --------------------------------- | ---------------------------- |
| **Workflow**   | Orchestrateur de processus métier | PlayerRankingsWorkflow       |
| **Activity**   | Unité de travail exécutable       | FetchAllDungeonRankings      |
| **Task Queue** | File d'attente de tâches          | warcraft-logs-sync           |
| **Worker**     | Processus qui exécute les tâches  | Go process avec Temporal SDK |
| **Schedule**   | Exécution programmée récurrente   | Quotidien à 12h UTC          |

### Gestion d'état

| Terme            | Définition                          | Exemple                      |
| ---------------- | ----------------------------------- | ---------------------------- |
| **Déterminisme** | Reproductibilité des workflows      | Même input = même résultat   |
| **Replay**       | Re-exécution d'un workflow existant | Après crash du worker        |
| **Heartbeat**    | Signal de vie d'une activity        | Toutes les 30s pendant fetch |
| **Checkpoint**   | Point de sauvegarde d'état          | Après chaque activity        |

### Retry et erreurs

| Terme                | Définition                           | Configuration                   |
| -------------------- | ------------------------------------ | ------------------------------- |
| **Retry Policy**     | Stratégie de tentatives automatiques | 3 attempts, backoff exponentiel |
| **Backoff**          | Délai croissant entre retries        | 5s, 10s, 20s                    |
| **Timeout**          | Temps limite d'exécution             | 24h pour le workflow            |
| **ApplicationError** | Erreur métier typée                  | "FETCH_ERROR", "DB_ERROR"       |

## 📊 Termes métriques

### Calculs statistiques

| Terme             | Définition                 | Application                  |
| ----------------- | -------------------------- | ---------------------------- |
| **TOP 10**        | 10 meilleures performances | Base des métriques           |
| **Moyenne (AVG)** | Score moyen des TOP 10     | avg_score par spécialisation |
| **Médiane**       | Valeur centrale            | Non calculée actuellement    |
| **Percentile**    | Rang dans la distribution  | Via WarcraftLogs             |

### Types de métriques

| Terme                   | Définition                     | Exemple                      |
| ----------------------- | ------------------------------ | ---------------------------- |
| **Métrique par donjon** | Stats d'une spé dans un donjon | Warrior Prot dans Cinderbrew |
| **Métrique globale**    | Stats d'une spé tous donjons   | Warrior Prot global          |
| **Ranking rôle**        | Classement dans le rôle        | 2ème meilleur tank           |
| **Ranking global**      | Classement toutes spés         | 15ème toutes spés confondues |

## 🌐 Termes API et données

### WarcraftLogs API

| Terme          | Définition                     | Usage                      |
| -------------- | ------------------------------ | -------------------------- |
| **GraphQL**    | Langage de requête API         | Requêtes vers WarcraftLogs |
| **Endpoint**   | Point d'accès API              | /api/v2/client             |
| **Token**      | Jeton d'authentification       | Bearer token OAuth         |
| **Rate Limit** | Limite de requêtes par seconde | 5 req/sec recommandé       |
| **Pagination** | Division en pages              | Page 1, Page 2             |

### Données de response

| Terme                 | Définition                     | Exemple                    |
| --------------------- | ------------------------------ | -------------------------- |
| **Ranking**           | Classement d'un joueur         | Position + score + détails |
| **CharacterRankings** | Ensemble de classements        | Liste de rankings par page |
| **Report**            | Rapport de combat WarcraftLogs | Code rapport + fight ID    |
| **Guild**             | Guilde du joueur               | <Exemple> sur Stormrage    |
| **Server**            | Serveur/Realm du joueur        | Stormrage-EU               |

## 🗄️ Structure de base de données

### Tables principales

| Table                              | Rôle                        | Taille typique   |
| ---------------------------------- | --------------------------- | ---------------- |
| **player_rankings**                | Données brutes quotidiennes | 45k lignes/jour  |
| **daily_spec_metrics_mythic_plus** | Métriques calculées         | 2.5k lignes/jour |
| **rankings_update_states**         | État des mises à jour       | 1 ligne          |

### Colonnes importantes

| Colonne             | Type          | Signification                  |
| ------------------- | ------------- | ------------------------------ |
| **dungeon_id**      | INTEGER       | ID du donjon (1200-1207)       |
| **hard_mode_level** | INTEGER       | Niveau de la clé (2-35+)       |
| **score**           | DECIMAL(10,3) | Score du run (ex: 3841.567)    |
| **server_region**   | VARCHAR(10)   | Région (EU, US, KR, TW)        |
| **is_global**       | BOOLEAN       | Métrique globale ou par donjon |
| **capture_date**    | DATE          | Date de capture des métriques  |

## 🔧 Outils et commandes

### Temporal CLI

| Commande                       | Usage                   | Exemple                                         |
| ------------------------------ | ----------------------- | ----------------------------------------------- |
| **temporal workflow list**     | Lister les workflows    | --query "WorkflowType='PlayerRankingsWorkflow'" |
| **temporal schedule describe** | Détails du schedule     | --schedule-id player-rankings-daily             |
| **temporal schedule trigger**  | Déclencher manuellement | --schedule-id player-rankings-daily             |

### Requêtes SQL courantes

| Requête                   | Objectif          | Fréquence           |
| ------------------------- | ----------------- | ------------------- |
| **Vérification doublons** | Sanity check      | Après chaque import |
| **Scores cohérents**      | Validation métier | Quotidienne         |
| **Distribution rôles**    | Health check      | Hebdomadaire        |

## 🚨 Codes d'erreur

### Erreurs Temporal

| Code                | Signification            | Action                         |
| ------------------- | ------------------------ | ------------------------------ |
| **TIMEOUT**         | Dépassement de temps     | Augmenter timeout ou optimiser |
| **RETRY_EXHAUSTED** | Tous les retries échoués | Vérifier cause racine          |
| **WORKFLOW_FAILED** | Échec définitif          | Analyser logs et re-trigger    |

### Erreurs métier

| Code              | Signification            | Cause probable            |
| ----------------- | ------------------------ | ------------------------- |
| **FETCH_ERROR**   | Échec récupération API   | WarcraftLogs indisponible |
| **DB_ERROR**      | Problème base de données | Connexion ou contrainte   |
| **METRICS_ERROR** | Échec calcul métriques   | Données incohérentes      |

---

Ce glossaire centralise tous les termes techniques et métier pour faciliter la **compréhension** et **maintenance** du système Player Rankings.
