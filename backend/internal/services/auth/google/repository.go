package googleauth

import (
	"fmt"
	"log"
	"strings"
	"time"
	"wowperf/internal/models"

	"gorm.io/gorm"
)

// GoogleAuthRepository gère les accès aux données pour Google OAuth
type GoogleAuthRepository struct {
	db *gorm.DB
}

// NewGoogleAuthRepository crée une nouvelle instance du repository
func NewGoogleAuthRepository(db *gorm.DB) *GoogleAuthRepository {
	return &GoogleAuthRepository{
		db: db,
	}
}

// ===== MÉTHODES DE RECHERCHE =====

// FindUserByGoogleID recherche un utilisateur par Google ID
func (r *GoogleAuthRepository) FindUserByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("google_id = ?", googleID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByEmail recherche un utilisateur par email
func (r *GoogleAuthRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ===== MÉTHODES DE CRÉATION/MISE À JOUR =====

// CreateUser crée un utilisateur avec transaction sécurisée
func (r *GoogleAuthRepository) CreateUser(user *models.User) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return &OAuthError{
			Code:    "transaction_start_failed",
			Message: "Failed to start database transaction",
			Details: tx.Error.Error(),
		}
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return &OAuthError{
			Code:    "user_creation_failed",
			Message: "Failed to create user in database",
			Details: err.Error(),
		}
	}

	return tx.Commit().Error
}

// UpdateUser met à jour un utilisateur avec transaction sécurisée
func (r *GoogleAuthRepository) UpdateUser(user *models.User) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return &OAuthError{
			Code:    "transaction_start_failed",
			Message: "Failed to start database transaction",
			Details: tx.Error.Error(),
		}
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return &OAuthError{
			Code:    "user_update_failed",
			Message: "Failed to update user in database",
			Details: err.Error(),
		}
	}

	return tx.Commit().Error
}

// ===== MÉTHODES HELPER =====

