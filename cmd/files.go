package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zarar/vaultfs/internal/markdown"
	"github.com/zarar/vaultfs/internal/output"
)

// FileInfo represents a file listing entry.
type FileInfo struct {
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

// ReadResult represents the output of the read command.
type ReadResult struct {
	Path       string         `json:"path"`
	Properties map[string]any `json:"properties"`
	Body       string         `json:"body"`
	Modified   string         `json:"modified"`
	Size       int64          `json:"size"`
}

// --- Core functions (testable without cobra) ---

func runCreate(vaultPath, path, content string, appendMode bool) error {
	if !strings.HasSuffix(path, ".md") {
		path = path + ".md"
	}
	fullPath := filepath.Join(vaultPath, path)

	if !appendMode {
		if _, err := os.Stat(fullPath); err == nil {
			return fmt.Errorf("file already exists: %s", path)
		}
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	if appendMode {
		if _, err := os.Stat(fullPath); err == nil {
			f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = f.WriteString(content)
			return err
		}
	}

	return os.WriteFile(fullPath, []byte(content), 0644)
}

func runRead(vaultPath, path string) (*ReadResult, error) {
	fullPath := filepath.Join(vaultPath, path)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	fm, body, err := markdown.ParseFrontmatter(data)
	if err != nil {
		return nil, err
	}

	return &ReadResult{
		Path:       path,
		Properties: fm,
		Body:       string(body),
		Modified:   info.ModTime().Format(time.RFC3339),
		Size:       info.Size(),
	}, nil
}

func runAppend(vaultPath, path, content string) error {
	fullPath := filepath.Join(vaultPath, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return os.WriteFile(fullPath, []byte(content), 0644)
	}

	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

func runPrepend(vaultPath, path, content string) error {
	fullPath := filepath.Join(vaultPath, path)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	fm, body, err := markdown.ParseFrontmatter(data)
	if err != nil {
		return err
	}

	var result string
	if len(fm) > 0 {
		// Reconstruct: original frontmatter + new content + body
		fmEnd := findFrontmatterEnd(data)
		result = string(data[:fmEnd]) + "\n" + content + string(body)
	} else {
		result = content + string(data)
	}

	return os.WriteFile(fullPath, []byte(result), 0644)
}

// findFrontmatterEnd returns the byte offset just after the closing ---\n
func findFrontmatterEnd(data []byte) int {
	if len(data) < 3 || string(data[:3]) != "---" {
		return 0
	}
	// Skip first ---
	i := 3
	for i < len(data) && data[i] == '\n' || (i < len(data) && data[i] == '\r') {
		i++
	}
	// Find closing ---
	rest := string(data[i:])
	idx := strings.Index(rest, "---")
	if idx < 0 {
		return 0
	}
	end := i + idx + 3
	// Include the newline after closing ---
	if end < len(data) && data[end] == '\n' {
		end++
	}
	return end
}

func runMove(vaultPath, from, to string) error {
	fromPath := filepath.Join(vaultPath, from)
	toPath := filepath.Join(vaultPath, to)

	if err := os.MkdirAll(filepath.Dir(toPath), 0755); err != nil {
		return err
	}

	return os.Rename(fromPath, toPath)
}

func runDelete(vaultPath, path string) error {
	return os.Remove(filepath.Join(vaultPath, path))
}

func runList(vaultPath, folder, ext string) ([]FileInfo, error) {
	var files []FileInfo

	searchRoot := vaultPath
	if folder != "" {
		searchRoot = filepath.Join(vaultPath, folder)
	}

	if ext == "" {
		ext = "md"
	}

	err := filepath.WalkDir(searchRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(vaultPath, path)

		// Skip .vaultfs
		if strings.HasPrefix(rel, ".vaultfs") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), "."+ext) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		files = append(files, FileInfo{
			Path:     filepath.ToSlash(rel),
			Size:     info.Size(),
			Modified: info.ModTime().Format(time.RFC3339),
		})
		return nil
	})

	return files, err
}

func runFolders(vaultPath string) ([]string, error) {
	var folders []string

	err := filepath.WalkDir(vaultPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() || path == vaultPath {
			return nil
		}

		rel, _ := filepath.Rel(vaultPath, path)

		if strings.HasPrefix(rel, ".vaultfs") {
			return filepath.SkipDir
		}

		folders = append(folders, filepath.ToSlash(rel))
		return nil
	})

	return folders, err
}

