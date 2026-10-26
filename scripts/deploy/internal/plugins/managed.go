package plugins

import (
	"crypto/sha256"
	"deploy/internal/external"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

type managedPlugin struct {
	Name   string
	Source string
}

func managedPlugins(root string) ([]managedPlugin, error) {
	entries, err := os.ReadDir(root)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	plugins := make([]managedPlugin, 0, len(entries))
	for _, entry := range entries {
		if entry.Type()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("managed plugin %q is a symlink", entry.Name())
		}
		if !entry.IsDir() || !good(entry.Name()) {
			return nil, fmt.Errorf("invalid managed plugin %q", entry.Name())
		}
		plugin := managedPlugin{Name: entry.Name(), Source: filepath.Join(root, entry.Name())}
		if err := validateManagedPlugin(plugin); err != nil {
			return nil, err
		}
		plugins = append(plugins, plugin)
	}
	sort.Slice(plugins, func(i, j int) bool { return plugins[i].Name < plugins[j].Name })
	return plugins, nil
}

func validateManagedPlugin(plugin managedPlugin) error {
	codex, err := readManifest(filepath.Join(plugin.Source, ".codex-plugin", "plugin.json"))
	if err != nil {
		return fmt.Errorf("managed plugin %q: %w", plugin.Name, err)
	}
	if codex["name"] != plugin.Name || codex["skills"] != "./skills/" {
		return fmt.Errorf("managed plugin %q has invalid Codex manifest", plugin.Name)
	}
	claude, err := readManifest(filepath.Join(plugin.Source, ".claude-plugin", "plugin.json"))
	if err != nil {
		return fmt.Errorf("managed plugin %q: %w", plugin.Name, err)
	}
	if claude["name"] != plugin.Name {
		return fmt.Errorf("managed plugin %q has invalid Claude manifest", plugin.Name)
	}
	return validateManagedSkills(filepath.Join(plugin.Source, "skills"))
}

func readManifest(path string) (map[string]any, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("manifest %q is not a regular file", path)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var manifest map[string]any
	if err := json.Unmarshal(b, &manifest); err != nil {
		return nil, fmt.Errorf("invalid manifest %q: %w", path, err)
	}
	if _, ok := manifest["version"].(string); !ok {
		return nil, fmt.Errorf("manifest %q has no string version", path)
	}
	return manifest, nil
}

func validateManagedSkills(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) == 0 {
		return fmt.Errorf("managed plugin skills directory %q is missing or empty", root)
	}
	for _, entry := range entries {
		if entry.Type()&os.ModeSymlink != 0 || !entry.IsDir() || !good(entry.Name()) {
			return fmt.Errorf("invalid managed skill %q", entry.Name())
		}
		info, err := os.Lstat(filepath.Join(root, entry.Name(), "SKILL.md"))
		if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
			return fmt.Errorf("managed skill %q has no regular SKILL.md", entry.Name())
		}
		if err := external.ValidateFetched(external.Skill{Name: entry.Name()}, filepath.Join(root, entry.Name())); err != nil {
			return err
		}
	}
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("refuse symlink %q", path)
		}
		return nil
	})
}

func writeManagedPlugin(root string, plugin managedPlugin) error {
	destination := filepath.Join(root, "plugins", plugin.Name)
	if err := copyTree(plugin.Source, destination); err != nil {
		return err
	}
	version, err := managedPluginVersion(plugin.Source)
	if err != nil {
		return err
	}
	for _, path := range []string{
		filepath.Join(destination, ".codex-plugin", "plugin.json"),
		filepath.Join(destination, ".claude-plugin", "plugin.json"),
	} {
		manifest, err := readManifest(path)
		if err != nil {
			return err
		}
		manifest["version"] = version
		if err := write(path, manifest); err != nil {
			return err
		}
	}
	return nil
}

func managedPluginVersion(root string) (string, error) {
	hash := sha256.New()
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("refuse symlink %q", path)
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if _, err := hash.Write([]byte(rel + info.Mode().Perm().String())); err != nil {
			return err
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if rel == filepath.Join(".codex-plugin", "plugin.json") || rel == filepath.Join(".claude-plugin", "plugin.json") {
			var manifest map[string]any
			if err := json.Unmarshal(body, &manifest); err != nil {
				return err
			}
			delete(manifest, "version")
			body, err = json.Marshal(manifest)
			if err != nil {
				return err
			}
		}
		_, err = hash.Write(body)
		return err
	})
	if err != nil {
		return "", err
	}
	return "0.0.0-" + fmt.Sprintf("%x", hash.Sum(nil))[:12], nil
}
