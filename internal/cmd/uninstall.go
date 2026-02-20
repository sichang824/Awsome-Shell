package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove ~/.Awesome-Shell (Awesome Shell)",
	RunE:  runUninstall,
}

func runUninstall(cmd *cobra.Command, args []string) error {
	home, _ := os.UserHomeDir()
	dest := filepath.Join(home, ".Awesome-Shell")
	fmt.Println("Uninstalling Awesome Shell...")
	if err := os.RemoveAll(dest); err != nil {
		return err
	}
	fmt.Println("Uninstall complete.")
	return nil
}
