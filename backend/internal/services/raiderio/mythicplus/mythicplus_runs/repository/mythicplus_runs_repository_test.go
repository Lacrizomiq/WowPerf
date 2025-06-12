// go test ./internal/services/raiderio/mythicplus/mythicplus_runs/repository -v
package raiderioMythicPlusRunsRepository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	models "wowperf/internal/models/raiderio/mythicplus_runs"
)

func TestMythicPlusRunsRepository(t *testing.T) {
	t.Run("ProcessRuns - Insert new runs successfully", func(t *testing.T) {
		db := setupTestDB(t) // DB fraîche pour chaque test
		repo := NewMythicPlusRunsRepository(db)

		testRuns := []*models.Run{
			createTestRunWithUniqueComposition(123456, "Demon Hunter", "Vengeance", "Priest", "Discipline"),
			createTestRunWithUniqueComposition(789012, "Paladin", "Protection", "Shaman", "Restoration"), // Composition différente
		}

		stats, err := repo.ProcessRuns(testRuns, "test-batch-1")

		require.NoError(t, err)
		assert.Equal(t, 2, stats.NewRuns)
		assert.Equal(t, 0, stats.SkippedRuns)
		assert.Equal(t, 2, stats.NewCompositions)

		// Vérifie que les runs sont en DB
		var dbRunsCount int64
		err = db.Model(&models.MythicPlusRuns{}).Count(&dbRunsCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(2), dbRunsCount)
	})

	t.Run("ProcessRuns - Skip duplicate runs", func(t *testing.T) {
		db := setupTestDB(t) // DB fraîche
		repo := NewMythicPlusRunsRepository(db)

		// D'abord insérer une run
		firstRun := createTestRunWithUniqueComposition(123456, "Demon Hunter", "Vengeance", "Priest", "Discipline")
		_, err := repo.ProcessRuns([]*models.Run{firstRun}, "initial-batch")
		require.NoError(t, err)

		// Maintenant tester avec un duplicate et un nouveau
		testRuns := []*models.Run{
			createTestRunWithUniqueComposition(123456, "Demon Hunter", "Vengeance", "Priest", "Discipline"), // Même keystone_run_id
			createTestRunWithUniqueComposition(999999, "Paladin", "Protection", "Shaman", "Restoration"),    // Nouveau
		}

		stats, err := repo.ProcessRuns(testRuns, "test-batch-2")

		require.NoError(t, err)
		assert.Equal(t, 1, stats.NewRuns)     // Seulement le nouveau
		assert.Equal(t, 1, stats.SkippedRuns) // L'existant est skippé

		// Total runs en DB = 1 (first) + 1 (new) = 2
		var count int64
		err = db.Model(&models.MythicPlusRuns{}).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})

	t.Run("ProcessRuns - Same team composition gets reused", func(t *testing.T) {
		db := setupTestDB(t) // DB fraîche
		repo := NewMythicPlusRunsRepository(db)

		// Crée 2 runs avec la MÊME composition
		run1 := createTestRunWithUniqueComposition(111111, "Demon Hunter", "Vengeance", "Priest", "Discipline")
		run2 := createTestRunWithUniqueComposition(222222, "Demon Hunter", "Vengeance", "Priest", "Discipline") // MÊME composition

		testRuns := []*models.Run{run1, run2}
		stats, err := repo.ProcessRuns(testRuns, "test-batch-same-comp")

		require.NoError(t, err)
		assert.Equal(t, 2, stats.NewRuns)
		assert.Equal(t, 1, stats.NewCompositions)      // 1 composition créée
		assert.Equal(t, 1, stats.ExistingCompositions) // 1 réutilisée
	})

	t.Run("ProcessRuns - Same team composition gets reused", func(t *testing.T) {
		db := setupTestDB(t) // DB fraîche
		repo := NewMythicPlusRunsRepository(db)

		// Crée 2 runs avec la MÊME composition
		run1 := createTestRunWithUniqueComposition(111111, "Demon Hunter", "Vengeance", "Priest", "Discipline")
		run2 := createTestRunWithUniqueComposition(222222, "Demon Hunter", "Vengeance", "Priest", "Discipline") // MÊME composition

		testRuns := []*models.Run{run1, run2}
		stats, err := repo.ProcessRuns(testRuns, "test-batch-same-comp")

		require.NoError(t, err)
		assert.Equal(t, 2, stats.NewRuns)
		assert.Equal(t, 1, stats.NewCompositions)      // 1 composition créée
		assert.Equal(t, 1, stats.ExistingCompositions) // 1 réutilisée

		// Vérifie que les 2 runs pointent vers la même composition
		var runsWithSameComp []models.MythicPlusRuns
		err = db.Where("keystone_run_id IN ?", []int64{111111, 222222}).Find(&runsWithSameComp).Error
		require.NoError(t, err)
		assert.Len(t, runsWithSameComp, 2)
		assert.Equal(t, runsWithSameComp[0].TeamCompositionID, runsWithSameComp[1].TeamCompositionID)
	})

	t.Run("ProcessRuns - Invalid runs are skipped", func(t *testing.T) {
		db := setupTestDB(t) // DB fraîche
		repo := NewMythicPlusRunsRepository(db)

		invalidRun := createTestRunWithUniqueComposition(333333, "Demon Hunter", "Vengeance", "Priest", "Discipline")
		invalidRun.Roster = invalidRun.Roster[:3] // Seulement 3 membres au lieu de 5

		testRuns := []*models.Run{invalidRun}
		stats, err := repo.ProcessRuns(testRuns, "test-batch-invalid")

		require.NoError(t, err)
		assert.Equal(t, 0, stats.NewRuns)
		assert.Equal(t, 1, stats.SkippedRuns)
	})
}

