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

// GetRaw retourne les données JSON brutes (optimisé)
func (s *RaiderIOService) GetRaw(endpoint string, params map[string]string) ([]byte, error) {
	return s.Client.GetRaw(endpoint, params)
}

// GetRateLimitStats expose les statistiques de rate limiting
func (s *RaiderIOService) GetRateLimitStats() (current int, total int64, remaining int, maxRequests int) {
	return s.Client.GetRateLimitStats()
}
