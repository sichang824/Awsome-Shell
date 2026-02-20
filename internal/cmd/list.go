package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available bin commands (shell scripts in AWESOME_SHELL_ROOT/bin)",
	Long:  "When AWESOME_SHELL_ROOT is set, lists .sh files in bin/.",
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	root := os.Getenv("AWESOME_SHELL_ROOT")
	if root == "" {
		home, _ := os.UserHomeDir()
		root = filepath.Join(home, ".Awesome-Shell")
	}
	binDir := filepath.Join(root, "bin")
	entries, err := os.ReadDir(binDir)
	if err != nil {
		return err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && len(e.Name()) > 3 && e.Name()[len(e.Name())-3:] == ".sh" {
			names = append(names, e.Name()[:len(e.Name())-3])
		}
	}
	sort.Strings(names)
	fmt.Println("Available commands in", binDir+":")
	fmt.Println("----------------------------------------")
	for i, n := range names {
		fmt.Printf("%d   %s\n", i, n)
	}
	return nil
}
