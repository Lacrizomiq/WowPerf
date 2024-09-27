package raiderio

func NewRaiderIOService() (*RaiderIOService, error) {
	client, err := NewRaiderIOClient()
	if err != nil {
		return nil, err
	}

	return &RaiderIOService{
		Client: client,
	}, nil
}
