package filesystem

import (
	"path/filepath"
	"strings"
)

// MatchesExtension checks if a given filename ends with one of the target extensions.
func MatchesExtension(filename string, targetExts []string) bool {
	if len(targetExts) == 0 {
		return false
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return false
	}

	for _, tExt := range targetExts {
		if strings.ToLower(tExt) == ext {
			return true
		}
	}
	return false
}

// MatchesSearch checks if a given filepath contains the search term.
func MatchesSearch(path string, searchTerm string) bool {
	if searchTerm == "" {
		return false
	}
	return strings.Contains(strings.ToLower(path), strings.ToLower(searchTerm))
}
