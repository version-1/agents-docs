package domain

import (
	"bytes"
	"fmt"
)

type Frontmatter struct {
	Name        string
	Description string
}

// ParseFrontmatter はマークダウンの YAML フロントマターを解析し構造体を返す。
func ParseFrontmatter(content []byte) (Frontmatter, error) {
	const sep = "---"

	content = bytes.TrimLeft(content, "\n\r")
	if !bytes.HasPrefix(content, []byte(sep)) {
		return Frontmatter{}, fmt.Errorf("frontmatter not found")
	}

	rest := content[len(sep):]
	endIdx := bytes.Index(rest, []byte("\n"+sep))
	if endIdx < 0 {
		return Frontmatter{}, fmt.Errorf("frontmatter closing delimiter not found")
	}

	var fm Frontmatter
	block := rest[:endIdx]
	for _, line := range bytes.Split(block, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if val, ok := parseField(line, "name"); ok {
			fm.Name = val
		} else if val, ok := parseField(line, "description"); ok {
			fm.Description = val
		}
	}

	if fm.Name == "" {
		return Frontmatter{}, fmt.Errorf("name field is empty or not found in frontmatter")
	}

	return fm, nil
}

func parseField(line []byte, key string) (string, bool) {
	prefix := []byte(key + ":")
	if !bytes.HasPrefix(line, prefix) {
		return "", false
	}
	val := string(bytes.TrimSpace(line[len(prefix):]))
	return val, true
}
