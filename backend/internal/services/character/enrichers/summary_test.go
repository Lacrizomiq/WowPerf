package enrichers

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"wowperf/internal/models"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/common"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// init charge les variables d'environnement depuis .env
func init() {
	// Chercher le fichier .env à la racine du projet
	// Remonter depuis internal/services/character/enrichers vers la racine
	envPaths := []string{
		".env",                // répertoire courant
		"../../../../.env",    // racine du projet depuis enrichers/
		"../../../../../.env", // au cas où la structure change
	}

	loaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Fichier .env chargé depuis: %s", path)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Println("Aucun fichier .env trouvé, utilisation des variables d'environnement système")
	}
}

// Test d'intégration - appelle vraiment l'API Blizzard
func TestFetchCharacterProfileData_Integration(t *testing.T) {
	// Skip ce test si pas en mode intégration
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Vérifier que les credentials sont disponibles
	clientID := os.Getenv("BLIZZARD_CLIENT_ID")
	clientSecret := os.Getenv("BLIZZARD_CLIENT_SECRET")
	region := os.Getenv("BLIZZARD_REGION")

	if clientID == "" || clientSecret == "" || region == "" {
		t.Skip("BLIZZARD_CLIENT_ID, BLIZZARD_CLIENT_SECRET et BLIZZARD_REGION doivent être définis pour les tests d'intégration")
	}

	// Créer un vrai Client puis ProfileService
	client, err := blizzard.NewClient()
	require.NoError(t, err, "Impossible de créer le Client")

	profileService := blizzard.NewProfileService(client)

	// Données de test - utilisez un personnage que vous savez qui existe
	testCases := []struct {
		name          string
		region        string
		realm         string
		characterName string
		namespace     string
		locale        string
	}{
		{
			name:          "Test Ouimagatée",
			region:        "eu",
			realm:         "silvermoon",
			characterName: "ouimagatée",
			namespace:     "profile-eu",
			locale:        "en_GB",
		},
		// Ajoutez d'autres cas de test si nécessaire
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Appel de la vraie fonction
			profile, err := common.FetchCharacterProfileData(
				profileService,
				tc.region,
				tc.realm,
				tc.characterName,
				tc.namespace,
				tc.locale,
			)

			// Vérifications
			require.NoError(t, err, "FetchCharacterProfileData ne doit pas retourner d'erreur")
			require.NotNil(t, profile, "Le profil ne doit pas être nil")

			// Vérifier les données essentielles
			assert.NotEmpty(t, profile.Name, "Le nom ne doit pas être vide")
			assert.NotEmpty(t, profile.Race, "La race ne doit pas être vide")
			assert.NotEmpty(t, profile.Class, "La classe ne doit pas être vide")
			assert.NotEmpty(t, profile.Faction, "La faction ne doit pas être vide")
			assert.NotEmpty(t, profile.ActiveSpecName, "Le spec ne doit pas être vide")
			assert.Greater(t, profile.SpecID, 0, "L'ID du spec doit être > 0")

			// Vérifier les URLs
			assert.NotEmpty(t, profile.AvatarURL, "L'URL de l'avatar ne doit pas être vide")
			assert.Contains(t, profile.AvatarURL, "render.worldofwarcraft.com", "L'URL doit être de Blizzard")

			// Log des données pour debug
			fmt.Printf("✅ Données récupérées pour %s:\n", profile.Name)
			fmt.Printf("   Race: %s\n", profile.Race)
			fmt.Printf("   Classe: %s\n", profile.Class)
			fmt.Printf("   Spé: %s (%s)\n", profile.ActiveSpecName, profile.ActiveSpecRole)
			fmt.Printf("   Faction: %s\n", profile.Faction)
			fmt.Printf("   Points succès: %d\n", profile.AchievementPoints)
			fmt.Printf("   Avatar: %s\n", profile.AvatarURL)
		})
	}
}

// Test d'intégration pour votre SummaryEnricher complet
func TestSummaryEnricher_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	clientID := os.Getenv("BLIZZARD_CLIENT_ID")
	clientSecret := os.Getenv("BLIZZARD_CLIENT_SECRET")
	region := os.Getenv("BLIZZARD_REGION")

	if clientID == "" || clientSecret == "" || region == "" {
		t.Skip("Variables d'environnement Blizzard manquantes")
	}

	// Setup
	client, err := blizzard.NewClient()
	require.NoError(t, err)

	profileService := blizzard.NewProfileService(client)

	enricher := NewSummaryEnricher(profileService)

	// Personnage à tester
	character := &models.UserCharacter{
		Name:   "ouimagatée",
		Realm:  "silvermoon",
		Region: "eu",
	}

	// Test complet
	ctx := context.Background()
	err = enricher.EnrichCharacter(ctx, character)

	// Vérifications
	require.NoError(t, err)

	assert.Equal(t, "Ouimagatée", character.Name) // Nom corrigé par l'API
	assert.Equal(t, "Night Elf", character.Race)
	assert.Equal(t, "Priest", character.Class)
	assert.Equal(t, "Alliance", character.Faction)
	assert.Equal(t, "Discipline", character.ActiveSpecName)
	assert.Equal(t, 256, character.ActiveSpecID)
	assert.Equal(t, "Healer", character.ActiveSpecRole)
	assert.Greater(t, character.AchievementPoints, 0)
	assert.NotEmpty(t, character.AvatarURL)
	assert.NotZero(t, character.LastAPIUpdate)

	fmt.Printf("✅ Enrichissement réussi pour %s - %s %s\n",
		character.Name, character.Race, character.Class)
}
