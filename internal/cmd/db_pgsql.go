package cmd

import (
	"fmt"
	"strings"

	"github.com/sichang824/awesome-shell/internal/config"
	"github.com/sichang824/awesome-shell/internal/exec"
	"github.com/spf13/cobra"
)

var pgsqlCmd = &cobra.Command{
	Use:   "pgsql",
	Short: "PostgreSQL commands (docker compose service: postgresql)",
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
		Short: "Connect to PostgreSQL (interactive)",
		RunE:  runPgsqlClient,
	}
	pgsqlLoginCmd = &cobra.Command{
		Use:   "login [username] [password]",
		Short: "Connect to PostgreSQL as user",
		Args:  cobra.ExactArgs(2),
		RunE:  runPgsqlLogin,
	}
)

func pgsqlUser() string { return config.GetEnv("PG_USER", "postgres") }

func runPgsqlCreateDB(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	owner, database := args[0], args[1]
	out, _, _ := exec.DockerComposeExec("postgresql", "psql", "-U", pgsqlUser(), "-t", "-A", "-c", "SELECT 1 FROM pg_roles WHERE rolname = '"+owner+"';")
	if strings.TrimSpace(out) == "" {
		fmt.Println("User '" + owner + "' does not exist.")
		return nil
	}
	out, _, _ = exec.DockerComposeExec("postgresql", "psql", "-U", pgsqlUser(), "-t", "-A", "-c", "SELECT 1 FROM pg_database WHERE datname='"+database+"';")
	if strings.TrimSpace(out) != "" {
		fmt.Println("Database '" + database + "' already exists.")
		return nil
	}
	exec.MustDockerOut("postgresql", "psql", "-U", pgsqlUser(), "-c",
		"CREATE DATABASE "+database+" WITH OWNER "+owner+" ENCODING 'UTF8' LC_COLLATE='en_US.utf8' LC_CTYPE='en_US.utf8' TEMPLATE=template0;")
	exec.MustDockerOut("postgresql", "psql", "-U", pgsqlUser(), "-c", "ALTER DATABASE "+database+" SET timezone TO 'Asia/Shanghai';")
	fmt.Println("Database '" + database + "' created.")
	return nil
}

func runPgsqlCreateUser(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	username := args[0]
	pw := genPassword()
	out, _, _ := exec.DockerComposeExec("postgresql", "psql", "-U", pgsqlUser(), "-t", "-A", "-c", "SELECT 1 FROM pg_roles WHERE rolname = '"+username+"';")
	if strings.TrimSpace(out) != "" {
		fmt.Println("User '" + username + "' already exists.")
		return nil
	}
	exec.MustDockerOut("postgresql", "psql", "-U", pgsqlUser(), "-c", "CREATE USER "+username+" WITH PASSWORD '"+pw+"';")
	fmt.Println("User:", username)
	fmt.Println("Password:", pw)
	fmt.Println("Save this password.")
	return nil
}

func runPgsqlDeleteDB(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	db := args[0]
	out, _, _ := exec.DockerComposeExec("postgresql", "psql", "-U", pgsqlUser(), "-t", "-A", "-c", "SELECT 1 FROM pg_database WHERE datname = '"+db+"';")
	if strings.TrimSpace(out) == "" {
		fmt.Println("Database '" + db + "' does not exist.")
		return nil
	}
	if !confirm("Type database name to confirm: ", db) {
		fmt.Println("Cancelled.")
		return nil
	}
	exec.MustDockerOut("postgresql", "psql", "-U", pgsqlUser(), "-c", "DROP DATABASE "+db+";")
	fmt.Println("Database deleted.")
	return nil
}

func runPgsqlDeleteUser(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	user := args[0]
	out, _, _ := exec.DockerComposeExec("postgresql", "psql", "-U", pgsqlUser(), "-t", "-A", "-c", "SELECT 1 FROM pg_roles WHERE rolname = '"+user+"';")
	if strings.TrimSpace(out) == "" {
		fmt.Println("User '" + user + "' does not exist.")
		return nil
	}
	if !confirm("Type username to confirm: ", user) {
		fmt.Println("Cancelled.")
		return nil
	}
	exec.MustDockerOut("postgresql", "psql", "-U", pgsqlUser(), "-c", "DROP USER "+user+";")
	fmt.Println("User deleted.")
	return nil
}

func runPgsqlGrant(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	db, user := args[0], args[1]
	exec.MustDockerOut("postgresql", "psql", "-U", pgsqlUser(), "-c", "GRANT ALL PRIVILEGES ON DATABASE \""+db+"\" TO \""+user+"\";")
	exec.MustDockerOut("postgresql", "psql", "-U", pgsqlUser(), "-c", "GRANT USAGE ON SCHEMA public TO \""+user+"\";")
	exec.MustDockerOut("postgresql", "psql", "-U", pgsqlUser(), "-c", "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO \""+user+"\";")
	exec.MustDockerOut("postgresql", "psql", "-U", pgsqlUser(), "-c", "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO \""+user+"\";")
	fmt.Println("Granted.")
	return nil
}

func runPgsqlClient(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	return exec.DockerComposeExecTTY("postgresql", "psql", "-U", pgsqlUser())
}

func runPgsqlLogin(cmd *cobra.Command, args []string) error {
	return exec.DockerComposeExecTTYWithEnv("postgresql", map[string]string{"PGPASSWORD": args[1]}, "psql", "-U", args[0])
}
