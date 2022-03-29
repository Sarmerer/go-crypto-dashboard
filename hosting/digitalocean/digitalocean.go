package digitalocean

import "github.com/sarmerer/go-crypto-dashboard/hosting"

type DigitalOceanConfig struct {
	Token string
}

type localServer struct {
	config *DigitalOceanConfig
}

func NewServer(config *DigitalOceanConfig) (hosting.Provider, error) {
	return &localServer{config}, nil
}

func (s *localServer) Start() error {
	return nil
}

func (s *localServer) Stop() error {
	return nil
}
