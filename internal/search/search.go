package search

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
)

// Index wraps a bleve index for vault search.
type Index struct {
	index     bleve.Index
	indexPath string
}

// SearchResult represents a single search hit.
type SearchResult struct {
	Path  string  `json:"path"`
	Score float64 `json:"score"`
}

// OpenOrCreate opens an existing bleve index or creates a new one.
func OpenOrCreate(vaultPath string) (*Index, error) {
	indexPath := filepath.Join(vaultPath, ".vaultfs", "index.bleve")

	idx, err := bleve.Open(indexPath)
	if err != nil {
		// Create new index
		mapping := bleve.NewIndexMapping()
		idx, err = bleve.New(indexPath, mapping)
		if err != nil {
			return nil, err
		}
	}

	return &Index{index: idx, indexPath: indexPath}, nil
}

// Close closes the index.
func (i *Index) Close() error {
	return i.index.Close()
}

// Rebuild indexes all markdown files in the vault. Returns the number of documents indexed.
func (i *Index) Rebuild(vaultPath string) (int, error) {
	count := 0

	err := filepath.WalkDir(vaultPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(vaultPath, path)

		if strings.HasPrefix(rel, ".vaultfs") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		doc := map[string]any{
			"path":    filepath.ToSlash(rel),
			"content": string(data),
		}

		if err := i.index.Index(filepath.ToSlash(rel), doc); err != nil {
			return err
		}
		count++
		return nil
	})

	// Touch the index marker to track freshness
	touchMarker(vaultPath)

	return count, err
}

// Search performs a full-text query. Optionally filters by folder prefix.
// Uses AND semantics by default: "abc def" requires both terms.
// Use exact=true for phrase matching: "abc def" matches that exact sequence.
func (i *Index) Search(query, folder string, limit int, exact bool) ([]SearchResult, error) {
	var req *bleve.SearchRequest
	if exact {
		q := bleve.NewMatchPhraseQuery(query)
		req = bleve.NewSearchRequestOptions(q, limit, 0, false)
	} else {
		// Convert to AND: prefix each term with +
		terms := strings.Fields(query)
		for j, t := range terms {
			if !strings.HasPrefix(t, "+") && !strings.HasPrefix(t, "-") {
				terms[j] = "+" + t
			}
		}
		q := bleve.NewQueryStringQuery(strings.Join(terms, " "))
		req = bleve.NewSearchRequestOptions(q, limit, 0, false)
	}
	req.Fields = []string{"path"}

	searchResults, err := i.index.Search(req)
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, hit := range searchResults.Hits {
		path := hit.ID
		if folder != "" && !strings.HasPrefix(path, folder+"/") {
			continue
		}
		results = append(results, SearchResult{
			Path:  path,
			Score: hit.Score,
		})
	}

	return results, nil
}

// IsStale checks if the index needs rebuilding (most recent file mtime > marker mtime).
func IsStale(vaultPath string) (bool, error) {
	markerPath := filepath.Join(vaultPath, ".vaultfs", "index.marker")
	markerInfo, err := os.Stat(markerPath)
	if err != nil {
		return true, nil // No marker = stale
	}

	var latestMtime time.Time
	filepath.WalkDir(vaultPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(vaultPath, path)
		if strings.HasPrefix(rel, ".vaultfs") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.ModTime().After(latestMtime) {
			latestMtime = info.ModTime()
		}
		return nil
	})

	return latestMtime.After(markerInfo.ModTime()), nil
}

func touchMarker(vaultPath string) {
	markerPath := filepath.Join(vaultPath, ".vaultfs", "index.marker")
	os.WriteFile(markerPath, []byte(time.Now().Format(time.RFC3339)), 0644)
}
