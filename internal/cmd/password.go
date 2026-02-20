package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var passwordLength int

var passwordCmd = &cobra.Command{
	Use:   "password [length]",
	Short: "Generate random password",
	Long:  "Generate a random password using crypto/rand, default length 32.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runPassword,
}

func init() {
	passwordCmd.Flags().IntVarP(&passwordLength, "length", "n", 32, "password length")
}

func runPassword(cmd *cobra.Command, args []string) error {
	n := passwordLength
	if len(args) > 0 {
		if _, err := fmt.Sscanf(args[0], "%d", &n); err != nil {
			n = 32
		}
	}
	if n <= 0 {
		n = 32
	}
	// generate more bytes then trim to avoid weak endings
	b := make([]byte, (n*3/2)+8)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	s := base64.RawURLEncoding.EncodeToString(b)
	s = strings.TrimRight(s, "+/=")
	if len(s) > n {
		s = s[:n]
	}
	fmt.Println(s)
	return nil
}
