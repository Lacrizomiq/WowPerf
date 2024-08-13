package blizzard

type Service struct {
	Client         *Client
	GameDataClient *GameDataClient
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
	}, nil
}
