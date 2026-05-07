package pathutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ResolveConfigPath(path string) (string, error) {
	expanded, err := ExpandHome(path)
	if err != nil {
		return "", err
	}
	return filepath.Abs(expanded)
}

func ResolveItemPath(configDir, path string) (string, error) {
	expanded, err := ExpandHome(path)
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(expanded) {
		return filepath.Clean(expanded), nil
	}
	return filepath.Clean(filepath.Join(configDir, expanded)), nil
}

func ResolveSourcePath(cwd, path string) (string, error) {
	expanded, err := ExpandHome(path)
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(expanded) {
		return filepath.Clean(expanded), nil
	}
	return filepath.Clean(filepath.Join(cwd, expanded)), nil
}

func ExpandHome(path string) (string, error) {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory: %w", err)
		}
		return home, nil
	}

	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory: %w", err)
		}
		return filepath.Join(home, strings.TrimPrefix(path, "~/")), nil
	}
	return path, nil
}