func runMkdir(vaultPath, path string) error {
	return os.MkdirAll(filepath.Join(vaultPath, path), 0755)
}

// --- Cobra commands ---

var createCmd = &cobra.Command{
	Use:   "create <path>",
	Short: "Create a new markdown file",
	Long:  "Create a new markdown file. Omit .md extension (added automatically). Parent directories are created automatically.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		content, _ := cmd.Flags().GetString("content")
		appendFlag, _ := cmd.Flags().GetBool("append")
		if err := runCreate(vaultPath, args[0], content, appendFlag); err != nil {
			return err
		}
		path := args[0]
		if !strings.HasSuffix(path, ".md") {
			path += ".md"
		}
		fmt.Printf("Created %s\n", path)
		return nil
	},
}

var readCmd = &cobra.Command{
	Use:   "read <path>",
	Short: "Read a markdown file with parsed frontmatter",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		result, err := runRead(vaultPath, args[0])
		if err != nil {
			return err
		}
		format := output.ResolveFormat(formatFlag, true)
		if format == output.FormatJSON {
			return output.WriteJSON(os.Stdout, result)
		}
		fmt.Println(result.Body)
		return nil
	},
}

var appendCmd = &cobra.Command{
	Use:   "append <path>",
	Short: "Append content to a file (creates if missing)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		content, _ := cmd.Flags().GetString("content")
		if err := runAppend(vaultPath, args[0], content); err != nil {
			return err
		}
		fmt.Printf("Appended to %s\n", args[0])
		return nil
	},
}

var prependCmd = &cobra.Command{
	Use:   "prepend <path>",
	Short: "Prepend content after frontmatter",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		content, _ := cmd.Flags().GetString("content")
		if err := runPrepend(vaultPath, args[0], content); err != nil {
			return err
		}
		fmt.Printf("Prepended to %s\n", args[0])
		return nil
	},
}

var moveCmd = &cobra.Command{
	Use:   "move <path>",
	Short: "Move/rename a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		to, _ := cmd.Flags().GetString("to")
		if to == "" {
			return fmt.Errorf("--to flag is required")
		}
		if err := runMove(vaultPath, args[0], to); err != nil {
			return err
		}
		fmt.Printf("Moved %s → %s\n", args[0], to)
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <path>",
	Short: "Delete a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		if err := runDelete(vaultPath, args[0]); err != nil {
			return err
		}
		fmt.Printf("Deleted %s\n", args[0])
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List files in the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		folder, _ := cmd.Flags().GetString("folder")
		ext, _ := cmd.Flags().GetString("ext")
		files, err := runList(vaultPath, folder, ext)
		if err != nil {
			return err
		}
		format := output.ResolveFormat(formatFlag, true)
		if format == output.FormatJSON {
			return output.WriteJSON(os.Stdout, files)
		}
		for _, f := range files {
			fmt.Println(f.Path)
		}
		return nil
	},
}

var foldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "List all directories in the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		folders, err := runFolders(vaultPath)
		if err != nil {
			return err
		}
		format := output.ResolveFormat(formatFlag, true)
		if format == output.FormatJSON {
			return output.WriteJSON(os.Stdout, folders)
		}
		for _, f := range folders {
			fmt.Println(f)
		}
		return nil
	},
}

var mkdirCmd = &cobra.Command{
	Use:   "mkdir <path>",
	Short: "Create directories recursively",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		if err := runMkdir(vaultPath, args[0]); err != nil {
			return err
		}
		fmt.Printf("Created directory %s\n", args[0])
		return nil
	},
}

func init() {
	createCmd.Flags().String("content", "", "File content")
	createCmd.Flags().Bool("append", false, "Append to existing file instead of erroring")
	rootCmd.AddCommand(createCmd)

	rootCmd.AddCommand(readCmd)

	appendCmd.Flags().String("content", "", "Content to append")
	rootCmd.AddCommand(appendCmd)

	prependCmd.Flags().String("content", "", "Content to prepend")
	rootCmd.AddCommand(prependCmd)

	moveCmd.Flags().String("to", "", "Target path")
	rootCmd.AddCommand(moveCmd)

	rootCmd.AddCommand(deleteCmd)

	listCmd.Flags().String("folder", "", "Filter by folder")
	listCmd.Flags().String("ext", "", "Filter by extension (default: md)")
	rootCmd.AddCommand(listCmd)

	rootCmd.AddCommand(foldersCmd)
	rootCmd.AddCommand(mkdirCmd)
}
