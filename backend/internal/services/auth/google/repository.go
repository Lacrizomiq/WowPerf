package googleauth

import (
	"fmt"
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

// ===== MÉTHODES MÉTIER SPÉCIFIQUES =====

// GenerateUsernameFromGoogle génère un nom d'utilisateur depuis les infos Google
func (r *GoogleAuthRepository) GenerateUsernameFromGoogle(userInfo *GoogleUserInfo) string {
	// Priorité : prénom > nom complet > partie email
	if userInfo.GivenName != "" {
		return r.ensureUniqueUsername(userInfo.GivenName)
	}

	if userInfo.Name != "" {
		// Prendre le premier mot du nom complet
		parts := strings.Fields(userInfo.Name)
		if len(parts) > 0 {
			return r.ensureUniqueUsername(parts[0])
		}
	}

	// Fallback : partie avant @ de l'email
	if parts := strings.Split(userInfo.Email, "@"); len(parts) > 0 {
		return r.ensureUniqueUsername(parts[0])
	}

	// Dernier fallback
	return r.ensureUniqueUsername("user")
}

// ensureUniqueUsername s'assure que le nom d'utilisateur est unique
func (r *GoogleAuthRepository) ensureUniqueUsername(baseUsername string) string {
	// Nettoyer le nom d'utilisateur (supprimer caractères spéciaux, etc.)
	username := r.cleanUsername(baseUsername)

	// Vérifier si le nom existe déjà
	var count int64
	r.db.Model(&models.User{}).Where("username = ?", username).Count(&count)

	// Si pas de conflit, retourner tel quel
	if count == 0 {
		return username
	}

	// Sinon, ajouter un numéro
	for i := 1; i <= 999; i++ {
		candidateUsername := fmt.Sprintf("%s%d", username, i)
		r.db.Model(&models.User{}).Where("username = ?", candidateUsername).Count(&count)
		if count == 0 {
			return candidateUsername
		}
	}

	// Fallback avec timestamp si vraiment pas de chance
	return fmt.Sprintf("%s%d", username, time.Now().Unix())
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

// CreateUserFromGoogle crée un nouvel utilisateur à partir des infos Google
func (r *GoogleAuthRepository) CreateUserFromGoogle(userInfo *GoogleUserInfo) (*models.User, error) {
	user := &models.User{
		Username: r.GenerateUsernameFromGoogle(userInfo),
		Email:    userInfo.Email,
		Password: "", // Pas de mot de passe pour les comptes Google uniquement
	}

	// Lier le compte Google
	user.LinkGoogleAccount(userInfo.ID, userInfo.Email)

	if err := r.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
