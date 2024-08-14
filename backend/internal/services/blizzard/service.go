package blizzard

type Service struct {
	Client         *Client
	GameDataClient *GameDataClient
	Profile        *ProfileService
	GameData       *GameDataService
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

func NewService() (*Service, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	gameDataClient, err := NewGameDataClient()
	if err != nil {
		return nil, err
	}

	return &Service{
		Client:         client,
		GameDataClient: gameDataClient,
		Profile:        NewProfileService(client),
		GameData:       NewGameDataService(gameDataClient),
	}, nil
}
