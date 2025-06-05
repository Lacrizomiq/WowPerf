package character

import (
	"fmt"
	"sync"
	"time"
)

// UserRateInfo stocke les informations de rate limiting pour un utilisateur
type UserRateInfo struct {
	LastSyncTime    time.Time
	SyncCountToday  int
	LastEnrichTime  time.Time
	EnrichCountHour int
}

// RateLimiter gère les limitations de taux en mémoire
type RateLimiter struct {
	users map[uint]*UserRateInfo
	mutex sync.RWMutex
}

// NewRateLimiter crée un nouveau rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		users: make(map[uint]*UserRateInfo),
		mutex: sync.RWMutex{},
	}
}

// CanSyncUser vérifie si un utilisateur peut synchroniser ses personnages
func (rl *RateLimiter) CanSyncUser(userID uint) (bool, string) {
	rl.mutex.RLock()
	userInfo, exists := rl.users[userID]
	rl.mutex.RUnlock()

	now := time.Now()

	// Si l'utilisateur n'existe pas, il peut synchroniser
	if !exists {
		return true, ""
	}

	// Vérifier le délai minimum entre syncs
	if now.Sub(userInfo.LastSyncTime) < MinDelayBetweenSync {
		remainingTime := MinDelayBetweenSync - now.Sub(userInfo.LastSyncTime)
		return false, fmt.Sprintf("Please wait %v before next sync", remainingTime.Round(time.Second))
	}

	// Vérifier le nombre de syncs par jour
	if isSameDay(userInfo.LastSyncTime, now) && userInfo.SyncCountToday >= MaxSyncPerDay {
		return false, fmt.Sprintf("Maximum %d syncs per day reached", MaxSyncPerDay)
	}

	return true, ""
}

// CanEnrichUser vérifie si un utilisateur peut enrichir ses personnages
func (rl *RateLimiter) CanEnrichUser(userID uint) (bool, string) {
	rl.mutex.RLock()
	userInfo, exists := rl.users[userID]
	rl.mutex.RUnlock()

	now := time.Now()

	// Si l'utilisateur n'existe pas, il peut enrichir
	if !exists {
		return true, ""
	}

	// Vérifier le nombre d'enrichissements par heure
	if isSameHour(userInfo.LastEnrichTime, now) && userInfo.EnrichCountHour >= MaxEnrichPerHour {
		return false, fmt.Sprintf("Maximum %d enrichments per hour reached", MaxEnrichPerHour)
	}

	return true, ""
}

// RecordSync enregistre une synchronisation
func (rl *RateLimiter) RecordSync(userID uint) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	userInfo, exists := rl.users[userID]

	if !exists {
		userInfo = &UserRateInfo{}
		rl.users[userID] = userInfo
	}

	// Reset counter si nouveau jour
	if !isSameDay(userInfo.LastSyncTime, now) {
		userInfo.SyncCountToday = 0
	}

	userInfo.LastSyncTime = now
	userInfo.SyncCountToday++
}

// RecordEnrich enregistre un enrichissement
func (rl *RateLimiter) RecordEnrich(userID uint) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	userInfo, exists := rl.users[userID]

	if !exists {
		userInfo = &UserRateInfo{}
		rl.users[userID] = userInfo
	}

	// Reset counter si nouvelle heure
	if !isSameHour(userInfo.LastEnrichTime, now) {
		userInfo.EnrichCountHour = 0
	}

	userInfo.LastEnrichTime = now
	userInfo.EnrichCountHour++
}

// CleanupOldEntries nettoie les entrées anciennes (à appeler périodiquement)
func (rl *RateLimiter) CleanupOldEntries() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-24 * time.Hour) // Garder les données pendant 24h

	for userID, userInfo := range rl.users {
		// Supprimer les utilisateurs qui n'ont pas synchronisé depuis 24h
		if userInfo.LastSyncTime.Before(cutoff) && userInfo.LastEnrichTime.Before(cutoff) {
			delete(rl.users, userID)
		}
	}
}

// GetUserInfo retourne les informations de rate limiting d'un utilisateur (pour debug)
func (rl *RateLimiter) GetUserInfo(userID uint) *UserRateInfo {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	userInfo, exists := rl.users[userID]
	if !exists {
		return &UserRateInfo{}
	}

	// Retourner une copie pour éviter les modifications concurrentes
	return &UserRateInfo{
		LastSyncTime:    userInfo.LastSyncTime,
		SyncCountToday:  userInfo.SyncCountToday,
		LastEnrichTime:  userInfo.LastEnrichTime,
		EnrichCountHour: userInfo.EnrichCountHour,
	}
}

// Fonctions utilitaires

func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func isSameHour(t1, t2 time.Time) bool {
	return t1.Format("2006-01-02 15") == t2.Format("2006-01-02 15")
}
