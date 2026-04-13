package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
	"github.com/zarar/vaultfs/internal/output"
	"github.com/zarar/vaultfs/internal/search"
)

// SearchResultItem is the output structure for search results.
type SearchResultItem struct {
	Path  string  `json:"path"`
	Score float64 `json:"score"`
}

// ContextResult is a search result with matching lines.
type ContextResult struct {
	Path    string        `json:"path"`
	Matches []ContextLine `json:"matches"`
}

// ContextLine represents a matching line with context.
type ContextLine struct {
	Line    int    `json:"line"`
	Content string `json:"content"`
}

func ensureIndex(vaultPath string) (*search.Index, error) {
	stale, _ := search.IsStale(vaultPath)
	idx, err := search.OpenOrCreate(vaultPath)
	if err != nil {
		return nil, err
	}
	if stale {
		if _, err := idx.Rebuild(vaultPath); err != nil {
			idx.Close()
			return nil, err
		}
	}
	return idx, nil
}

func runSearch(vaultPath, query, folder string, limit int, fuzzyMode, exact bool) ([]SearchResultItem, error) {
	var results []SearchResultItem

	if fuzzyMode {
		// Fuzzy filename matching
		fuzzyResults, err := fuzzySearchFiles(vaultPath, query)
		if err != nil {
			return nil, err
		}
		for _, r := range fuzzyResults {
			if folder != "" && !strings.HasPrefix(r.Path, folder+"/") {
				continue
			}
			results = append(results, r)
			if len(results) >= limit {
				break
			}
		}
		return results, nil
	}

	// Full-text search via bleve
	idx, err := ensureIndex(vaultPath)
	if err != nil {
		return nil, err
	}
	defer idx.Close()

	bleveResults, err := idx.Search(query, folder, limit, exact)
	if err != nil {
		return nil, err
	}

	for _, r := range bleveResults {
		results = append(results, SearchResultItem{
			Path:  r.Path,
			Score: r.Score,
		})
	}

	return results, nil
}

func runSearchContext(vaultPath, query string, limit int) ([]ContextResult, error) {
	// Walk files and do line-by-line text matching for context
	queryLower := strings.ToLower(query)
	var results []ContextResult

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

		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()

		var matches []ContextLine
		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if strings.Contains(strings.ToLower(line), queryLower) {
				matches = append(matches, ContextLine{
					Line:    lineNum,
					Content: line,
				})
			}
		}

		if len(matches) > 0 {
			results = append(results, ContextResult{
				Path:    filepath.ToSlash(rel),
				Matches: matches,
			})
		}

		if len(results) >= limit {
			return filepath.SkipAll
		}

		return nil
	})

	return results, err
}

func fuzzySearchFiles(vaultPath, query string) ([]SearchResultItem, error) {
	var files []string

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
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}

	matches := fuzzy.Find(query, files)
	var results []SearchResultItem
	for _, m := range matches {
		results = append(results, SearchResultItem{
			Path:  files[m.Index],
			Score: float64(m.Score),
		})
	}
	return results, nil
}

func runIndexRebuild(vaultPath string) (int, error) {
	idx, err := search.OpenOrCreate(vaultPath)
	if err != nil {
		return 0, err
	}
	defer idx.Close()
	return idx.Rebuild(vaultPath)
}

// --- Cobra commands ---

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search vault content",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		folder, _ := cmd.Flags().GetString("folder")
		limit, _ := cmd.Flags().GetInt("limit")
		fuzzyFlag, _ := cmd.Flags().GetBool("fuzzy")
		exactFlag, _ := cmd.Flags().GetBool("exact")

		results, err := runSearch(vaultPath, args[0], folder, limit, fuzzyFlag, exactFlag)
		if err != nil {
			return err
		}

		format := output.ResolveFormat(formatFlag, true)
		if format == output.FormatJSON {
			return output.WriteJSON(os.Stdout, results)
		}
		for _, r := range results {
			fmt.Printf("%.2f  %s\n", r.Score, r.Path)
		}
		return nil
	},
}

var searchContextCmd = &cobra.Command{
	Use:   "search:context <query>",
	Short: "Search with matching line context",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		limit, _ := cmd.Flags().GetInt("limit")

		results, err := runSearchContext(vaultPath, args[0], limit)
		if err != nil {
			return err
		}

		format := output.ResolveFormat(formatFlag, true)
		if format == output.FormatJSON {
			return output.WriteJSON(os.Stdout, results)
		}
		for _, r := range results {
			for _, m := range r.Matches {
				fmt.Printf("%s:%d: %s\n", r.Path, m.Line, m.Content)
			}
		}
		return nil
	},
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index operations",
}

var indexRebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "Rebuild the search index",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		count, err := runIndexRebuild(vaultPath)
		if err != nil {
			return err
		}
		fmt.Printf("Index rebuilt: %d files indexed\n", count)
		return nil
	},
}

func init() {
	searchCmd.Flags().String("folder", "", "Filter by folder")
	searchCmd.Flags().Int("limit", 10, "Maximum results")
	searchCmd.Flags().Bool("fuzzy", false, "Use fuzzy filename matching")
	searchCmd.Flags().Bool("exact", false, "Match exact phrase instead of AND-ing terms")
	rootCmd.AddCommand(searchCmd)

	searchContextCmd.Flags().Int("limit", 10, "Maximum results")
	rootCmd.AddCommand(searchContextCmd)

	indexCmd.AddCommand(indexRebuildCmd)
	rootCmd.AddCommand(indexCmd)
}
