package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Config стурктура для хранение конфигурационных данных
type Config struct {
	Symbols    []string `yaml:"symbols"`
	MaxWorkers int      `yaml:"max_workers"`
}

// LoadConfig загружает и возвращает структуру Config
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
