# 📊 Mythic+ Analytics API Documentation

## Vue d'ensemble

Cette API fournit des analyses statistiques complètes sur les runs Mythic+ basées sur les données collectées depuis Raider.IO. Elle permet d'analyser les spécialisations populaires, les compositions d'équipe, et les tendances par donjon, niveau de clé, et région.

**Base URL :** `/raiderio/mythicplus/analytics`  
**Cache :** 24 heures sur tous les endpoints  
**Format de réponse :** JSON

---

## 🌍 Analyses Globales

### Spécialisations par Rôle

#### `GET /specs/tank`

**Description :** Retourne les spécialisations Tank les plus utilisées globalement  
**Paramètres :** Aucun  
**Réponse :** `Array<SpecializationStats>`

#### `GET /specs/healer`

**Description :** Retourne les spécialisations Healer les plus utilisées globalement  
**Paramètres :** Aucun  
**Réponse :** `Array<SpecializationStats>`

#### `GET /specs/dps`

**Description :** Retourne les spécialisations DPS les plus utilisées globalement  
**Paramètres :** Aucun  
**Réponse :** `Array<SpecializationStats>`

#### `GET /specs/{role}`

**Description :** Retourne les spécialisations pour un rôle spécifique  
**Paramètres :**

- `role` (path) : `tank`, `healer`, ou `dps`

**Réponse :** `Array<SpecializationStats>`

### Compositions d'Équipe

#### `GET /compositions`

**Description :** Retourne les compositions d'équipe les plus populaires avec scores moyens  
**Paramètres :**

- `limit` (query, optional) : Nombre maximum de résultats (défaut: 20)
- `min_usage` (query, optional) : Nombre minimum d'utilisations pour filtrer (défaut: 5)

**Réponse :** `Array<CompositionStats>`

**Exemple :** `/compositions?limit=10&min_usage=15`

---

## 🏰 Analyses par Donjon

### Spécialisations par Donjon

#### `GET /dungeons/specs/tank`

**Description :** Retourne les spécialisations Tank populaires pour chaque donjon  
**Paramètres :**

- `top_n` (query, optional) : Nombre de spécs par donjon (0 = toutes, défaut: 0)

**Réponse :** `Array<DungeonSpecStats>`

#### `GET /dungeons/specs/healer`

**Description :** Retourne les spécialisations Healer populaires pour chaque donjon  
**Paramètres :**

- `top_n` (query, optional) : Nombre de spécs par donjon (0 = toutes, défaut: 0)

**Réponse :** `Array<DungeonSpecStats>`

#### `GET /dungeons/specs/dps`

**Description :** Retourne les spécialisations DPS populaires pour chaque donjon  
**Paramètres :**

- `top_n` (query, optional) : Nombre de spécs par donjon (0 = toutes, défaut: 0)

**Réponse :** `Array<DungeonSpecStats>`

### Analyses Spécifiques par Donjon

#### `GET /dungeons/{dungeon_slug}/specs/{role}`

**Description :** Retourne les spécialisations pour un rôle dans un donjon spécifique  
**Paramètres :**

- `dungeon_slug` (path) : Identifiant du donjon (ex: `theater-of-pain`)
- `role` (path) : `tank`, `healer`, ou `dps`
- `top_n` (query, optional) : Nombre de spécs (0 = toutes, défaut: 0)

**Réponse :** `Array<DungeonSpecStats>`

**Exemple :** `/dungeons/theater-of-pain/specs/tank?top_n=3`

### Compositions par Donjon

#### `GET /dungeons/compositions`

**Description :** Retourne les compositions populaires pour chaque donjon  
**Paramètres :**

- `top_n` (query, optional) : Nombre de compositions par donjon (0 = toutes, défaut: 0)
- `min_usage` (query, optional) : Usage minimum pour filtrer (défaut: 3)

**Réponse :** `Array<DungeonCompositionStats>`

---

## 🔥 Analyses Avancées

### Analyse par Niveau de Clé

#### `GET /key-levels`

**Description :** Retourne les spécialisations par bracket de niveau de clé  
**Brackets :**

