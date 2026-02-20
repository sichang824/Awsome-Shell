package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "as",
	Short: "Awesome Shell CLI toolkit",
	Long:  "Awesome Shell is a CLI toolkit for database management, devops and daily tasks.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("Awesome Shell version {{.Version}}\n")
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(passwordCmd)
	rootCmd.AddCommand(dockerCmd)
	rootCmd.AddCommand(pycacheCmd)
	rootCmd.AddCommand(ipCmd)
	rootCmd.AddCommand(sshCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(navicatCmd)
}
