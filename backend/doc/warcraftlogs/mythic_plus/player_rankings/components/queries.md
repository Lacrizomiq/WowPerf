# Queries Component - Player Rankings

## 📖 Vue d'ensemble

Le composant **Queries** est responsable de l'interface avec l'API WarcraftLogs. Il gère les requêtes GraphQL, le parsing des réponses et la déduplication initiale des données.

## 📁 Structure

```
queries/
└── player_rankings.go    # Interface WarcraftLogs API
```

## 🎯 Responsabilités

### ✅ Ce que fait ce composant

- 📞 **Requêtes GraphQL** vers l'API WarcraftLogs
- 🔄 **Parsing des réponses** JSON vers les structures Go
- ⚠️ **Gestion des erreurs** API (timeouts, rate limits)
- 🎯 **Déduplication par joueur** (nom + serveur)
- 📊 **Agrégation des spécialisations** (39 specs différentes)

### ❌ Ce qu'il ne fait PAS

- Stockage en base de données
- Logique métier complexe
- Gestion de la concurrence entre donjons
- Calculs de métriques

## 🔍 Fonction principale

### GetDungeonLeaderboardByPlayer

```go
func GetDungeonLeaderboardByPlayer(
    s *service.WarcraftLogsClientService,
    params LeaderboardParams
) (*playerLeaderboardModels.DungeonLogs, error)
```

#### Paramètres d'entrée

```go
type LeaderboardParams struct {
    EncounterID  int    // ID du donjon (ex: 1200 pour Cinderbrewery)
    Page         int    // Numéro de page (1-2 typiquement)
    ServerRegion string // Région (optionnel, ex: "EU")
    ServerSlug   string // Serveur spécifique (optionnel)
    ClassName    string // Classe WoW (optionnel, ex: "Warrior")
    SpecName     string // Spécialisation (optionnel, ex: "Protection")
}
```

#### Données retournées

```go
type DungeonLogs struct {
    Page         int       // Page traitée
    HasMorePages bool      // Indicateur de pagination
    Count        int       // Nombre total de résultats
    Rankings     []Ranking // Classements des joueurs
}
```

## 🔄 Logique de traitement

### 1. Gestion des spécialisations

```go
// Si aucune spécialisation spécifiée, récupère TOUTES les spécialisations
if params.ClassName == "" && params.SpecName == "" {
    specsToFetch = Specializations // 39 spécialisations WoW
} else {
    // Sinon, seulement la spécialisation demandée
    specsToFetch = []ClassSpec{{ClassName: params.ClassName, SpecName: params.SpecName}}
}
```

**Impact :** Pour une requête sans filtre, **39 requêtes GraphQL** sont exécutées en séquence.

### 2. Requête GraphQL par spécialisation

```graphql
query getDungeonLeaderboard(
  $encounterId: Int!
  $page: Int!
  $serverRegion: String
  $serverSlug: String
  $className: String
  $specName: String
) {
  worldData {
    encounter(id: $encounterId) {
      name
      characterRankings(
        leaderboard: Any
        page: $page
        serverRegion: $serverRegion
        serverSlug: $serverSlug
        className: $className
        specName: $specName
      )
    }
  }
}
```

### 3. Déduplication des résultats

```go
// Map pour éviter les doublons par joueur
playerMap := make(map[string]playerLeaderboardModels.Ranking)

for _, ranking := range characterRankings.Rankings {
    // Clé unique : nom + serveur (évite les homonymes)
    playerKey := fmt.Sprintf("%s-%s", ranking.Name, ranking.Server.Name)

    if existing, exists := playerMap[playerKey]; exists {
        // Garder le meilleur score seulement
        if ranking.Score > existing.Score {
            playerMap[playerKey] = ranking
        }
    } else {
        playerMap[playerKey] = ranking
    }
}
```

**Pourquoi cette déduplication ?**

- Un même joueur peut apparaître dans plusieurs réponses de spécialisations
- On veut garder seulement son **meilleur score** pour ce donjon
- Évite les scores gonflés dans les calculs finaux

## 📊 Structure des données

### Ranking (joueur individuel)

```go
type Ranking struct {
    Name          string  // Nom du joueur
    Class         string  // Classe WoW (ex: "Warrior")
    Spec          string  // Spécialisation (ex: "Protection")
    Score         float64 // Score M+ (ex: 3841.2)
    Amount        int     // Points de score bruts
    HardModeLevel int     // Niveau de clé (ex: 15)
    Duration      int64   // Durée en millisecondes
    StartTime     int64   // Timestamp de début
    Medal         string  // Médaille obtenue ("gold", "silver", etc.)
    Affixes       []int   // IDs des affixes de la semaine
    BracketData   int     // Données de bracket
    Faction       int     // Faction (Alliance/Horde)

    // Données de guilde
    Guild struct {
        ID      int    // ID de la guilde
        Name    string // Nom de la guilde
        Faction int    // Faction de la guilde
    }

    // Données de serveur
    Server struct {
        ID     int    // ID du serveur
        Name   string // Nom du serveur
        Region string // Région (EU, US, etc.)
    }

    // Données de rapport WarcraftLogs
    Report struct {
        Code      string // Code du rapport (ex: "aBcDeFg123")
        FightID   int    // ID du combat dans le rapport
        StartTime int64  // Timestamp de début du rapport
    }
}
```

