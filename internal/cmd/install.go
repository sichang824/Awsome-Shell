package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [local|remote]",
	Short: "Install Awesome Shell to ~/.Awesome-Shell",
	Long:  "local: copy current dir; remote: git clone. Then configure shell (zsh/bash/fish).",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInstall,
}

func init() {
	installCmd.Flags().BoolP("force", "f", false, "overwrite existing directory without prompt")
}

func runInstall(cmd *cobra.Command, args []string) error {
	mode := "remote"
	if len(args) > 0 {
		mode = args[0]
	}
	home, _ := os.UserHomeDir()
	dest := filepath.Join(home, ".Awesome-Shell")

	fmt.Println("Installing Awesome Shell...")
	fmt.Println("Mode:", mode)

	if mode == "local" {
		if _, err := os.Stat(dest); err == nil {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Println("Directory already exists:", dest)
				fmt.Print("Remove and continue? (y/N): ")
				scanner := bufio.NewScanner(os.Stdin)
				if !scanner.Scan() {
					return nil
				}
				if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
					fmt.Println("Cancelled.")
					return nil
				}
			}
			if err := os.RemoveAll(dest); err != nil {
				return err
			}
		}
		cwd, _ := os.Getwd()
		root := cwd
		for {
			if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
				break
			}
			parent := filepath.Dir(root)
			if parent == root {
				root = cwd
				break
			}
			root = parent
		}
		if err := copyDir(root, dest); err != nil {
			return err
		}
		fmt.Println("Copied to", dest)
	} else {
		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("%s already exists; remove it first or use install local", dest)
		}
		c := exec.Command("git", "clone", "https://github.com/sichang824/Awsome-Shell.git", dest)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return err
		}
		fmt.Println("Cloned to", dest)
	}

	shell := filepath.Base(os.Getenv("SHELL"))
	configLine := "export AWESOME_SHELL_ROOT=$HOME/.Awesome-Shell"
	aliasLine := "alias as=\"$HOME/.Awesome-Shell/core/main.sh\""
	if shell == "" {
		shell = "zsh"
	}

	var configFile string
	switch shell {
	case "fish":
		configFile = filepath.Join(home, ".config", "fish", "config.fish")
		configLine = "set -gx AWESOME_SHELL_ROOT $HOME/.Awesome-Shell"
		aliasLine = "alias as \"$HOME/.Awesome-Shell/core/main.sh\""
	case "bash":
		configFile = filepath.Join(home, ".bashrc")
	case "zsh":
		configFile = filepath.Join(home, ".zshrc")
	default:
		fmt.Println("Shell", shell, "not auto-configured. Add manually:")
		fmt.Println("  export AWESOME_SHELL_ROOT=$HOME/.Awesome-Shell")
		fmt.Println("  alias as=$HOME/.Awesome-Shell/core/main.sh")
		fmt.Println("Install complete.")
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
		return err
	}
	if err := replaceOrAppendLine(configFile, `^export AWESOME_SHELL_ROOT=`, configLine); err != nil {
		return err
	}
	if err := replaceOrAppendLine(configFile, `^alias as=`, aliasLine); err != nil {
		return err
	}
	if shell == "fish" {
		replaceOrAppendLine(configFile, `^set -gx AWESOME_SHELL_ROOT`, configLine)
	}
	fmt.Println("Configured", configFile)
	fmt.Println("Install complete. Run 'as' to use (or source your config).")
	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Sync()
}

func replaceOrAppendLine(configFile, pattern, line string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return os.WriteFile(configFile, []byte(line+"\n"), 0644)
		}
		return err
	}
	lines := strings.Split(string(content), "\n")
	found := false
	for i, l := range lines {
		if re.MatchString(l) {
			lines[i] = line
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, line)
	}
	return os.WriteFile(configFile, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
