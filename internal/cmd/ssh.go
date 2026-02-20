package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH / Git public key utilities",
}

var (
	sshCreateKeyCmd = &cobra.Command{
		Use:   "create-key",
		Short: "Create SSH key pair and add to agent",
		RunE:  runSSHCreateKey,
	}
	sshCheckCmd = &cobra.Command{
		Use:   "check [github_domain]",
		Short: "Test SSH connection (e.g. git@github.com)",
		RunE:  runSSHCheck,
	}
	sshConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "Add Host to ~/.ssh/config.d/git",
		RunE:  runSSHConfig,
	}
	sshListKeysCmd = &cobra.Command{
		Use:   "list-keys",
		Short: "List keys in ssh-agent",
		RunE:  runSSHListKeys,
	}
)

func init() {
	sshCmd.AddCommand(sshCreateKeyCmd, sshCheckCmd, sshConfigCmd, sshListKeysCmd)
}

func prompt(p string) (string, error) {
	fmt.Print(p)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", scanner.Err()
	}
	return strings.TrimSpace(scanner.Text()), nil
}

func runSSHCreateKey(cmd *cobra.Command, args []string) error {
	fmt.Println("Create a new SSH key pair and add to agent.")
	email, err := prompt("Email: ")
	if err != nil || email == "" {
		return err
	}
	sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return err
	}
	filename := filepath.Join(sshDir, "id_rsa_"+time.Now().Format("20060102150405"))
	fmt.Println("Key path:", filename)
	// ssh-keygen -t rsa -b 4096 -C email -f filename -N ""
	c := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-C", email, "-f", filename, "-N", "")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	// ssh-add filename
	add := exec.Command("ssh-add", filename)
	add.Stdout = os.Stdout
	add.Stderr = os.Stderr
	return add.Run()
}

func runSSHCheck(cmd *cobra.Command, args []string) error {
	domain := "github.com"
	if len(args) > 0 {
		domain = args[0]
	}
	c := exec.Command("ssh", "-T", domain)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

func runSSHConfig(cmd *cobra.Command, args []string) error {
	filename, _ := prompt("Key file path: ")
	githubUser, _ := prompt("GitHub username: ")
	githubDomain, _ := prompt("GitHub domain (e.g. github.com): ")
	if filename == "" || githubUser == "" || githubDomain == "" {
		return fmt.Errorf("all fields required")
	}
	configDir := filepath.Join(os.Getenv("HOME"), ".ssh", "config.d")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}
	configFile := filepath.Join(configDir, "git")
	block := fmt.Sprintf("\nHost %s\n    HostName %s\n    User %s\n    Port 22\n    IdentityFile %s\n",
		githubDomain, githubDomain, githubUser, filename)
	f, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, err = f.WriteString(block)
	f.Close()
	if err != nil {
		return err
	}
	fmt.Println("Added to", configFile)
	return nil
}

func runSSHListKeys(cmd *cobra.Command, args []string) error {
	c := exec.Command("ssh-add", "-l")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