func TestDataMappingAndInsertion(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMythicPlusRunsRepository(db)

	t.Run("Verify complete data mapping from API to DB", func(t *testing.T) {
		// Crée une run avec des données spécifiques pour vérifier le mapping
		testRun := &models.Run{
			KeystoneRunID:   999888777,
			Season:          "season-tww-2",
			Status:          "finished",
			MythicLevel:     25,
			ClearTimeMs:     1234567,
			KeystoneTimeMs:  1800000,
			CompletedAt:     time.Date(2025, 5, 27, 15, 30, 45, 0, time.UTC),
			NumChests:       3,
			TimeRemainingMs: 565433,
			Dungeon: models.DungeonInfo{
				Slug: "test-dungeon",
				Name: "Test Dungeon Name",
			},
			Roster: []models.RosterMember{
				{
					Role: "tank",
					Character: models.Character{
						Class:  models.ClassInfo{Name: "Death Knight"},
						Spec:   models.SpecInfo{Name: "Blood"},
						Region: models.RegionInfo{Slug: "us"},
					},
				},
				{
					Role: "healer",
					Character: models.Character{
						Class:  models.ClassInfo{Name: "Paladin"},
						Spec:   models.SpecInfo{Name: "Holy"},
						Region: models.RegionInfo{Slug: "us"},
					},
				},
				{
					Role: "dps",
					Character: models.Character{
						Class:  models.ClassInfo{Name: "Hunter"},
						Spec:   models.SpecInfo{Name: "Beast Mastery"},
						Region: models.RegionInfo{Slug: "us"},
					},
				},
				{
					Role: "dps",
					Character: models.Character{
						Class:  models.ClassInfo{Name: "Mage"},
						Spec:   models.SpecInfo{Name: "Fire"},
						Region: models.RegionInfo{Slug: "us"},
					},
				},
				{
					Role: "dps",
					Character: models.Character{
						Class:  models.ClassInfo{Name: "Warlock"},
						Spec:   models.SpecInfo{Name: "Destruction"},
						Region: models.RegionInfo{Slug: "us"},
					},
				},
			},
		}

		// Process la run
		stats, err := repo.ProcessRuns([]*models.Run{testRun}, "data-mapping-test")
		require.NoError(t, err)
		assert.Equal(t, 1, stats.NewRuns)

		// 1. Vérifie que la run est bien insérée avec tous les champs corrects
		var dbRun models.MythicPlusRuns
		err = db.Where("keystone_run_id = ?", 999888777).First(&dbRun).Error
		require.NoError(t, err)

		// Vérifie TOUS les champs mappés
		assert.Equal(t, int64(999888777), dbRun.KeystoneRunID)
		assert.Equal(t, "season-tww-2", dbRun.Season)
		assert.Equal(t, "us", dbRun.Region) // Extrait du premier membre du roster
		assert.Equal(t, "test-dungeon", dbRun.DungeonSlug)
		assert.Equal(t, "Test Dungeon Name", dbRun.DungeonName)
		assert.Equal(t, 25, dbRun.MythicLevel)
		assert.Equal(t, "finished", dbRun.Status)
		assert.Equal(t, int64(1234567), dbRun.ClearTimeMs)
		assert.Equal(t, int64(1800000), dbRun.KeystoneTimeMs)
		assert.Equal(t, 3, dbRun.NumChests)
		assert.Equal(t, int64(565433), dbRun.TimeRemainingMs)
		assert.NotNil(t, dbRun.TeamCompositionID)

		// Vérifie le timestamp (avec une tolérance)
		expectedTime := time.Date(2025, 5, 27, 15, 30, 45, 0, time.UTC)
		assert.True(t, dbRun.CompletedAt.Equal(expectedTime))

		// 2. Vérifie que la team composition est correctement créée
		var teamComp models.MythicPlusTeamComposition
		err = db.Where("id = ?", *dbRun.TeamCompositionID).First(&teamComp).Error
		require.NoError(t, err)

		assert.Equal(t, "Death Knight", teamComp.TankClass)
		assert.Equal(t, "Blood", teamComp.TankSpec)
		assert.Equal(t, "Paladin", teamComp.HealerClass)
		assert.Equal(t, "Holy", teamComp.HealerSpec)

		// DPS triés alphabétiquement par classe puis spec
		assert.Equal(t, "Hunter", teamComp.Dps1Class)
		assert.Equal(t, "Beast Mastery", teamComp.Dps1Spec)
		assert.Equal(t, "Mage", teamComp.Dps2Class)
		assert.Equal(t, "Fire", teamComp.Dps2Spec)
		assert.Equal(t, "Warlock", teamComp.Dps3Class)
		assert.Equal(t, "Destruction", teamComp.Dps3Spec)

		// 3. Vérifie que les 5 entrées roster sont créées
		var rosterEntries []models.MythicPlusRunRoster
		err = db.Where("team_composition_id = ?", teamComp.ID).Find(&rosterEntries).Error
		require.NoError(t, err)
		assert.Len(t, rosterEntries, 5)

		// Vérifie chaque entrée roster
		roleCount := make(map[string]int)
		classCount := make(map[string]int)
		for _, entry := range rosterEntries {
			assert.Equal(t, teamComp.ID, entry.TeamCompositionID)
			roleCount[entry.Role]++
			classCount[entry.ClassName]++
		}

		// Vérifie la répartition des rôles
		assert.Equal(t, 1, roleCount["tank"])
		assert.Equal(t, 1, roleCount["healer"])
		assert.Equal(t, 3, roleCount["dps"])

		// Vérifie que toutes les classes sont présentes
		assert.Equal(t, 1, classCount["Death Knight"])
		assert.Equal(t, 1, classCount["Paladin"])
		assert.Equal(t, 1, classCount["Hunter"])
		assert.Equal(t, 1, classCount["Mage"])
		assert.Equal(t, 1, classCount["Warlock"])
	})

	t.Run("Verify foreign key relationships", func(t *testing.T) {
		// Insère une run
		testRun := createTestRunWithUniqueComposition(123123123, "Warrior", "Protection", "Druid", "Restoration")
		_, err := repo.ProcessRuns([]*models.Run{testRun}, "fk-test")
		require.NoError(t, err)

		// Récupère la run avec ses relations
		var dbRun models.MythicPlusRuns
		err = db.Preload("TeamComposition").Preload("TeamComposition.RunRoster").
			Where("keystone_run_id = ?", 123123123).First(&dbRun).Error
		require.NoError(t, err)

		// Vérifie que les relations sont bien chargées
		assert.NotNil(t, dbRun.TeamComposition)
		assert.Equal(t, "Warrior", dbRun.TeamComposition.TankClass)
		assert.Equal(t, "Druid", dbRun.TeamComposition.HealerClass)

	})

	t.Run("Verify edge cases and data types", func(t *testing.T) {
		// Test avec des valeurs limites
		edgeCaseRun := createTestRunWithUniqueComposition(999999999, "Monk", "Brewmaster", "Evoker", "Preservation")
		edgeCaseRun.MythicLevel = 50          // Niveau très élevé
		edgeCaseRun.ClearTimeMs = 0           // Temps minimal
		edgeCaseRun.TimeRemainingMs = -500000 // Temps négatif (run échouée)
		edgeCaseRun.NumChests = 0             // Pas de coffres

		_, err := repo.ProcessRuns([]*models.Run{edgeCaseRun}, "edge-case-test")
		require.NoError(t, err)

		// Vérifie que les valeurs limites sont bien stockées
		var dbRun models.MythicPlusRuns
		err = db.Where("keystone_run_id = ?", 999999999).First(&dbRun).Error
		require.NoError(t, err)

		assert.Equal(t, 50, dbRun.MythicLevel)
		assert.Equal(t, int64(0), dbRun.ClearTimeMs)
		assert.Equal(t, int64(-500000), dbRun.TimeRemainingMs)
		assert.Equal(t, 0, dbRun.NumChests)
	})
}

