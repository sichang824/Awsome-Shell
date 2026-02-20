package cmd

import (
	"bufio"
	"fmt"
	"os"
	osexec "os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate commit message with Ollama and optionally commit",
	Long:  "Uses Ollama to generate a commit message from staged changes. Optionally runs git commit.",
	RunE:  runCommit,
}

var commitModel string

func init() {
	commitCmd.Flags().StringVar(&commitModel, "model", "deepseek-r1:14b", "Ollama model ID")
}

func runCommit(cmd *cobra.Command, args []string) error {
	// git status -s
	statusOut, err := osexec.Command("git", "status", "-s").Output()
	if err != nil {
		return err
	}
	// git diff --staged
	diffOut, _ := osexec.Command("git", "diff", "--staged").Output()
	prompt := fmt.Sprintf("Changes:\n%s\n\nChange Contents:\n%s", string(statusOut), string(diffOut))
	// ollama run model generate "prompt"
	c := osexec.Command("ollama", "run", commitModel, "generate", prompt)
	c.Stderr = os.Stderr
	out, err := c.Output()
	if err != nil {
		return err
	}
	msg := strings.TrimSpace(string(out))
	fmt.Println("\nGenerated commit message:\n" + msg)
	fmt.Print("Commit now? (y/n) ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil
	}
	if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
		return nil
	}
	osexec.Command("git", "add", ".").Run()
	c2 := osexec.Command("git", "commit", "-m", msg)
	c2.Stdout = os.Stdout
	c2.Stderr = os.Stderr
	return c2.Run()
}
