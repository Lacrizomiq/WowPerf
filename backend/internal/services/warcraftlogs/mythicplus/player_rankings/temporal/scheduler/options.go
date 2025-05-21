package warcraftlogsPlayerRankingsTemporalScheduler

import (
	"fmt"
	"time"
)

// PlayerRankingsScheduleConfig définit la configuration pour l'exécution quotidienne
type PlayerRankingsScheduleConfig struct {
	Hour      int    // Heure au format 24h UTC
	Minute    int    // Minute
	TaskQueue string // Nom de la queue pour le workflow
}

// DefaultPlayerRankingsScheduleConfig fournit la configuration par défaut (12h00 UTC)
var DefaultPlayerRankingsScheduleConfig = PlayerRankingsScheduleConfig{
	Hour:      12, // 12h
	Minute:    0,  // 0 minute
	TaskQueue: "warcraft-logs-sync",
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
			InitialInterval:    time.Minute,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Hour,
			MaximumAttempts:    5,
		},
		Timeout: 6 * time.Hour, // Temps suffisant pour l'exécution
		Paused:  false,
	}
}

// ValidateScheduleConfig valide la configuration du schedule
func ValidateScheduleConfig(config *PlayerRankingsScheduleConfig) error {
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
