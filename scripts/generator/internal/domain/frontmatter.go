package domain

import (
	"bytes"
	"fmt"
)

// ParseName はマークダウンの YAML フロントマターから name フィールドを抽出する。
func ParseName(content []byte) (string, error) {
	const sep = "---"

	content = bytes.TrimLeft(content, "\n\r")
	if !bytes.HasPrefix(content, []byte(sep)) {
		return "", fmt.Errorf("frontmatter not found")
	}

	rest := content[len(sep):]
	endIdx := bytes.Index(rest, []byte("\n"+sep))
	if endIdx < 0 {
		return "", fmt.Errorf("frontmatter closing delimiter not found")
	}

	block := rest[:endIdx]
	for _, line := range bytes.Split(block, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if bytes.HasPrefix(line, []byte("name:")) {
			val := bytes.TrimSpace(line[len("name:"):])
			if len(val) == 0 {
				return "", fmt.Errorf("name field is empty")
			}
			return string(val), nil
		}
	}

	return "", fmt.Errorf("name field not found in frontmatter")
}
