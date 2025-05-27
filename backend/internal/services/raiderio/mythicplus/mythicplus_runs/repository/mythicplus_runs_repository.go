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
	SkippedRuns          int // Runs déjà existants
	NewCompositions      int
	ExistingCompositions int
}

// ProcessRuns traite les runs Mythic+ et les enregistre dans la base de données
// Logique de traitement: si keystone_run_id existe déjà, on l'ignore
// Sinon, on l'enregistre
func (r *MythicPlusRunsRepository) ProcessRuns(runs []*models.Run, batchID string) (*ProcessingStats, error) {
	stats := &ProcessingStats{}

	return stats, r.db.Transaction(func(tx *gorm.DB) error {
		for _, run := range runs {
			// 1. Valide la run (logique métier basique)
			if !r.isValidRun(run) {
				stats.SkippedRuns++
				continue
			}

			// 2. Gère la composition d'équipe
			teamComp, isNewComp, err := r.getOrCreateTeamComposition(tx, run.Roster)
			if err != nil {
				return fmt.Errorf("failed to handle team composition: %w", err)
			}

			if isNewComp {
				stats.NewCompositions++
			} else {
				stats.ExistingCompositions++
			}

			// 3. Tente d'insérer la run (skip si existe déjà grâce à la contrainte UNIQUE)
			if err := r.insertRun(tx, run, teamComp.ID, batchID); err != nil {
				// Si c'est une violation de contrainte unique, on skip
				if r.isDuplicateError(err) {
					stats.SkippedRuns++
					continue
				}
				return fmt.Errorf("failed to insert run: %w", err)
			}

			stats.NewRuns++
		}

		return nil
	})
}

// insertRun insère une nouvelle run
func (r *MythicPlusRunsRepository) insertRun(tx *gorm.DB, run *models.Run, teamCompID uint, batchID string) error {
	dbRun := &models.MythicPlusRuns{
		KeystoneRunID:     run.KeystoneRunID,
		Season:            run.Season,
		Region:            r.extractRegion(run),
		DungeonSlug:       run.Dungeon.Slug,
		DungeonName:       run.Dungeon.Name,
		MythicLevel:       run.MythicLevel,
		Status:            run.Status,
		ClearTimeMs:       run.ClearTimeMs,
		KeystoneTimeMs:    run.KeystoneTimeMs,
		CompletedAt:       run.CompletedAt,
		NumChests:         run.NumChests,
		TimeRemainingMs:   run.TimeRemainingMs,
		TeamCompositionID: &teamCompID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := tx.Create(dbRun).Error; err != nil {
		return err
	}

	// Crée les entrées roster
	return r.createRosterEntries(tx, teamCompID, run.Roster)
}

// getOrCreateTeamComposition gère la composition d'équipe avec déduplication
func (r *MythicPlusRunsRepository) getOrCreateTeamComposition(tx *gorm.DB, roster []models.RosterMember) (*models.MythicPlusTeamComposition, bool, error) {
	// 1. Extrait et trie les rôles
	tank, healer, dps := r.extractRoles(roster)

	// 2. Crée le hash de composition
	compHash := r.createCompositionHash(tank, healer, dps)

	// 3. Cherche si elle existe déjà
	var existing models.MythicPlusTeamComposition
	err := tx.Where("composition_hash = ?", compHash).First(&existing).Error

	if err == nil {
		return &existing, false, nil // Existe déjà
	}

	if err != gorm.ErrRecordNotFound {
		return nil, false, fmt.Errorf("failed to query team composition: %w", err)
	}

	// 4. Crée nouvelle composition
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

// createRosterEntries crée les entrées détaillées du roster
func (r *MythicPlusRunsRepository) createRosterEntries(tx *gorm.DB, teamCompID uint, roster []models.RosterMember) error {
	for _, member := range roster {
		rosterEntry := &models.MythicPlusRunRoster{
			TeamCompositionID: teamCompID,
			Role:              member.Role,
			ClassName:         member.Character.Class.Name,
			SpecName:          member.Character.Spec.Name,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		if err := tx.Create(rosterEntry).Error; err != nil {
			return fmt.Errorf("failed to create roster entry: %w", err)
		}
	}

	return nil
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

// isDuplicateError vérifie si l'erreur est due à une contrainte unique
func (r *MythicPlusRunsRepository) isDuplicateError(err error) bool {
	// À adapter selon ta DB (PostgreSQL, MySQL, etc.)
	errorStr := err.Error()
	return strings.Contains(errorStr, "duplicate key") ||
		strings.Contains(errorStr, "UNIQUE constraint") ||
		strings.Contains(errorStr, "Duplicate entry")
}
