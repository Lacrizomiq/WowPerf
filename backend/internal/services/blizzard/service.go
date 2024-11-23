package blizzard

import (
	"os"
	"wowperf/internal/services/blizzard/auth"

	"github.com/go-redis/redis/v8"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type Service struct {
	Client          *Client
	ProtectedClient *ProtectedClient
	GameDataClient  *GameDataClient
	Profile         *ProfileService
	GameData        *GameDataService
	BattleNetAuth   *auth.BattleNetAuthService
}

type ProfileService struct {
	Client *Client
}

type GameDataService struct {
	Client *GameDataClient
}

func NewProfileService(client *Client) *ProfileService {
	return &ProfileService{
		Client: client,
	}
}

func NewGameDataService(client *GameDataClient) *GameDataService {
	return &GameDataService{
		Client: client,
	}
}

func NewService(db *gorm.DB, redisClient *redis.Client) (*Service, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	gameDataClient, err := NewGameDataClient()
	if err != nil {
		return nil, err
	}

	battleNetAuth, err := auth.NewBattleNetAuthService(db, redisClient)
	if err != nil {
		return nil, err
	}

	oauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("BLIZZARD_CLIENT_ID"),
		ClientSecret: os.Getenv("BLIZZARD_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("BLIZZARD_REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.battle.net/authorize",
			TokenURL: "https://oauth.battle.net/token",
		},
		Scopes: []string{"wow.profile"},
	}

	protectedClient := NewProtectedClient(os.Getenv("BLIZZARD_REGION"), oauthConfig)

	return &Service{
		Client:          client,
		ProtectedClient: protectedClient,
		GameDataClient:  gameDataClient,
		Profile:         NewProfileService(client),
		GameData:        NewGameDataService(gameDataClient),
		BattleNetAuth:   battleNetAuth,
	}, nil
}
