package config

import (
	"log"
	"os"
	"scavenger/internal/models"

	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	
	log.Printf("Config loadet: %+v", &cfg)
	return &cfg, nil
}

func SaveConfig(cfg *models.Config) error {
	_, err := yaml.Marshal(cfg)
	return err
}
