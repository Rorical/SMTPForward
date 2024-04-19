package util

import (
	"encoding/json"
	"github.com/Rorical/SMTPForward/client"
	"github.com/Rorical/SMTPForward/server"
	"os"
)

func ReadClientConfig(path string) (*client.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := &client.Config{}
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func ReadServerConfig(path string) (*server.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := &server.Config{}
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
