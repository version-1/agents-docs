package deploy

import (
	"fmt"
	"slices"

	"deploy/internal/config"
)

func selectTaggedItems(cfg config.Config, tag string) (config.Config, error) {
	if tag == "" {
		return cfg, nil
	}
	var selected config.Config
	for _, item := range cfg.Items {
		if slices.Contains(item.Tags, tag) {
			selected.Items = append(selected.Items, item)
		}
	}
	if len(selected.Items) == 0 {
		return config.Config{}, fmt.Errorf("no deployment items match tag %q", tag)
	}
	return selected, nil
}
