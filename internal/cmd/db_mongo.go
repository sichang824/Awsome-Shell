package cmd

import (
	"fmt"
	"strings"

	"github.com/sichang824/awesome-shell/internal/config"
	"github.com/sichang824/awesome-shell/internal/exec"
	"github.com/spf13/cobra"
)

var mongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "MongoDB commands (docker compose service: mongo)",
}

var (
	mongoCreateDBCmd = &cobra.Command{
		Use:   "create-db [database]",
		Short: "Create MongoDB database",
		Args:  cobra.ExactArgs(1),
		RunE:  runMongoCreateDB,
	}
	mongoCreateUserCmd = &cobra.Command{
		Use:   "create-user [username] [role] [database]",
		Short: "Create MongoDB user (role/database optional, default admin)",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runMongoCreateUser,
	}
	mongoDeleteDBCmd = &cobra.Command{
		Use:   "delete-db [database]",
		Short: "Delete MongoDB database",
		Args:  cobra.ExactArgs(1),
		RunE:  runMongoDeleteDB,
	}
	mongoDeleteUserCmd = &cobra.Command{
		Use:   "delete-user [username]",
		Short: "Delete MongoDB user",
		Args:  cobra.ExactArgs(1),
		RunE:  runMongoDeleteUser,
	}
	mongoGrantCmd = &cobra.Command{
		Use:   "grant [username] [role] [database]",
		Short: "Grant role on database to user",
		Args:  cobra.ExactArgs(3),
		RunE:  runMongoGrant,
	}
	mongoClientCmd = &cobra.Command{
		Use:   "client",
		Short: "Connect to MongoDB (interactive)",
		RunE:  runMongoClient,
	}
	mongoLoginCmd = &cobra.Command{
		Use:   "login [username] [password]",
		Short: "Connect to MongoDB as user",
		Args:  cobra.ExactArgs(2),
		RunE:  runMongoLogin,
	}
)

func mongoUser() string  { return config.GetEnv("MONGO_INITDB_ROOT_USERNAME", "root") }
func mongoPass() string { return config.GetEnv("MONGO_INITDB_ROOT_PASSWORD", "") }

func mongoEval(expr string) (string, string, error) {
	return exec.DockerComposeExec("mongo", "mongosh", "--quiet",
		"--username", mongoUser(), "--password", mongoPass(), "--eval", expr)
}

func runMongoCreateDB(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	db := args[0]
	out, _, _ := mongoEval("db.getMongo().getDBNames().indexOf('" + db + "') >= 0")
	if strings.TrimSpace(out) == "true" {
		fmt.Println("Database '" + db + "' already exists.")
		return nil
	}
	mongoEval("db.getSiblingDB('" + db + "').createCollection('init_collection')")
	mongoEval("db.getSiblingDB('" + db + "').init_collection.insertOne({ initialized: true })")
	fmt.Println("Database '" + db + "' created.")
	return nil
}

func runMongoCreateUser(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	username := args[0]
	role, database := "readWrite", "admin"
	if len(args) >= 3 {
		role, database = args[1], args[2]
	} else if len(args) >= 2 {
		role = args[1]
	}
	pw := genPassword()
	out, _, _ := mongoEval("db.getSiblingDB('admin').getUser('" + username + "')")
	if strings.TrimSpace(out) != "null" && out != "" {
		fmt.Println("User '" + username + "' already exists.")
		return nil
	}
	js := "db.getSiblingDB('admin').createUser({ user: '" + username + "', pwd: '" + pw + "', roles: [{ role: '" + role + "', db: '" + database + "' }] })"
	mongoEval(js)
	fmt.Println("User:", username)
	fmt.Println("Password:", pw)
	fmt.Println("Save this password.")
	return nil
}

func runMongoDeleteDB(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	db := args[0]
	out, _, _ := mongoEval("db.getMongo().getDBNames().indexOf('" + db + "') >= 0")
	if strings.TrimSpace(out) != "true" {
		fmt.Println("Database '" + db + "' does not exist.")
		return nil
	}
	mongoEval("db.getSiblingDB('" + db + "').dropDatabase()")
	fmt.Println("Database deleted.")
	return nil
}

func runMongoDeleteUser(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	user := args[0]
	out, _, _ := mongoEval("db.getSiblingDB('admin').getUser('" + user + "')")
	if strings.TrimSpace(out) == "null" || out == "" {
		fmt.Println("User '" + user + "' does not exist.")
		return nil
	}
	mongoEval("db.getSiblingDB('admin').dropUser('" + user + "')")
	fmt.Println("User deleted.")
	return nil
}

func runMongoGrant(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	user, role, db := args[0], args[1], args[2]
	mongoEval("db.getSiblingDB('admin').grantRolesToUser('" + user + "', [{ role: '" + role + "', db: '" + db + "' }])")
	fmt.Println("Granted.")
	return nil
}

func runMongoClient(cmd *cobra.Command, args []string) error {
	config.LoadEnv()
	return exec.DockerComposeExecTTY("mongo", "mongosh", "--username", mongoUser(), "--password", mongoPass())
}

func runMongoLogin(cmd *cobra.Command, args []string) error {
	return exec.DockerComposeExecTTY("mongo", "mongosh", "--username", args[0], "--password", args[1], "--authenticationDatabase", "admin")
}
