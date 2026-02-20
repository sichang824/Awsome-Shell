package cmd

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/sichang824/awesome-shell/internal/config"
	"github.com/sichang824/awesome-shell/internal/db"
	"github.com/sichang824/awesome-shell/internal/exec"
	"github.com/spf13/cobra"
	_ "github.com/lib/pq"
)

var (
	pgHost, pgPort, pgUser, pgPassword string
)

var pgsqlCmd = &cobra.Command{
	Use:   "pgsql",
	Short: "PostgreSQL commands (native connection)",
}

func init() {
	pgsqlCmd.PersistentFlags().StringVar(&pgHost, "host", "127.0.0.1", "PostgreSQL host")
	pgsqlCmd.PersistentFlags().StringVar(&pgPort, "port", "5432", "PostgreSQL port")
	pgsqlCmd.PersistentFlags().StringVar(&pgUser, "user", "postgres", "PostgreSQL user")
	pgsqlCmd.PersistentFlags().StringVar(&pgPassword, "password", "", "PostgreSQL password (default from PG_PASSWORD env)")
}

func getPgConfig() db.PgConfig {
	config.LoadEnv()
	host := pgHost
	if v := config.GetEnv("PGHOST", ""); v != "" {
		host = v
	}
	port := pgPort
	if v := config.GetEnv("PGPORT", ""); v != "" {
		port = v
	}
	pw := pgPassword
	if pw == "" {
		pw = config.GetEnv("PG_PASSWORD", config.GetEnv("PGPASSWORD", ""))
	}
	return db.PgConfig{
		Host:     host,
		Port:     port,
		User:     pgUser,
		Password: pw,
		Database: "postgres",
	}
}

func openPg(cfg db.PgConfig) (*sql.DB, error) {
	return sql.Open("postgres", cfg.DSN())
}

var (
	pgsqlCreateDBCmd = &cobra.Command{
		Use:   "create-db [owner] [database]",
		Short: "Create PostgreSQL database",
		Args:  cobra.ExactArgs(2),
		RunE:  runPgsqlCreateDB,
	}
	pgsqlCreateUserCmd = &cobra.Command{
		Use:   "create-user [username]",
		Short: "Create PostgreSQL user with generated password",
		Args:  cobra.ExactArgs(1),
		RunE:  runPgsqlCreateUser,
	}
	pgsqlDeleteDBCmd = &cobra.Command{
		Use:   "delete-db [database]",
		Short: "Delete PostgreSQL database (with confirmation)",
		Args:  cobra.ExactArgs(1),
		RunE:  runPgsqlDeleteDB,
	}
	pgsqlDeleteUserCmd = &cobra.Command{
		Use:   "delete-user [username]",
		Short: "Delete PostgreSQL user (with confirmation)",
		Args:  cobra.ExactArgs(1),
		RunE:  runPgsqlDeleteUser,
	}
	pgsqlGrantCmd = &cobra.Command{
		Use:   "grant [database] [username]",
		Short: "Grant all on database to user",
		Args:  cobra.ExactArgs(2),
		RunE:  runPgsqlGrant,
	}
	pgsqlClientCmd = &cobra.Command{
		Use:   "client",
		Short: "Connect to PostgreSQL (interactive, runs local psql)",
		RunE:  runPgsqlClient,
	}
	pgsqlLoginCmd = &cobra.Command{
		Use:   "login [username] [password]",
		Short: "Connect to PostgreSQL as user",
		Args:  cobra.ExactArgs(2),
		RunE:  runPgsqlLogin,
	}
)

