package aws

import "github.com/sarmerer/go-crypto-dashboard/hosting"

type AWSConfig struct {
	AccessKey string
	SecretKey string
	Region    string
}

type awsServer struct {
	config *AWSConfig
}

func NewServer(config *AWSConfig) (hosting.Hosting, error) {
	return &awsServer{config}, nil
}

func (s *awsServer) Start() error {
	return nil
}

func (s *awsServer) Stop() error {
	return nil
}
