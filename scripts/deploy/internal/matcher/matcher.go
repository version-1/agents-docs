package matcher

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Matcher struct {
	patterns []excludePattern
}

type excludePattern struct {
	raw        string
	re         *regexp.Regexp
	basenameRe *regexp.Regexp
}

func New(patterns []string) (Matcher, error) {
	matcher := Matcher{}
	for _, pattern := range patterns {
		normalized := normalizeGlobPattern(pattern)
		if normalized == "" {
			continue
		}

		re, err := compileGlob(normalized)
		if err != nil {
			return Matcher{}, fmt.Errorf("compile exclude pattern %q: %w", pattern, err)
		}

		exclude := excludePattern{raw: normalized, re: re}
		if !strings.Contains(normalized, "/") {
			exclude.basenameRe = re
		}
		matcher.patterns = append(matcher.patterns, exclude)
	}
	return matcher, nil
}

func (m Matcher) Match(rel string) bool {
	rel = normalizeRelPath(rel)
	if rel == "." || rel == "" {
		return false
	}

	for _, pattern := range m.patterns {
		if pattern.re.MatchString(rel) || pattern.re.MatchString(rel+"/") {
			return true
		}
		if pattern.basenameRe != nil && pattern.basenameRe.MatchString(path.Base(rel)) {
			return true
		}
	}
	return false
}

func normalizeRelPath(rel string) string {
	return filepath.ToSlash(filepath.Clean(rel))
}

func normalizeGlobPattern(pattern string) string {
	pattern = strings.TrimSpace(filepath.ToSlash(pattern))
	pattern = strings.TrimPrefix(pattern, "./")
	pattern = strings.TrimSuffix(pattern, "/")
	return pattern
}

// CacheKeyPatterns normalizes exclude patterns for order-independent cache keys.
// Matcher semantics are currently order-independent; update this if that changes.
func CacheKeyPatterns(patterns []string) []string {
	normalized := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		pattern = normalizeGlobPattern(pattern)
		if pattern != "" {
			normalized = append(normalized, pattern)
		}
	}
	sort.Strings(normalized)
	return normalized
}

func compileGlob(pattern string) (*regexp.Regexp, error) {
	var b strings.Builder
	b.WriteString("^")

	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				i++
				if i+1 < len(pattern) && pattern[i+1] == '/' {
					b.WriteString("(?:.*/)?")
					i++
				} else {
					b.WriteString(".*")
				}
			} else {
				b.WriteString("[^/]*")
			}
		case '?':
			b.WriteString("[^/]")
		default:
			b.WriteString(regexp.QuoteMeta(string(pattern[i])))
		}
	}

	b.WriteString("$")
	return regexp.Compile(b.String())
}
