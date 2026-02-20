package cmd

import (
	"bufio"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sichang824/awesome-shell/internal/config"
	"github.com/sichang824/awesome-shell/internal/db"
	"github.com/sichang824/awesome-shell/internal/exec"
	"github.com/spf13/cobra"
	"github.com/go-sql-driver/mysql"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations (MySQL, PostgreSQL, MongoDB) via native connection",
	Long:  "Connect with --host, --port, --user, --password. No docker required.",
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

// safeIdent allows only alphanumeric and underscore for SQL identifiers to avoid injection.
var safeIdent = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func requireSafeIdent(name, kind string) error {
	if !safeIdent.MatchString(name) {
		return fmt.Errorf("invalid %s name (only letters, numbers, underscore allowed)", kind)
	}
	return nil
}

// --- MySQL ---
var (
	mysqlHost, mysqlPort, mysqlUser, mysqlPassword string
)

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "MySQL commands (native connection)",
}

func init() {
	mysqlCmd.PersistentFlags().StringVar(&mysqlHost, "host", "127.0.0.1", "MySQL host")
	mysqlCmd.PersistentFlags().StringVar(&mysqlPort, "port", "3306", "MySQL port")
	mysqlCmd.PersistentFlags().StringVar(&mysqlUser, "user", "root", "MySQL user")
	mysqlCmd.PersistentFlags().StringVar(&mysqlPassword, "password", "", "MySQL password (default from MYSQL_ROOT_PASSWORD env)")
	mysqlCmd.AddCommand(
		mysqlCreateDBCmd, mysqlCreateUserCmd, mysqlDeleteDBCmd, mysqlDeleteUserCmd,
		mysqlGrantCmd, mysqlDbsCmd, mysqlUsersCmd, mysqlTablesCmd,
		mysqlClientCmd, mysqlLoginCmd,
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

func getMySQLConfig() db.MySQLConfig {
	// Read from env first so direnv/shell exports (and special chars like +) are not overwritten by LoadEnv
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	pw := os.Getenv("MYSQL_ROOT_PASSWORD")
	config.LoadEnv()
	if host == "" {
		host = config.GetEnv("MYSQL_HOST", mysqlHost)
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = config.GetEnv("MYSQL_PORT", mysqlPort)
	}
	if port == "" {
		port = "3306"
	}
	if pw == "" {
		pw = mysqlPassword
	}
	if pw == "" {
		pw = config.GetEnv("MYSQL_ROOT_PASSWORD", "")
	}
	return db.MySQLConfig{
		Host:     host,
		Port:     port,
		User:     mysqlUser,
		Password: pw,
		Database: "",
	}
}

func openMySQL(cfg db.MySQLConfig) (*sql.DB, error) {
	mc := &mysql.Config{
		User:                 cfg.User,
		Passwd:               cfg.Password,
		Net:                  "tcp",
		Addr:                 cfg.Host + ":" + cfg.Port,
		DBName:               cfg.Database,
		AllowNativePasswords: true, // required for MariaDB / mysql_native_password
		AllowCleartextPasswords: true,
		TLSConfig:            "false", // skip TLS for local/docker
	}
	return sql.Open("mysql", mc.FormatDSN())
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
	mysqlDbsCmd = &cobra.Command{
		Use:   "dbs",
		Short: "List databases",
		Args:  cobra.NoArgs,
		RunE:  runMysqlDbs,
	}
	mysqlUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "List users",
		Args:  cobra.NoArgs,
		RunE:  runMysqlUsers,
	}
	mysqlTablesCmd = &cobra.Command{
		Use:   "tables [database]",
		Short: "List tables in a database",
		Args:  cobra.ExactArgs(1),
		RunE:  runMysqlTables,
	}
	mysqlClientCmd = &cobra.Command{
		Use:   "client",
		Short: "Connect to MySQL (interactive, runs local mysql client)",
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

func runMysqlCreateDB(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "database"); err != nil {
		return err
	}
	database := args[0]
	cfg := getMySQLConfig()
	conn, err := openMySQL(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	var exists int
	err = conn.QueryRow("SELECT 1 FROM information_schema.schemata WHERE schema_name = ?", database).Scan(&exists)
	if err == nil {
		fmt.Println("Database '" + database + "' already exists.")
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}
	_, err = conn.Exec("CREATE DATABASE IF NOT EXISTS `" + database + "`")
	if err != nil {
		return err
	}
	fmt.Println("Database '" + database + "' created.")
	return nil
}

func runMysqlCreateUser(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "username"); err != nil {
		return err
	}
	username := args[0]
	pw := genPassword()
	cfg := getMySQLConfig()
	conn, err := openMySQL(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	var count int
	err = conn.QueryRow("SELECT COUNT(*) FROM mysql.user WHERE user = ?", username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		fmt.Println("User '" + username + "' already exists.")
		return nil
	}
	// MySQL does not support placeholders for identifiers in CREATE USER
	pwEsc := strings.ReplaceAll(pw, "'", "''")
	_, err = conn.Exec("CREATE USER IF NOT EXISTS `" + username + "`@'%' IDENTIFIED BY '" + pwEsc + "'")
	if err != nil {
		return err
	}
	fmt.Println("User:", username)
	fmt.Println("Password:", pw)
	fmt.Println("Save this password.")
	return nil
}

func runMysqlDeleteDB(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "database"); err != nil {
		return err
	}
	database := args[0]
	cfg := getMySQLConfig()
	conn, err := openMySQL(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	var exists int
	err = conn.QueryRow("SELECT 1 FROM information_schema.schemata WHERE schema_name = ?", database).Scan(&exists)
	if err == sql.ErrNoRows {
		fmt.Println("Database '" + database + "' does not exist.")
		return nil
	}
	if err != nil {
		return err
	}
	if !confirm("Type database name to confirm: ", database) {
		fmt.Println("Cancelled.")
		return nil
	}
	_, err = conn.Exec("DROP DATABASE `" + database + "`")
	if err != nil {
		return err
	}
	fmt.Println("Database deleted.")
	return nil
}

func runMysqlDeleteUser(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "username"); err != nil {
		return err
	}
	username := args[0]
	cfg := getMySQLConfig()
	conn, err := openMySQL(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	var count int
	err = conn.QueryRow("SELECT COUNT(*) FROM mysql.user WHERE user = ?", username).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		fmt.Println("User '" + username + "' does not exist.")
		return nil
	}
	if !confirm("Type username to confirm: ", username) {
		fmt.Println("Cancelled.")
		return nil
	}
	_, err = conn.Exec("DROP USER `" + username + "`@'%'")
	if err != nil {
		return err
	}
	fmt.Println("User deleted.")
	return nil
}

func runMysqlGrant(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "database"); err != nil {
		return err
	}
	if err := requireSafeIdent(args[1], "username"); err != nil {
		return err
	}
	database, username := args[0], args[1]
	cfg := getMySQLConfig()
	conn, err := openMySQL(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec("GRANT ALL PRIVILEGES ON `" + database + "`.* TO `" + username + "`@'%'")
	if err != nil {
		return err
	}
	_, err = conn.Exec("FLUSH PRIVILEGES")
	if err != nil {
		return err
	}
	fmt.Println("Granted.")
	return nil
}

func runMysqlDbs(cmd *cobra.Command, args []string) error {
	cfg := getMySQLConfig()
	conn, err := openMySQL(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()
	rows, err := conn.Query("SELECT schema_name FROM information_schema.schemata ORDER BY schema_name")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		fmt.Println(name)
	}
	return rows.Err()
}

func runMysqlUsers(cmd *cobra.Command, args []string) error {
	cfg := getMySQLConfig()
	conn, err := openMySQL(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()
	rows, err := conn.Query("SELECT user, host FROM mysql.user ORDER BY user, host")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var user, host string
		if err := rows.Scan(&user, &host); err != nil {
			return err
		}
		fmt.Printf("%s@%s\n", user, host)
	}
	return rows.Err()
}

func runMysqlTables(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "database"); err != nil {
		return err
	}
	database := args[0]
	cfg := getMySQLConfig()
	cfg.Database = database
	conn, err := openMySQL(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()
	rows, err := conn.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = ? ORDER BY table_name", database)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		fmt.Println(name)
	}
	return rows.Err()
}

func runMysqlClient(cmd *cobra.Command, args []string) error {
	cfg := getMySQLConfig()
	argsCli := []string{"-h", cfg.Host, "-P", cfg.Port, "-u", cfg.User}
	if cfg.Password != "" {
		argsCli = append(argsCli, "-p"+cfg.Password)
	}
	return exec.RunInherit("mysql", argsCli...)
}

func runMysqlLogin(cmd *cobra.Command, args []string) error {
	cfg := getMySQLConfig()
	argsCli := []string{"-h", cfg.Host, "-P", cfg.Port, "-u", args[0], "-p" + args[1]}
	return exec.RunInherit("mysql", argsCli...)
}
