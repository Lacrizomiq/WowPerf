package raiderioMythicPlusRunsTemporalScheduler

import (
	"fmt"
	"time"
)

// MythicPlusRunsScheduleConfig définit la configuration pour l'exécution quotidienne
type MythicPlusRunsScheduleConfig struct {
	Hour      int    // Heure au format 24h UTC
	Minute    int    // Minute
	Day       int    // Jour de la semaine (1 = lundi, 7 = dimanche)
	TaskQueue string // Nom de la queue pour le workflow
}

// DefaultMythicPlusRunsScheduleConfig fournit la configuration par défaut (06h00 UTC)
var DefaultMythicPlusRunsScheduleConfig = MythicPlusRunsScheduleConfig{
	Hour:      6,                    // 6h du matin UTC
	Minute:    0,                    // 0 minute
	Day:       1,                    // Lundi
	TaskQueue: "warcraft-logs-sync", // TODO: a voir si doit etre similaire a celle de warcraft-logs
}

// RetryPolicy définit comment les échecs sont gérés
type RetryPolicy struct {
	InitialInterval    time.Duration // Intervalle initial de retry
	BackoffCoefficient float64       // Multiplicateur pour les retry suivants
	MaximumInterval    time.Duration // Intervalle maximum de retry
	MaximumAttempts    int           // Nombre maximum de tentatives de retry
}

// ScheduleOptions combine toutes les options de configuration
type ScheduleOptions struct {
	Retry   RetryPolicy   // Politique de retry
	Timeout time.Duration // Temps d'exécution maximum
	Paused  bool          // Si le schedule démarre en pause
}

// DefaultScheduleOptions retourne la configuration par défaut
func DefaultScheduleOptions() *ScheduleOptions {
	return &ScheduleOptions{
		Retry: RetryPolicy{
			InitialInterval:    30 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Minute,
			MaximumAttempts:    5,
		},
		Timeout: 8 * time.Hour, // Temps suffisant pour traiter toutes les régions/donjons
		Paused:  false,
	}
}

// ValidateScheduleConfig valide la configuration du schedule
func ValidateScheduleConfig(config *MythicPlusRunsScheduleConfig) error {
	if config == nil {
		return fmt.Errorf("schedule config cannot be nil")
	}
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