func runPgsqlCreateDB(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "owner"); err != nil {
		return err
	}
	if err := requireSafeIdent(args[1], "database"); err != nil {
		return err
	}
	owner, database := args[0], args[1]
	cfg := getPgConfig()
	conn, err := openPg(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	var exists int
	err = conn.QueryRow("SELECT 1 FROM pg_roles WHERE rolname = $1", owner).Scan(&exists)
	if err == sql.ErrNoRows {
		fmt.Println("User '" + owner + "' does not exist.")
		return nil
	}
	if err != nil {
		return err
	}
	err = conn.QueryRow("SELECT 1 FROM pg_database WHERE datname = $1", database).Scan(&exists)
	if err == nil {
		fmt.Println("Database '" + database + "' already exists.")
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}
	// Identifiers in PostgreSQL: use quote_ident or safe concat; we validated with safeIdent
	_, err = conn.Exec(`CREATE DATABASE "` + database + `" WITH OWNER "` + owner + `" ENCODING 'UTF8' LC_COLLATE='en_US.utf8' LC_CTYPE='en_US.utf8' TEMPLATE=template0`)
	if err != nil {
		return err
	}
	_, err = conn.Exec(`ALTER DATABASE "` + database + `" SET timezone TO 'Asia/Shanghai'`)
	if err != nil {
		return err
	}
	fmt.Println("Database '" + database + "' created.")
	return nil
}

func runPgsqlCreateUser(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "username"); err != nil {
		return err
	}
	username := args[0]
	pw := genPassword()
	cfg := getPgConfig()
	conn, err := openPg(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	var exists int
	err = conn.QueryRow("SELECT 1 FROM pg_roles WHERE rolname = $1", username).Scan(&exists)
	if err == nil {
		fmt.Println("User '" + username + "' already exists.")
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}
	pwEsc := strings.ReplaceAll(pw, "'", "''")
	_, err = conn.Exec("CREATE USER \"" + username + "\" WITH PASSWORD '" + pwEsc + "'")
	if err != nil {
		return err
	}
	fmt.Println("User:", username)
	fmt.Println("Password:", pw)
	fmt.Println("Save this password.")
	return nil
}

func runPgsqlDeleteDB(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "database"); err != nil {
		return err
	}
	database := args[0]
	cfg := getPgConfig()
	conn, err := openPg(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	var exists int
	err = conn.QueryRow("SELECT 1 FROM pg_database WHERE datname = $1", database).Scan(&exists)
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
	_, err = conn.Exec(`DROP DATABASE "` + database + `"`)
	if err != nil {
		return err
	}
	fmt.Println("Database deleted.")
	return nil
}

func runPgsqlDeleteUser(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "username"); err != nil {
		return err
	}
	username := args[0]
	cfg := getPgConfig()
	conn, err := openPg(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	var exists int
	err = conn.QueryRow("SELECT 1 FROM pg_roles WHERE rolname = $1", username).Scan(&exists)
	if err == sql.ErrNoRows {
		fmt.Println("User '" + username + "' does not exist.")
		return nil
	}
	if err != nil {
		return err
	}
	if !confirm("Type username to confirm: ", username) {
		fmt.Println("Cancelled.")
		return nil
	}
	_, err = conn.Exec(`DROP USER "` + username + `"`)
	if err != nil {
		return err
	}
	fmt.Println("User deleted.")
	return nil
}

func runPgsqlGrant(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "database"); err != nil {
		return err
	}
	if err := requireSafeIdent(args[1], "username"); err != nil {
		return err
	}
	database, username := args[0], args[1]
	cfg := getPgConfig()
	conn, err := openPg(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	for _, q := range []string{
		`GRANT ALL PRIVILEGES ON DATABASE "` + database + `" TO "` + username + `"`,
		`GRANT USAGE ON SCHEMA public TO "` + username + `"`,
		`GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "` + username + `"`,
		`GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "` + username + `"`,
	} {
		_, err = conn.Exec(q)
		if err != nil {
			return err
		}
	}
	fmt.Println("Granted.")
	return nil
}

func runPgsqlClient(cmd *cobra.Command, args []string) error {
	cfg := getPgConfig()
	argsCli := []string{"-h", cfg.Host, "-p", cfg.Port, "-U", cfg.User}
	if cfg.Password != "" {
		return exec.RunInheritWithEnv(map[string]string{"PGPASSWORD": cfg.Password}, "psql", argsCli...)
	}
	return exec.RunInherit("psql", argsCli...)
}

func runPgsqlLogin(cmd *cobra.Command, args []string) error {
	cfg := getPgConfig()
	argsCli := []string{"-h", cfg.Host, "-p", cfg.Port, "-U", args[0]}
	return exec.RunInheritWithEnv(map[string]string{"PGPASSWORD": args[1]}, "psql", argsCli...)
}
