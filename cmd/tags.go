package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zarar/vaultfs/internal/markdown"
	"github.com/zarar/vaultfs/internal/output"
)

// TagResult represents a tag with optional count.
type TagResult struct {
	Name  string `json:"name"`
	Count int    `json:"count,omitempty"`
}

func runTags(vaultPath string, counts bool, sortBy string) ([]TagResult, error) {
	tagCounts := make(map[string]int)

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

		tags, _ := markdown.ExtractAllTags(data)
		for _, tag := range tags {
			tagCounts[tag]++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var results []TagResult
	for name, count := range tagCounts {
		r := TagResult{Name: name}
		if counts {
			r.Count = count
		}
		results = append(results, r)
	}

	if sortBy == "count" {
		sort.Slice(results, func(i, j int) bool {
			return results[i].Count > results[j].Count
		})
	} else {
		sort.Slice(results, func(i, j int) bool {
			return results[i].Name < results[j].Name
		})
	}

	return results, nil
}

func runTagByName(vaultPath, tagName string) ([]string, error) {
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

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		tags, _ := markdown.ExtractAllTags(data)
		for _, tag := range tags {
			if tag == tagName {
				files = append(files, filepath.ToSlash(rel))
				break
			}
		}
		return nil
	})

	return files, err
}

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags in the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		counts, _ := cmd.Flags().GetBool("counts")
		sortBy, _ := cmd.Flags().GetString("sort")
		results, err := runTags(vaultPath, counts, sortBy)
		if err != nil {
			return err
		}
		format := output.ResolveFormat(formatFlag, true)
		if format == output.FormatJSON {
			return output.WriteJSON(os.Stdout, results)
		}
		for _, r := range results {
			if counts {
				fmt.Printf("%s (%d)\n", r.Name, r.Count)
			} else {
				fmt.Println(r.Name)
			}
		}
		return nil
	},
}

var tagCmd = &cobra.Command{
	Use:   "tag <name>",
	Short: "List files with a specific tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		files, err := runTagByName(vaultPath, args[0])
		if err != nil {
			return err
		}
		format := output.ResolveFormat(formatFlag, true)
		if format == output.FormatJSON {
			return output.WriteJSON(os.Stdout, files)
		}
		for _, f := range files {
			fmt.Println(f)
		}
		return nil
	},
}

func init() {
	tagsCmd.Flags().Bool("counts", false, "Show usage counts")
	tagsCmd.Flags().String("sort", "", "Sort by: count")
	rootCmd.AddCommand(tagsCmd)
	rootCmd.AddCommand(tagCmd)
}
