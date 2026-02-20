package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sichang824/awesome-shell/internal/exec"
	"github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Docker utilities",
}

var dockerRmNoneCmd = &cobra.Command{
	Use:   "rm-none-images",
	Short: "Remove Docker images with <none> tag",
	Long:  "List dangling images and optionally remove them after confirmation.",
	RunE:  runDockerRmNone,
}

func init() {
	dockerCmd.AddCommand(dockerRmNoneCmd)
}

func runDockerRmNone(cmd *cobra.Command, args []string) error {
	out, _, err := exec.Run("docker", "images", "-f", "dangling=true", "-q")
	if err != nil {
		return err
	}
	ids := strings.Fields(strings.TrimSpace(out))
	if len(ids) == 0 {
		fmt.Println("No dangling images found.")
		return nil
	}
	fmt.Println("Found dangling images:")
	listOut, _, _ := exec.Run("docker", "images", "-f", "dangling=true")
	fmt.Print(listOut)
	fmt.Print("Remove these images? (y/n): ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return nil
	}
	if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
		fmt.Println("Cancelled.")
		return nil
	}
	args = append([]string{"rmi"}, ids...)
	_, stderr, err := exec.Run("docker", args...)
	if err != nil {
		fmt.Fprint(os.Stderr, stderr)
		return err
	}
	fmt.Println("Done.")
	return nil
}
