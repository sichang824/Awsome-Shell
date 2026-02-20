package cmd

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/sichang824/awesome-shell/internal/config"
	"github.com/sichang824/awesome-shell/internal/exec"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations (MySQL, PostgreSQL, MongoDB via docker compose)",
}

func init() {
	config.LoadEnv()
	dbCmd.AddCommand(dbGenPasswordCmd, mysqlCmd, pgsqlCmd, mongoCmd)
}

func genPassword() string {
	b := make([]byte, 32)
	rand.Read(b)
	s := base64.RawURLEncoding.EncodeToString(b)
	s = strings.TrimRight(s, "+/=")
	if len(s) > 32 {
		return s[:32]
	}
	return s
}

func confirm(prompt, expected string) bool {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}
	return strings.TrimSpace(scanner.Text()) == expected
}

// --- MySQL ---
var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "MySQL commands (docker compose service: mysql)",
}

func init() {
	mysqlCmd.AddCommand(
		mysqlCreateDBCmd, mysqlCreateUserCmd, mysqlDeleteDBCmd, mysqlDeleteUserCmd,
		mysqlGrantCmd, mysqlClientCmd, mysqlLoginCmd,
	)
	pgsqlCmd.AddCommand(
		pgsqlCreateDBCmd, pgsqlCreateUserCmd, pgsqlDeleteDBCmd, pgsqlDeleteUserCmd,
		pgsqlGrantCmd, pgsqlClientCmd, pgsqlLoginCmd,
	)
	mongoCmd.AddCommand(
		mongoCreateDBCmd, mongoCreateUserCmd, mongoDeleteDBCmd, mongoDeleteUserCmd,
		mongoGrantCmd, mongoClientCmd, mongoLoginCmd,
	)
}

var (
	mysqlCreateDBCmd = &cobra.Command{
		Use:   "create-db [database]",
		Short: "Create MySQL database",
		Args:  cobra.ExactArgs(1),
		RunE:  runMysqlCreateDB,
	}
	mysqlCreateUserCmd = &cobra.Command{
		Use:   "create-user [username]",
		Short: "Create MySQL user with generated password",
		Args:  cobra.ExactArgs(1),
		RunE:  runMysqlCreateUser,
	}
	mysqlDeleteDBCmd = &cobra.Command{
		Use:   "delete-db [database]",
		Short: "Delete MySQL database (with confirmation)",
		Args:  cobra.ExactArgs(1),
		RunE:  runMysqlDeleteDB,
	}
	mysqlDeleteUserCmd = &cobra.Command{
		Use:   "delete-user [username]",
		Short: "Delete MySQL user (with confirmation)",
		Args:  cobra.ExactArgs(1),
		RunE:  runMysqlDeleteUser,
	}
	mysqlGrantCmd = &cobra.Command{
		Use:   "grant [database] [username]",
		Short: "Grant all on database to user",
		Args:  cobra.ExactArgs(2),
		RunE:  runMysqlGrant,
	}
	mysqlClientCmd = &cobra.Command{
		Use:   "client",
		Short: "Connect to MySQL as root (interactive)",
		RunE:  runMysqlClient,
	}
	mysqlLoginCmd = &cobra.Command{
		Use:   "login [username] [password]",
		Short: "Connect to MySQL as user",
		Args:  cobra.ExactArgs(2),
		RunE:  runMysqlLogin,
	}
)

var dbGenPasswordCmd = &cobra.Command{
	Use:   "gen-password",
	Short: "Generate random password",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(genPassword())
		return nil
	},
}

func mysqlRootPW() string {
	return config.GetEnv("MYSQL_ROOT_PASSWORD", "")
}

func runMysqlCreateDB(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	db := args[0]
	out := exec.MustDockerOut("mysql", "mysql", "-uroot", "-p"+mysqlRootPW(), "-e", "SHOW DATABASES LIKE '"+db+"';")
	if strings.TrimSpace(out) != "" {
		fmt.Println("Database '" + db + "' already exists.")
		return nil
	}
	exec.MustDockerOut("mysql", "mysql", "-uroot", "-p"+mysqlRootPW(), "-e", "CREATE DATABASE IF NOT EXISTS "+db+";")
	fmt.Println("Database '" + db + "' created.")
	return nil
}

func runMysqlCreateUser(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	user := args[0]
	pw := genPassword()
	out := exec.MustDockerOut("mysql", "mysql", "-uroot", "-p"+mysqlRootPW(), "-e", "SELECT user FROM mysql.user WHERE user = '"+user+"';")
	if strings.Contains(out, user) {
		fmt.Println("User '" + user + "' already exists.")
		return nil
	}
	exec.MustDockerOut("mysql", "mysql", "-uroot", "-p"+mysqlRootPW(), "-e",
		"CREATE USER IF NOT EXISTS '"+user+"'@'%' IDENTIFIED BY '"+pw+"';")
	fmt.Println("User:", user)
	fmt.Println("Password:", pw)
	fmt.Println("Save this password.")
	return nil
}

func runMysqlDeleteDB(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	db := args[0]
	out := exec.MustDockerOut("mysql", "mysql", "-uroot", "-p"+mysqlRootPW(), "-e", "SHOW DATABASES LIKE '"+db+"';")
	if strings.TrimSpace(out) == "" {
		fmt.Println("Database '" + db + "' does not exist.")
		return nil
	}
	if !confirm("Type database name to confirm: ", db) {
		fmt.Println("Cancelled.")
		return nil
	}
	exec.MustDockerOut("mysql", "mysql", "-uroot", "-p"+mysqlRootPW(), "-e", "DROP DATABASE "+db+";")
	fmt.Println("Database deleted.")
	return nil
}

func runMysqlDeleteUser(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	user := args[0]
	out := exec.MustDockerOut("mysql", "mysql", "-uroot", "-p"+mysqlRootPW(), "-e", "SELECT user FROM mysql.user WHERE user = '"+user+"';")
	if !strings.Contains(out, user) {
		fmt.Println("User '" + user + "' does not exist.")
		return nil
	}
	if !confirm("Type username to confirm: ", user) {
		fmt.Println("Cancelled.")
		return nil
	}
	exec.MustDockerOut("mysql", "mysql", "-uroot", "-p"+mysqlRootPW(), "-e", "DROP USER '"+user+"'@'%';")
	fmt.Println("User deleted.")
	return nil
}

func runMysqlGrant(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	db, user := args[0], args[1]
	exec.MustDockerOut("mysql", "mysql", "-uroot", "-p"+mysqlRootPW(), "-e",
		"GRANT ALL PRIVILEGES ON "+db+".* TO '"+user+"'@'%'; FLUSH PRIVILEGES;")
	fmt.Println("Granted.")
	return nil
}

func runMysqlClient(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	return exec.DockerComposeExecTTY("mysql", "mysql", "-uroot", "-p"+mysqlRootPW())
}

func runMysqlLogin(cmd *cobra.Command, args []string) error {
	return exec.DockerComposeExecTTY("mysql", "mysql", "-u"+args[0], "-p"+args[1])
}
