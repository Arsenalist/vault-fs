package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zarar/vaultfs/internal/markdown"
	"github.com/zarar/vaultfs/internal/output"
)

// VaultTask extends markdown.Task with the file path.
type VaultTask struct {
	File     string   `json:"file"`
	Line     int      `json:"line"`
	Text     string   `json:"text"`
	Done     bool     `json:"done"`
	Priority string   `json:"priority,omitempty"`
	Due      string   `json:"due,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Mentions []string `json:"mentions,omitempty"`
}

func runTasks(vaultPath, filter, folder string) ([]VaultTask, error) {
	searchRoot := vaultPath
	if folder != "" {
		searchRoot = filepath.Join(vaultPath, folder)
	}

	var result []VaultTask

	err := filepath.WalkDir(searchRoot, func(path string, d os.DirEntry, err error) error {
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

		tasks := markdown.ExtractTasks(data)
		for _, task := range tasks {
			if filter == "pending" && task.Done {
				continue
			}
			if filter == "done" && !task.Done {
				continue
			}

			result = append(result, VaultTask{
				File:     filepath.ToSlash(rel),
				Line:     task.Line,
				Text:     task.Text,
				Done:     task.Done,
				Priority: task.Priority,
				Due:      task.Due,
				Tags:     task.Tags,
				Mentions: task.Mentions,
			})
		}
		return nil
	})

	return result, err
}

func runTaskToggle(vaultPath, path string, line int) error {
	fullPath := filepath.Join(vaultPath, path)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	if line < 1 || line > len(lines) {
		return fmt.Errorf("line %d out of range (1-%d)", line, len(lines))
	}

	idx := line - 1 // 0-indexed
	l := lines[idx]

	if strings.Contains(l, "- [ ]") {
		lines[idx] = strings.Replace(l, "- [ ]", "- [x]", 1)
	} else if strings.Contains(l, "- [x]") {
		lines[idx] = strings.Replace(l, "- [x]", "- [ ]", 1)
	} else if strings.Contains(l, "- [X]") {
		lines[idx] = strings.Replace(l, "- [X]", "- [ ]", 1)
	} else {
		return fmt.Errorf("line %d is not a task checkbox", line)
	}

	return os.WriteFile(fullPath, []byte(strings.Join(lines, "\n")), 0644)
}

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List tasks across the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		filter := ""
		if pending, _ := cmd.Flags().GetBool("pending"); pending {
			filter = "pending"
		} else if done, _ := cmd.Flags().GetBool("done"); done {
			filter = "done"
		}
		folder, _ := cmd.Flags().GetString("folder")

		tasks, err := runTasks(vaultPath, filter, folder)
		if err != nil {
			return err
		}
		format := output.ResolveFormat(formatFlag, true)
		if format == output.FormatJSON {
			return output.WriteJSON(os.Stdout, tasks)
		}
		for _, t := range tasks {
			status := "[ ]"
			if t.Done {
				status = "[x]"
			}
			fmt.Printf("%s %s  %s:%d\n", status, t.Text, t.File, t.Line)
		}
		return nil
	},
}

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Task operations",
}

var taskToggleCmd = &cobra.Command{
	Use:   "toggle <path>",
	Short: "Toggle a task checkbox",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultPath, err := resolveVault()
		if err != nil {
			return err
		}
		line, _ := cmd.Flags().GetInt("line")
		if err := runTaskToggle(vaultPath, args[0], line); err != nil {
			return err
		}
		fmt.Printf("Toggled task at %s:%d\n", args[0], line)
		return nil
	},
}

func init() {
	tasksCmd.Flags().Bool("pending", false, "Show only pending tasks")
	tasksCmd.Flags().Bool("done", false, "Show only completed tasks")
	tasksCmd.Flags().String("folder", "", "Filter by folder")
	rootCmd.AddCommand(tasksCmd)

	taskToggleCmd.Flags().Int("line", 0, "Line number of the task")
	taskCmd.AddCommand(taskToggleCmd)
	rootCmd.AddCommand(taskCmd)
}
