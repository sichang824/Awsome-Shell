package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/sichang824/awesome-shell/internal/config"
	"github.com/sichang824/awesome-shell/internal/db"
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
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	user := os.Getenv("PGUSER")
	if user == "" {
		user = os.Getenv("PG_USER")
	}
	pw := os.Getenv("PGPASSWORD")
	if pw == "" {
		pw = os.Getenv("PG_PASSWORD")
	}
	if pw == "" {
		pw = os.Getenv("PG_PASS")
	}
	config.LoadEnv()
	if host == "" {
		host = config.GetEnv("PGHOST", pgHost)
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = config.GetEnv("PGPORT", pgPort)
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = pgUser
	}
	if user == "" {
		user = config.GetEnv("PGUSER", config.GetEnv("PG_USER", "postgres"))
	}
	if pw == "" {
		pw = pgPassword
	}
	if pw == "" {
		pw = config.GetEnv("PG_PASSWORD", config.GetEnv("PGPASSWORD", config.GetEnv("PG_PASS", "")))
	}
	return db.PgConfig{
		Host:     host,
		Port:     port,
		User:     user,
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
	pgsqlDbsCmd = &cobra.Command{
		Use:   "dbs",
		Short: "List databases",
		Args:  cobra.NoArgs,
		RunE:  runPgsqlDbs,
	}
	pgsqlUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "List users",
		Args:  cobra.NoArgs,
		RunE:  runPgsqlUsers,
	}
	pgsqlTablesCmd = &cobra.Command{
		Use:   "tables [database]",
		Short: "List tables in a database",
		Args:  cobra.ExactArgs(1),
		RunE:  runPgsqlTables,
	}
	pgsqlClientCmd = &cobra.Command{
		Use:   "client",
		Short: "Connect to PostgreSQL (interactive, runs local psql)",
		RunE:  runPgsqlClient,
	}
	pgsqlLoginCmd = &cobra.Command{
		Use:   "login [username] [password]",
		Short: "Connect to PostgreSQL as user (username/password from args, or from --user/--password/env)",
		Args:  cobra.RangeArgs(0, 2),
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

func runPgsqlDbs(cmd *cobra.Command, args []string) error {
	cfg := getPgConfig()
	conn, err := openPg(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()
	rows, err := conn.Query("SELECT datname FROM pg_database WHERE datistemplate = false ORDER BY datname")
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

func runPgsqlUsers(cmd *cobra.Command, args []string) error {
	cfg := getPgConfig()
	conn, err := openPg(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()
	rows, err := conn.Query("SELECT usename FROM pg_user ORDER BY usename")
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

func runPgsqlTables(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "database"); err != nil {
		return err
	}
	database := args[0]
	cfg := getPgConfig()
	cfg.Database = database
	conn, err := openPg(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()
	rows, err := conn.Query("SELECT tablename FROM pg_tables WHERE schemaname = 'public' ORDER BY tablename")
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

func runPgREPL(conn *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	var buf strings.Builder
	fmt.Fprintln(os.Stderr, "Go driver REPL (\\q to quit)")
	for {
		if buf.Len() > 0 {
			fmt.Fprint(os.Stderr, "... ")
		} else {
			fmt.Fprint(os.Stderr, "pgsql> ")
		}
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		buf.WriteString(line)
		buf.WriteString("\n")
		trimmed := strings.TrimRightFunc(buf.String(), unicode.IsSpace)
		if trimmed == "" {
			buf.Reset()
			continue
		}
		if strings.TrimSpace(trimmed) == "\\q" || strings.EqualFold(trimmed, "quit") || strings.EqualFold(trimmed, "exit") {
			break
		}
		if !strings.HasSuffix(strings.TrimSpace(trimmed), ";") {
			continue
		}
		stmt := strings.TrimSuffix(trimmed, ";")
		stmt = strings.TrimRightFunc(stmt, unicode.IsSpace)
		buf.Reset()
		if stmt == "" {
			continue
		}
		rows, err := conn.Query(stmt)
		if err != nil {
			result, execErr := conn.Exec(stmt)
			if execErr != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				continue
			}
			affected, _ := result.RowsAffected()
			fmt.Println("OK", affected, "row(s) affected")
			continue
		}
		cols, _ := rows.Columns()
		if len(cols) == 0 {
			_ = rows.Close()
			continue
		}
		fmt.Println(strings.Join(cols, "\t"))
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		for rows.Next() {
			if err := rows.Scan(ptrs...); err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				break
			}
			parts := make([]string, len(cols))
			for i, v := range vals {
				if v == nil {
					parts[i] = "NULL"
				} else {
					parts[i] = fmt.Sprint(v)
				}
			}
			fmt.Println(strings.Join(parts, "\t"))
		}
		if err := rows.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
		}
		_ = rows.Close()
	}
	return scanner.Err()
}

func runPgsqlClient(cmd *cobra.Command, args []string) error {
	cfg := getPgConfig()
	conn, err := openPg(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()
	return runPgREPL(conn)
}

func runPgsqlLogin(cmd *cobra.Command, args []string) error {
	cfg := getPgConfig()
	user := cfg.User
	password := cfg.Password
	if len(args) >= 1 {
		user = args[0]
	}
	if len(args) >= 2 {
		password = args[1]
	}
	if user == "" {
		return fmt.Errorf("username required: pass [username] [password] or use --user and --password (or set PGUSER/PG_PASSWORD)")
	}
	cfg.User = user
	cfg.Password = password
	conn, err := openPg(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()
	return runPgREPL(conn)
}
