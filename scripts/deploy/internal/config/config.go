package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Items []Item `json:"items"`
}

type Item struct {
	Source      string   `json:"source"`
	Destination string   `json:"destination"`
	Exclude     []string `json:"exclude"`
	Replace     bool     `json:"replace"`
	Flatten     bool     `json:"flatten"`
	Template    bool     `json:"template"`
	MergeJSON   bool     `json:"mergeJSON"`
}

func Load(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %q: %w", path, err)
	}
	if len(cfg.Items) == 0 {
		return Config{}, fmt.Errorf("config must include at least one item")
	}

	for i, item := range cfg.Items {
		if item.Source == "" {
			return Config{}, fmt.Errorf("items[%d].source is required", i)
		}
		if item.Destination == "" {
			return Config{}, fmt.Errorf("items[%d].destination is required", i)
		}
		if item.Template && item.Flatten {
			return Config{}, fmt.Errorf("items[%d]: template and flatten cannot be used together", i)
		}
		if item.MergeJSON && (item.Template || item.Flatten || item.Replace) {
			return Config{}, fmt.Errorf("items[%d]: mergeJSON cannot be combined with template, flatten, or replace", i)
		}
	}
	return cfg, nil
}
