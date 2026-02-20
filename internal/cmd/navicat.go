package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var navicatCmd = &cobra.Command{
	Use:   "navicat-reset",
	Short: "Reset Navicat Premium trial (macOS only)",
	Long:  "Detects Navicat version and clears trial data. macOS only.",
	RunE:  runNavicatReset,
}

func runNavicatReset(cmd *cobra.Command, args []string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("navicat-reset is only supported on macOS")
	}
	plistPath := "/Applications/Navicat Premium.app/Contents/Info.plist"
	if _, err := os.Stat(plistPath); err != nil {
		return fmt.Errorf("Navicat Premium not found at %s", plistPath)
	}
	out, err := exec.Command("defaults", "read", plistPath).Output()
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`CFBundleShortVersionString\s*=\s*"([^"]+)`)
	matches := re.FindStringSubmatch(string(out))
	if len(matches) < 2 {
		return fmt.Errorf("could not detect version")
	}
	version := strings.Split(matches[1], ".")[0]
	fmt.Println("Detected Navicat Premium version", version)

	var prefsFile string
	switch version {
	case "17", "16":
		prefsFile = filepath.Join(os.Getenv("HOME"), "Library/Preferences/com.navicat.NavicatPremium.plist")
	case "15":
		prefsFile = filepath.Join(os.Getenv("HOME"), "Library/Preferences/com.prect.NavicatPremium15.plist")
	default:
		return fmt.Errorf("version %s not supported", version)
	}

	fmt.Print("Resetting trial...")
	prefsOut, _ := exec.Command("defaults", "read", prefsFile).Output()
	hashRe := regexp.MustCompile(`([0-9A-Z]{32})\s*=`)
	for _, m := range hashRe.FindAllStringSubmatch(string(prefsOut), -1) {
		if len(m) >= 2 {
			exec.Command("defaults", "delete", prefsFile, m[1]).Run()
		}
	}
	navicatDir := filepath.Join(os.Getenv("HOME"), "Library/Application Support/PremiumSoft CyberTech/Navicat CC/Navicat Premium")
	dirEntries, _ := os.ReadDir(navicatDir)
	for _, e := range dirEntries {
		if strings.HasPrefix(e.Name(), ".") && len(e.Name()) == 33 {
			_ = os.Remove(filepath.Join(navicatDir, e.Name()))
		}
	}
	fmt.Println(" Done")
	return nil
}
