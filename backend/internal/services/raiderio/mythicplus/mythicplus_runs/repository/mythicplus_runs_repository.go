package raiderioMythicPlusRunsRepository

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
	"time"

	models "wowperf/internal/models/raiderio/mythicplus_runs"

	"gorm.io/gorm"
)

// MythicPlusRunsRepository représente le repository pour les runs Mythic+
type MythicPlusRunsRepository struct {
	db *gorm.DB
}

// NewMythicPlusRunsRepository crée un nouveau repository Mythic+
func NewMythicPlusRunsRepository(db *gorm.DB) *MythicPlusRunsRepository {
	return &MythicPlusRunsRepository{db: db}
}

// ProcessingStats contient les stats de traitement pour Temporal (optionnel)
type ProcessingStats struct {
	NewRuns              int
	UpdatedRuns          int // Runs mises à jour
	SkippedRuns          int // Runs déjà existants
	NewCompositions      int
	ExistingCompositions int
}

// ProcessRuns récupère les runs depuis l'API et les insère dans la DB
func (r *MythicPlusRunsRepository) ProcessRuns(runs []*models.Run, batchID string) (*ProcessingStats, error) {
	stats := &ProcessingStats{}

	return stats, r.db.Transaction(func(tx *gorm.DB) error {
		const batchSize = 100

		// Collections pour batch processing
		var validRuns []*models.MythicPlusRuns
		teamCompositionCache := make(map[string]*models.MythicPlusTeamComposition)

		// 1. Prépare toutes les données en mémoire d'abord
		for _, run := range runs {
			// Valide la run
			if !r.isValidRun(run) {
				stats.SkippedRuns++
				continue
			}

			// Gère la composition d'équipe avec cache
			teamComp, isNewComp, err := r.getOrCreateTeamCompositionCached(tx, run.Roster, teamCompositionCache)
			if err != nil {
				return fmt.Errorf("failed to handle team composition: %w", err)
			}

			if isNewComp {
				stats.NewCompositions++
			} else {
				stats.ExistingCompositions++
			}

			// Prépare la run pour batch insert - AVEC LE SCORE
			dbRun := &models.MythicPlusRuns{
				KeystoneRunID:     run.KeystoneRunID,
				Season:            run.Season,
				Region:            r.extractRegion(run),
				DungeonSlug:       run.Dungeon.Slug,
				DungeonName:       run.Dungeon.Name,
				MythicLevel:       run.MythicLevel,
				Score:             run.Score,
				Status:            run.Status,
				ClearTimeMs:       run.ClearTimeMs,
				KeystoneTimeMs:    run.KeystoneTimeMs,
				CompletedAt:       run.CompletedAt,
				NumChests:         run.NumChests,
				TimeRemainingMs:   run.TimeRemainingMs,
				TeamCompositionID: &teamComp.ID,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}

			validRuns = append(validRuns, dbRun)
		}

		// 2. Batch insert/update des runs avec stats détaillées
		if len(validRuns) > 0 {
			newRuns, updatedRuns, err := r.batchUpsertRuns(tx, validRuns, batchSize)
			if err != nil {
				return fmt.Errorf("failed to batch upsert runs: %w", err)
			}
			stats.NewRuns = newRuns
			stats.UpdatedRuns = updatedRuns
		}

		return nil
	})
}

// batchUpsertRuns gère l'insertion/update par batch avec UPSERT
// Retourne (nouvelles_runs, runs_mises_à_jour, erreur)
func (r *MythicPlusRunsRepository) batchUpsertRuns(tx *gorm.DB, runs []*models.MythicPlusRuns, batchSize int) (int, int, error) {
	newRuns := 0
	updatedRuns := 0

	for i := 0; i < len(runs); i += batchSize {
		end := i + batchSize
		if end > len(runs) {
			end = len(runs)
		}

		batch := runs[i:end]

		// Utilise UPSERT pour mettre à jour les runs existantes
		for _, run := range batch {
			var existingRun models.MythicPlusRuns
			result := tx.Where("keystone_run_id = ?", run.KeystoneRunID).First(&existingRun)

			if result.Error == gorm.ErrRecordNotFound {
				// Run n'existe pas, on l'insère
				if err := tx.Create(run).Error; err != nil {
					return newRuns, updatedRuns, fmt.Errorf("failed to create run %d: %w", run.KeystoneRunID, err)
				}
				newRuns++
			} else if result.Error == nil {
				// Run existe, on la met à jour avec les nouvelles données
				if err := tx.Model(&existingRun).
					Updates(map[string]interface{}{
						"score":               run.Score,
						"season":              run.Season,
						"region":              run.Region,
						"dungeon_slug":        run.DungeonSlug,
						"dungeon_name":        run.DungeonName,
						"mythic_level":        run.MythicLevel,
						"status":              run.Status,
						"clear_time_ms":       run.ClearTimeMs,
						"keystone_time_ms":    run.KeystoneTimeMs,
						"completed_at":        run.CompletedAt,
						"num_chests":          run.NumChests,
						"time_remaining_ms":   run.TimeRemainingMs,
						"team_composition_id": run.TeamCompositionID,
						"updated_at":          time.Now(),
					}).Error; err != nil {
					return newRuns, updatedRuns, fmt.Errorf("failed to update run %d: %w", run.KeystoneRunID, err)
				}
				updatedRuns++
			} else {
				// Erreur lors de la recherche
				return newRuns, updatedRuns, fmt.Errorf("failed to check if run %d exists: %w", run.KeystoneRunID, result.Error)
			}
		}
	}
	return newRuns, updatedRuns, nil
}