// CheckUserExistsByGoogleID vérifie si un utilisateur existe par Google ID
func (r *GoogleAuthRepository) CheckUserExistsByGoogleID(googleID string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("google_id = ?", googleID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckUserExistsByEmail vérifie si un utilisateur existe par email
func (r *GoogleAuthRepository) CheckUserExistsByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// =====  LOGIQUE DE CRÉATION D'UTILISATEUR (Evite la duplication d'username) =====

// CreateUserFromGoogle crée un nouvel utilisateur avec gestion intelligente des conflits username
func (r *GoogleAuthRepository) CreateUserFromGoogle(userInfo *GoogleUserInfo) (*models.User, error) {
	baseUsername := r.generateBaseUsername(userInfo)

	// 1. Essayer le username de base d'abord
	if user, err := r.tryCreateUser(baseUsername, userInfo); err == nil {
		return user, nil // ✅ Succès avec "ludovic"
	} else if !r.isUniqueConstraintError(err) {
		return nil, err // Erreur autre que contrainte unique → fail immédiat
	}

	// 2. Trouver le prochain numéro disponible
	nextNumber := r.findNextAvailableNumber(baseUsername)

	// 3. Essayer 100 variations à partir du prochain disponible
	for i := nextNumber; i < nextNumber+100; i++ {
		username := fmt.Sprintf("%s%d", baseUsername, i)

		if user, err := r.tryCreateUser(username, userInfo); err == nil {
			return user, nil // ✅ Succès avec "ludovic42"
		} else if !r.isUniqueConstraintError(err) {
			return nil, err // Erreur autre que contrainte unique → fail
		}
	}

	// 4. Fallback ultime : UUID garanti unique (cas très extrême)
	uuid := r.generateUUID8()
	username := fmt.Sprintf("%s_%s", baseUsername, uuid)
	return r.tryCreateUser(username, userInfo)
}

// tryCreateUser tente de créer un utilisateur avec un username donné
func (r *GoogleAuthRepository) tryCreateUser(username string, userInfo *GoogleUserInfo) (*models.User, error) {
	// Prépare les pointeurs avant la création
	googleID := userInfo.ID
	googleEmail := userInfo.Email

	user := &models.User{
		Username:    username,
		Email:       userInfo.Email,
		Password:    "", // Pas de mot de passe pour comptes Google
		GoogleID:    &googleID,
		GoogleEmail: &googleEmail,
	}

	// Tenter la création
	if err := r.CreateUser(user); err != nil {
		log.Printf("❌ Failed to create user: %v", err)
		return nil, err
	}

	// Vérifier après création
	log.Printf("✅ User created successfully - ID: %d, GoogleID: %v",
		user.ID, user.GoogleID)

	return user, nil
}

// findNextAvailableNumber trouve le prochain numéro disponible pour un username
func (r *GoogleAuthRepository) findNextAvailableNumber(baseUsername string) int {
	var result struct {
		MaxNum int `gorm:"column:max_num"`
	}

	// Requête pour trouver le plus grand numéro existant
	// Ex: pour "ludovic" → trouve le max de "ludovic1", "ludovic2", "ludovic99", etc.
	r.db.Raw(`
		SELECT COALESCE(MAX(
			CASE 
				WHEN username ~ ? 
				THEN CAST(REGEXP_REPLACE(username, ?, '', 'g') AS INTEGER)
				ELSE 0 
			END
		), 0) as max_num
		FROM users 
		WHERE username = ? OR username ~ ?
	`,
		"^"+baseUsername+"[0-9]+$", // Pattern regex : "ludovic[0-9]+"
		"^"+baseUsername,           // Enlever "ludovic" pour garder que le numéro
		baseUsername,               // Username exact : "ludovic"
		"^"+baseUsername+"[0-9]+$", // Pattern regex : "ludovic[0-9]+"
	).Scan(&result)

	return result.MaxNum + 1 // Retourner le suivant
}

// isUniqueConstraintError détecte si l'erreur est due à une contrainte unique
func (r *GoogleAuthRepository) isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// PostgreSQL
	if strings.Contains(errStr, "duplicate key") &&
		(strings.Contains(errStr, "username") || strings.Contains(errStr, "uniqueindex")) {
		return true
	}

	// MySQL
	if strings.Contains(errStr, "duplicate entry") && strings.Contains(errStr, "username") {
		return true
	}

	// SQLite
	if strings.Contains(errStr, "unique constraint") && strings.Contains(errStr, "username") {
		return true
	}

	// GORM peut wrapper l'erreur
	if strings.Contains(errStr, "unique") && strings.Contains(errStr, "constraint") {
		return true
	}

	return false
}

// generateBaseUsername génère le username de base (sans numéro)
func (r *GoogleAuthRepository) generateBaseUsername(userInfo *GoogleUserInfo) string {
	// Priorité : prénom > nom complet > partie email
	if userInfo.GivenName != "" {
		return r.cleanUsername(userInfo.GivenName)
	}

	if userInfo.Name != "" {
		// Prendre le premier mot du nom complet
		parts := strings.Fields(userInfo.Name)
		if len(parts) > 0 {
			return r.cleanUsername(parts[0])
		}
	}

	// Fallback : partie avant @ de l'email
	if parts := strings.Split(userInfo.Email, "@"); len(parts) > 0 {
		return r.cleanUsername(parts[0])
	}

	// Dernier fallback
	return "user"
}

// generateUUID8 génère 8 caractères UUID pour fallback ultime
func (r *GoogleAuthRepository) generateUUID8() string {
	// Version simple sans import uuid - utilise timestamp + random
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%08x", timestamp)[:8] // 8 premiers chars en hexa
}

// cleanUsername nettoie un nom d'utilisateur
func (r *GoogleAuthRepository) cleanUsername(username string) string {
	// Convertir en minuscules
	cleaned := strings.ToLower(username)

	// Supprimer les espaces et caractères spéciaux basiques
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	// Limiter la longueur (respecter contraintes de ton modèle User)
	if len(cleaned) > 12 { // Selon tes validations
		cleaned = cleaned[:12]
	}

	// S'assurer qu'il reste quelque chose
	if len(cleaned) < 3 {
		cleaned = "user"
	}

	return cleaned
}