## 🎮 Spécialisations WoW supportées

### Liste complète (39 spécialisations)

```go
var Specializations = []workflowModels.ClassSpec{
    // Priest (3 specs)
    {ClassName: "Priest", SpecName: "Discipline"},
    {ClassName: "Priest", SpecName: "Holy"},
    {ClassName: "Priest", SpecName: "Shadow"},

    // Death Knight (3 specs)
    {ClassName: "DeathKnight", SpecName: "Blood"},
    {ClassName: "DeathKnight", SpecName: "Frost"},
    {ClassName: "DeathKnight", SpecName: "Unholy"},

    // Druid (4 specs)
    {ClassName: "Druid", SpecName: "Balance"},
    {ClassName: "Druid", SpecName: "Feral"},
    {ClassName: "Druid", SpecName: "Guardian"},
    {ClassName: "Druid", SpecName: "Restoration"},

    // ... (continue pour toutes les classes)

    // Evoker (3 specs)
    {ClassName: "Evoker", SpecName: "Devastation"},
    {ClassName: "Evoker", SpecName: "Preservation"},
    {ClassName: "Evoker", SpecName: "Augmentation"},
}
```

### Répartition par rôle

- **Tanks** : 6 spécialisations (Protection, Blood, Vengeance, Guardian, Brewmaster)
- **Healers** : 7 spécialisations (Holy, Discipline, Restoration, Mistweaver, Preservation)
- **DPS** : 26 spécialisations (toutes les autres)

## ⚠️ Gestion des erreurs

### Types d'erreurs gérées

#### 1. Erreurs GraphQL

```go
var errorResponse struct {
    Errors []struct {
        Message string `json:"message"`
    } `json:"errors"`
}

if len(errorResponse.Errors) > 0 {
    log.Printf("GraphQL error for class %s, spec %s: %s",
        spec.ClassName, spec.SpecName, errorResponse.Errors[0].Message)
    continue // Continue avec la spécialisation suivante
}
```

#### 2. Erreurs de parsing JSON

```go
if err := json.Unmarshal(response, &result); err != nil {
    return nil, fmt.Errorf("failed to unmarshal response for class %s, spec %s: %w",
        spec.ClassName, spec.SpecName, err)
}
```

#### 3. Réponses vides

```go
if characterRankings.Rankings == nil {
    log.Printf("No rankings found for class %s, spec %s", spec.ClassName, spec.SpecName)
    continue // Pas une erreur bloquante
}
```

### Stratégie de résilience

- **Continue on error** : Si une spécialisation échoue, les autres continuent
- **Logs détaillés** : Chaque erreur est loggée avec le contexte
- **Pas d'échec total** : La fonction retourne les données collectées même si certaines specs échouent

## 🚀 Performance et optimisations

### Optimisations appliquées

#### 1. Déduplication en mémoire

```go
// Map plus efficace que des boucles de recherche
playerMap := make(map[string]playerLeaderboardModels.Ranking)
```

#### 2. Early exit sur pagination

```go
if !dungeonData.HasMorePages {
    break // Arrête de requêter s'il n'y a plus de pages
}
```

#### 3. Logs avec compteurs

```go
log.Printf("Deduplicated rankings: %d unique players from %d total entries",
    len(allRankings), totalCount)
```

### Métriques typiques

- **Temps par requête** : 200-500ms par spécialisation
- **Temps total** : 39 × 0.3s = ~12s par donjon/page
- **Déduplication** : ~10-20% de doublons éliminés typiquement
- **Taille des réponses** : 50-100 rankings par spécialisation/page

## 🔧 Configuration et limites

### Limites de l'API WarcraftLogs

- **Rate limiting** : ~3-5 requêtes/seconde recommandé
- **Timeout** : 30s par requête maximum
- **Pagination** : Maximum 100 résultats par page
- **Historique** : Données disponibles sur ~2 semaines

### Usage typique dans Player Rankings

```go
// Appelé depuis FetchAllDungeonRankings Activity
params := LeaderboardParams{
    EncounterID: dungeonID,     // 1 des 8 donjons M+
    Page:        page,          // 1 ou 2
    // Pas de filtre classe/spec = récupère tout
}

dungeonData, err := GetDungeonLeaderboardByPlayer(client, params)
```

**Résultat :** Pour 8 donjons × 2 pages = **624 requêtes GraphQL** par exécution quotidienne (39 × 8 × 2).

---

Ce composant est **critique** pour la qualité des données : une mauvaise déduplication ici se répercute sur tous les calculs en aval.