// getOrCreateTeamCompositionCached fait une mise en cache pour éviter les requêtes répétées
func (r *MythicPlusRunsRepository) getOrCreateTeamCompositionCached(
	tx *gorm.DB,
	roster []models.RosterMember,
	cache map[string]*models.MythicPlusTeamComposition,
) (*models.MythicPlusTeamComposition, bool, error) {

	// 1. Extrait et trie les rôles
	tank, healer, dps := r.extractRoles(roster)

	// 2. Crée le hash de composition
	compHash := r.createCompositionHash(tank, healer, dps)

	// 3. Vérifie le cache d'abord
	if cachedComp, exists := cache[compHash]; exists {
		return cachedComp, false, nil
	}

	// 4. Cherche en DB
	var existing models.MythicPlusTeamComposition
	err := tx.Where("composition_hash = ?", compHash).First(&existing).Error

	if err == nil {
		// Ajoute au cache et retourne
		cache[compHash] = &existing
		return &existing, false, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, false, fmt.Errorf("failed to query team composition: %w", err)
	}

	// 5. Crée nouvelle composition
	newComp := &models.MythicPlusTeamComposition{
		CompositionHash: compHash,
		TankClass:       tank.Character.Class.Name,
		TankSpec:        tank.Character.Spec.Name,
		HealerClass:     healer.Character.Class.Name,
		HealerSpec:      healer.Character.Spec.Name,
		Dps1Class:       dps[0].Character.Class.Name,
		Dps1Spec:        dps[0].Character.Spec.Name,
		Dps2Class:       dps[1].Character.Class.Name,
		Dps2Spec:        dps[1].Character.Spec.Name,
		Dps3Class:       dps[2].Character.Class.Name,
		Dps3Spec:        dps[2].Character.Spec.Name,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := tx.Create(newComp).Error; err != nil {
		return nil, false, fmt.Errorf("failed to create team composition: %w", err)
	}

	// Ajoute au cache
	cache[compHash] = newComp
	return newComp, true, nil
}

// extractRoles extrait et organise les rôles du roster
func (r *MythicPlusRunsRepository) extractRoles(roster []models.RosterMember) (tank, healer models.RosterMember, dps []models.RosterMember) {
	for _, member := range roster {
		switch member.Role {
		case "tank":
			tank = member
		case "healer":
			healer = member
		case "dps":
			dps = append(dps, member)
		}
	}

	// Trie les DPS pour cohérence du hash
	sort.Slice(dps, func(i, j int) bool {
		if dps[i].Character.Class.Name == dps[j].Character.Class.Name {
			return dps[i].Character.Spec.Name < dps[j].Character.Spec.Name
		}
		return dps[i].Character.Class.Name < dps[j].Character.Class.Name
	})

	return tank, healer, dps
}

// createCompositionHash crée un hash unique pour la composition
func (r *MythicPlusRunsRepository) createCompositionHash(tank, healer models.RosterMember, dps []models.RosterMember) string {
	parts := []string{
		tank.Character.Class.Name + "_" + tank.Character.Spec.Name,
		healer.Character.Class.Name + "_" + healer.Character.Spec.Name,
	}

	for _, d := range dps {
		parts = append(parts, d.Character.Class.Name+"_"+d.Character.Spec.Name)
	}

	composition := strings.Join(parts, "|")
	return fmt.Sprintf("%x", md5.Sum([]byte(composition)))
}

// isValidRun valide qu'une run est correcte (logique métier basique)
func (r *MythicPlusRunsRepository) isValidRun(run *models.Run) bool {
	// Vérifie qu'on a bien 5 membres
	if len(run.Roster) != 5 {
		return false
	}

	// Vérifie la composition (1 tank, 1 heal, 3 dps)
	roleCount := make(map[string]int)
	for _, member := range run.Roster {
		roleCount[member.Role]++
	}

	return roleCount["tank"] == 1 && roleCount["healer"] == 1 && roleCount["dps"] == 3
}

// extractRegion extrait la région depuis les données de la run
func (r *MythicPlusRunsRepository) extractRegion(run *models.Run) string {
	if len(run.Roster) > 0 {
		return run.Roster[0].Character.Region.Slug
	}
	return "unknown"
}
