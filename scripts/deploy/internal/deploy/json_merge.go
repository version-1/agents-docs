package deploy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func mergedJSONFile(src, dst string) ([]byte, os.FileMode, error) {
	source, err := os.ReadFile(src)
	if err != nil {
		return nil, 0, fmt.Errorf("read JSON source %q: %w", src, err)
	}
	var incoming map[string]json.RawMessage
	if err := json.Unmarshal(source, &incoming); err != nil || incoming == nil {
		return nil, 0, fmt.Errorf("JSON source %q must be an object: %v", src, err)
	}
	mode := os.FileMode(0600)
	current := make(map[string]json.RawMessage)
	info, err := os.Lstat(dst)
	if err != nil && !os.IsNotExist(err) {
		return nil, 0, fmt.Errorf("stat JSON destination %q: %w", dst, err)
	}
	if err == nil {
		if !info.Mode().IsRegular() {
			return nil, 0, fmt.Errorf("JSON destination %q must be a regular file", dst)
		}
		mode = info.Mode().Perm()
		contents, err := os.ReadFile(dst)
		if err != nil {
			return nil, 0, fmt.Errorf("read JSON destination %q: %w", dst, err)
		}
		if err := json.Unmarshal(contents, &current); err != nil || current == nil {
			return nil, 0, fmt.Errorf("JSON destination %q must be an object: %v", dst, err)
		}
	}
	if err := preservePermissionDenials(current, incoming); err != nil {
		return nil, 0, fmt.Errorf("merge JSON permissions for %q: %w", dst, err)
	}
	for key, value := range incoming {
		current[key] = value
	}
	merged, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return nil, 0, fmt.Errorf("marshal JSON destination %q: %w", dst, err)
	}
	return append(merged, '\n'), mode, nil
}

func preservePermissionDenials(current, incoming map[string]json.RawMessage) error {
	previousRaw, exists := current["permissions"]
	if !exists {
		return nil
	}
	nextRaw, exists := incoming["permissions"]
	if !exists {
		return nil
	}
	var previous, next map[string]json.RawMessage
	if err := json.Unmarshal(previousRaw, &previous); err != nil || previous == nil {
		return fmt.Errorf("existing permissions must be an object: %v", err)
	}
	if err := json.Unmarshal(nextRaw, &next); err != nil || next == nil {
		return fmt.Errorf("source permissions must be an object: %v", err)
	}
	var previousDeny, nextDeny []string
	if raw, exists := previous["deny"]; exists {
		if err := json.Unmarshal(raw, &previousDeny); err != nil {
			return fmt.Errorf("existing permissions.deny must be a string array: %w", err)
		}
	}
	if raw, exists := next["deny"]; exists {
		if err := json.Unmarshal(raw, &nextDeny); err != nil {
			return fmt.Errorf("source permissions.deny must be a string array: %w", err)
		}
	}
	seen := make(map[string]bool, len(nextDeny))
	for _, rule := range nextDeny {
		seen[rule] = true
	}
	for _, rule := range previousDeny {
		if !seen[rule] {
			nextDeny = append(nextDeny, rule)
			seen[rule] = true
		}
	}
	mergedDeny, err := json.Marshal(nextDeny)
	if err != nil {
		return err
	}
	next["deny"] = mergedDeny
	mergedPermissions, err := json.Marshal(next)
	if err != nil {
		return err
	}
	incoming["permissions"] = mergedPermissions
	return nil
}

func writeJSONAtomically(dst string, contents []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("mkdir JSON destination %q: %w", dst, err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(dst), ".deploy-json-*")
	if err != nil {
		return fmt.Errorf("create temporary JSON destination %q: %w", dst, err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(contents); err != nil {
		tmp.Close()
		return fmt.Errorf("write temporary JSON destination %q: %w", dst, err)
	}
	if err := tmp.Chmod(mode); err != nil {
		tmp.Close()
		return fmt.Errorf("chmod temporary JSON destination %q: %w", dst, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temporary JSON destination %q: %w", dst, err)
	}
	if err := os.Rename(tmp.Name(), dst); err != nil {
		return fmt.Errorf("replace JSON destination %q: %w", dst, err)
	}
	return nil
}
