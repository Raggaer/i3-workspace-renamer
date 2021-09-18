package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Separator            string            `json:"separator"`
	ClassShortNames      map[string]string `json:"class"`
	WindowNameShortNames map[string]string `json:"window"`
}

func loadConfigurationFile(path string) (*config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
