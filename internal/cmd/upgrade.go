package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade Awesome Shell (git pull in AWESOME_SHELL_ROOT)",
	RunE:  runUpgrade,
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	root := os.Getenv("AWESOME_SHELL_ROOT")
	if root == "" {
		home, _ := os.UserHomeDir()
		root = filepath.Join(home, ".Awesome-Shell")
	}
	if _, err := os.Stat(root); err != nil {
		return fmt.Errorf("directory does not exist: %s", root)
	}
	c := exec.Command("git", "-C", root, "fetch", "--all")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	c = exec.Command("git", "-C", root, "reset", "--hard", "origin/main")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	c = exec.Command("git", "-C", root, "pull")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