- Very High Keys (20+)
- High Keys (18-19)
- Mid Keys (16-17)
- Other Keys (<16)

**Paramètres :**

- `min_usage` (query, optional) : Usage minimum pour filtrer (défaut: 5)

**Réponse :** `Array<KeyLevelStats>`

### Analyse par Région

#### `GET /regions`

**Description :** Retourne les spécialisations populaires par région (US, EU, KR, TW)  
**Paramètres :** Aucun  
**Réponse :** `Array<RegionStats>`

---

## 📈 Statistiques Utilitaires

### Statistiques Générales

#### `GET /stats/overall`

**Description :** Retourne les statistiques générales du dataset  
**Inclut :** Total runs, runs avec score, compositions uniques, donjons, régions, score moyen, etc.  
**Paramètres :** Aucun  
**Réponse :** `OverallStats`

### Distributions

#### `GET /stats/key-levels`

**Description :** Retourne la distribution des runs par niveau de clé mythique  
**Paramètres :** Aucun  
**Réponse :** `Array<KeyLevelDistribution>`

#### `GET /stats/dungeons`

**Description :** Retourne la distribution des runs par donjon  
**Paramètres :** Aucun  
**Réponse :** `Array<DungeonDistribution>`

#### `GET /stats/regions`

**Description :** Retourne la distribution des runs par région  
**Paramètres :** Aucun  
**Réponse :** `Array<RegionDistribution>`

---

## 📋 Structures de Données

### SpecializationStats

```json
{
  "class": "Warrior",
  "spec": "Protection",
  "display": "Warrior - Protection",
  "usage_count": 1247,
  "percentage": 23.5
}
```

### CompositionStats

```json
{
  "tank": "Warrior - Protection",
  "healer": "Priest - Holy",
  "dps1": "Mage - Fire",
  "dps2": "Hunter - Beast Mastery",
  "dps3": "Rogue - Assassination",
  "usage_count": 89,
  "percentage": 2.1,
  "avg_score": 456.7
}
```

### DungeonSpecStats

```json
{
  "dungeon_slug": "theater-of-pain",
  "dungeon_name": "Theater of Pain",
  "class": "Warrior",
  "spec": "Protection",
  "display": "Warrior - Protection",
  "usage_count": 234,
  "percentage": 18.9,
  "rank_in_dungeon": 1
}
```

### KeyLevelStats

```json
{
  "role": "Tank",
  "key_level_bracket": "Very High Keys (20+)",
  "class": "Warrior",
  "spec": "Protection",
  "display": "Warrior - Protection",
  "usage_count": 156,
  "percentage": 34.2,
  "avg_score": 478.3
}
```

---

## 💡 Exemples d'Utilisation

### 1. Méta Tank global

```bash
GET /raiderio/mythicplus/analytics/specs/tank
```

### 2. Top 10 compositions populaires

```bash
GET /raiderio/mythicplus/analytics/compositions?limit=10&min_usage=20
```

### 3. Toutes les spécs DPS pour "Theater of Pain"

```bash
GET /raiderio/mythicplus/analytics/dungeons/theater-of-pain/specs/dps
```

### 4. Top 3 compositions par donjon

```bash
GET /raiderio/mythicplus/analytics/dungeons/compositions?top_n=3&min_usage=5
```

### 5. Méta par niveau de clé élevé

```bash
GET /raiderio/mythicplus/analytics/key-levels?min_usage=10
```

### 6. Stats générales pour dashboard

```bash
GET /raiderio/mythicplus/analytics/stats/overall
```

---

## ⚡ Notes Importantes

- **Cache :** Toutes les données sont mises en cache pendant 24h
- **Paramètre `top_n = 0`** : Retourne TOUS les résultats (recommandé pour analyses complètes)
- **Paramètre `min_usage`** : Filtre les données avec peu d'utilisations pour éviter le bruit
- **Format des spécialisations** : `"Classe - Spécialisation"` pour faciliter l'affichage
- **Classement** : Les résultats sont triés par usage décroissant sauf indication contraire
- **Données source** : Basées sur les runs collectés depuis Raider.IO avec scores réels
