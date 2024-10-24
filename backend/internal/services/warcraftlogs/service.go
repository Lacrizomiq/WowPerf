// package warcraftlogs/service.go
package warcraftlogs

type Service struct {
	Client *Client
}

func NewService() (*Service, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	return &Service{
		Client: client,
	}, nil
}
