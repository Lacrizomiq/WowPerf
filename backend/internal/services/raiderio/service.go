package raiderio

import "time"

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

// GetRequestStats expose les statistiques de rate limiting
func (s *RaiderIOService) GetRequestStats() (total int64, duration time.Duration, avgPerHour float64) {
	return s.Client.GetRequestStats()
}

// LogRequestSummary log les statistiques de rate limiting
func (s *RaiderIOService) LogRequestSummary() {
	s.Client.LogRequestSummary()
}
