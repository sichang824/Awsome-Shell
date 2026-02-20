package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var pycacheCmd = &cobra.Command{
	Use:   "pycache",
	Short: "Remove __pycache__ directories",
	Long:  "Find and remove __pycache__ directories, excluding venv and .venv.",
	RunE:  runPycache,
}

var pycacheDir string

func init() {
	pycacheCmd.Flags().StringVarP(&pycacheDir, "dir", "d", ".", "directory to search")
}

func runPycache(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		pycacheDir = args[0]
	}
	root, err := filepath.Abs(pycacheDir)
	if err != nil {
		return err
	}
	var removed int
	err = filepath.Walk(root, func(path string, info os.FileInfo, errWalk error) error {
		if errWalk != nil {
			return errWalk
		}
		if !info.IsDir() {
			return nil
		}
		if info.Name() != "__pycache__" {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		if rel == "" {
			rel = path
		}
		norm := filepath.ToSlash(rel) + "/"
		if strings.Contains(norm, "/venv/") || strings.Contains(norm, "/.venv/") || strings.HasPrefix(norm, "venv/") || strings.HasPrefix(norm, ".venv/") {
			return filepath.SkipDir
		}
		if err := os.RemoveAll(path); err != nil {
			return err
		}
		removed++
		fmt.Println(path)
		return filepath.SkipDir
	})
	if err != nil {
		return err
	}
	fmt.Printf("Removed %d __pycache__ director(y/ies).\n", removed)
	return nil
}

