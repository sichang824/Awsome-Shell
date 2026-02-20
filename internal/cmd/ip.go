package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Get public IP address",
}

var ipv4Cmd = &cobra.Command{
	Use:   "ipv4 [url]",
	Short: "Get IPv4 address",
	Long:  "Default URL: ipinfo.io/ip",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runIPv4,
}

var ipv6Cmd = &cobra.Command{
	Use:   "ipv6 [url]",
	Short: "Get IPv6 address",
	Long:  "Default URL: ifconfig.me",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runIPv6,
}

func init() {
	ipCmd.AddCommand(ipv4Cmd, ipv6Cmd)
}

func runIPv4(cmd *cobra.Command, args []string) error {
	url := "https://ipinfo.io/ip"
	if len(args) > 0 {
		url = args[0]
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b := make([]byte, 64)
	n, _ := resp.Body.Read(b)
	if n > 0 {
		fmt.Print(string(b[:n]))
	}
	return nil
}

func runIPv6(cmd *cobra.Command, args []string) error {
	url := "https://ifconfig.me"
	if len(args) > 0 {
		url = args[0]
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b := make([]byte, 64)
	n, _ := resp.Body.Read(b)
	if n > 0 {
		fmt.Print(string(b[:n]))
	}
	return nil
}
