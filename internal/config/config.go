package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Server  ServerConfig  `json:"server"`
	Weather WeatherConfig `json:"weather"`
}

type ServerConfig struct {
	Port                int `json:"port"`
	ReadTimeoutSeconds  int `json:"read_timeout_seconds"`
	WriteTimeoutSeconds int `json:"write_timeout_seconds"`
}

type WeatherConfig struct {
	TemperatureRanges TemperatureRanges `json:"temperature_ranges"`
}

type TemperatureRanges struct {
	HotThreshold  int `json:"hot_threshold"`
	ColdThreshold int `json:"cold_threshold"`
}

func New() (*Config, error) {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.json: %w", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config.json: %w", err)
	}

	return &config, nil
}

func (c *Config) GetReadTimeout() time.Duration {
	return time.Duration(c.Server.ReadTimeoutSeconds) * time.Second
}

func (c *Config) GetWriteTimeout() time.Duration {
	return time.Duration(c.Server.WriteTimeoutSeconds) * time.Second
}
