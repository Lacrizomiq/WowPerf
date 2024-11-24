package blizzard

import (
	"os"
	"wowperf/internal/services/blizzard/auth"
	protectedProfile "wowperf/internal/services/blizzard/protected/profile"
	"wowperf/internal/services/blizzard/types"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Service struct {
	Client           *Client
	ProtectedClient  *ProtectedClient
	GameDataClient   *GameDataClient
	Profile          *ProfileService
	ProtectedProfile *protectedProfile.ProtectedProfileService
	GameData         *GameDataService
	BattleNetAuth    *auth.BattleNetAuthService
}

var _ types.ProtectedClientInterface = (*ProtectedClient)(nil)

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

	protectedClient := NewProtectedClient(os.Getenv("BLIZZARD_REGION"), battleNetAuth)
	protectedProfileService := protectedProfile.NewProtectedProfileService(protectedClient)
	return &Service{
		Client:           client,
		ProtectedClient:  protectedClient,
		GameDataClient:   gameDataClient,
		Profile:          NewProfileService(client),
		GameData:         NewGameDataService(gameDataClient),
		BattleNetAuth:    battleNetAuth,
		ProtectedProfile: protectedProfileService,
	}, nil
}