func TestTeamCompositionHashing(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMythicPlusRunsRepository(db)

	// Crée 2 rosters identiques mais dans un ordre différent
	roster1 := []models.RosterMember{
		createRosterMember("tank", "Demon Hunter", "Vengeance", "eu"),
		createRosterMember("healer", "Priest", "Discipline", "eu"),
		createRosterMember("dps", "Mage", "Arcane", "eu"),
		createRosterMember("dps", "Warrior", "Arms", "eu"),
		createRosterMember("dps", "Death Knight", "Unholy", "eu"),
	}

	roster2 := []models.RosterMember{
		createRosterMember("healer", "Priest", "Discipline", "eu"),
		createRosterMember("dps", "Death Knight", "Unholy", "eu"), // Ordre différent
		createRosterMember("tank", "Demon Hunter", "Vengeance", "eu"),
		createRosterMember("dps", "Warrior", "Arms", "eu"),
		createRosterMember("dps", "Mage", "Arcane", "eu"),
	}

	// Les 2 compositions doivent produire le même hash
	tx := db.Begin()
	defer tx.Rollback()

	comp1, _, err1 := repo.getOrCreateTeamCompositionCached(tx, roster1, nil)
	require.NoError(t, err1)

	comp2, isNew, err2 := repo.getOrCreateTeamCompositionCached(tx, roster2, nil)
	require.NoError(t, err2)

	assert.Equal(t, comp1.CompositionHash, comp2.CompositionHash)
	assert.False(t, isNew) // La 2e composition n'est pas nouvelle
}

// Helper functions
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate les tables
	err = db.AutoMigrate(
		&models.MythicPlusRuns{},
		&models.MythicPlusTeamComposition{},
		&models.MythicPlusRunRoster{},
	)
	require.NoError(t, err)

	return db
}

func createTestRun(keystoneRunID int64) *models.Run {
	return &models.Run{
		KeystoneRunID:   keystoneRunID,
		Season:          "season-tww-2",
		Status:          "finished",
		MythicLevel:     20,
		ClearTimeMs:     1464036,
		KeystoneTimeMs:  1860999,
		CompletedAt:     time.Now(),
		NumChests:       2,
		TimeRemainingMs: 396963,
		Dungeon: models.DungeonInfo{
			Slug: "darkflame-cleft",
			Name: "Darkflame Cleft",
		},
		Roster: createStandardRoster(),
	}
}

func createTestRunWithUniqueComposition(keystoneRunID int64, tankClass, tankSpec, healerClass, healerSpec string) *models.Run {
	return &models.Run{
		KeystoneRunID:   keystoneRunID,
		Season:          "season-tww-2",
		Status:          "finished",
		MythicLevel:     20,
		ClearTimeMs:     1464036,
		KeystoneTimeMs:  1860999,
		CompletedAt:     time.Now(),
		NumChests:       2,
		TimeRemainingMs: 396963,
		Dungeon: models.DungeonInfo{
			Slug: "darkflame-cleft",
			Name: "Darkflame Cleft",
		},
		Roster: []models.RosterMember{
			createRosterMember("tank", tankClass, tankSpec, "eu"),
			createRosterMember("healer", healerClass, healerSpec, "eu"),
			createRosterMember("dps", "Mage", "Arcane", "eu"),
			createRosterMember("dps", "Warrior", "Arms", "eu"),
			createRosterMember("dps", "Death Knight", "Unholy", "eu"),
		},
	}
}

func createStandardRoster() []models.RosterMember {
	return []models.RosterMember{
		createRosterMember("tank", "Demon Hunter", "Vengeance", "eu"),
		createRosterMember("healer", "Priest", "Discipline", "eu"),
		createRosterMember("dps", "Mage", "Arcane", "eu"),
		createRosterMember("dps", "Warrior", "Arms", "eu"),
		createRosterMember("dps", "Death Knight", "Unholy", "eu"),
	}
}

func createRosterMember(role, className, specName, region string) models.RosterMember {
	return models.RosterMember{
		Role: role,
		Character: models.Character{
			Class: models.ClassInfo{
				Name: className,
			},
			Spec: models.SpecInfo{
				Name: specName,
			},
			Region: models.RegionInfo{
				Slug: region,
			},
		},
	}
}
