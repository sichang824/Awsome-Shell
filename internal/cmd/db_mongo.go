package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sichang824/awesome-shell/internal/config"
	"github.com/sichang824/awesome-shell/internal/db"
	"github.com/sichang824/awesome-shell/internal/exec"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoHost, mongoPort, mongoUser, mongoPassword string
)

var mongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "MongoDB commands (native connection)",
}

func init() {
	mongoCmd.PersistentFlags().StringVar(&mongoHost, "host", "127.0.0.1", "MongoDB host")
	mongoCmd.PersistentFlags().StringVar(&mongoPort, "port", "27017", "MongoDB port")
	mongoCmd.PersistentFlags().StringVar(&mongoUser, "user", "root", "MongoDB user")
	mongoCmd.PersistentFlags().StringVar(&mongoPassword, "password", "", "MongoDB password (default from MONGO_INITDB_ROOT_PASSWORD env)")
}

func getMongoConfig() db.MongoConfig {
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	config.LoadEnv()
	if host == "" {
		host = config.GetEnv("MONGO_HOST", mongoHost)
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = config.GetEnv("MONGO_PORT", mongoPort)
	}
	if port == "" {
		port = "27017"
	}
	pw := mongoPassword
	if pw == "" {
		pw = config.GetEnv("MONGO_INITDB_ROOT_PASSWORD", "")
	}
	user := mongoUser
	if user == "" {
		user = config.GetEnv("MONGO_INITDB_ROOT_USERNAME", "root")
	}
	return db.MongoConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: pw,
	}
}

func openMongo(ctx context.Context, cfg db.MongoConfig) (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI(cfg.URI())
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}
	return client, nil
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
		Short: "Connect to MongoDB (interactive, runs local mongosh)",
		RunE:  runMongoClient,
	}
	mongoLoginCmd = &cobra.Command{
		Use:   "login [username] [password]",
		Short: "Connect to MongoDB as user",
		Args:  cobra.ExactArgs(2),
		RunE:  runMongoLogin,
	}
)

func mongoDBExists(ctx context.Context, client *mongo.Client, name string) (bool, error) {
	list, err := client.ListDatabases(ctx, bson.M{"name": name})
	if err != nil {
		return false, err
	}
	for _, d := range list.Databases {
		if d.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func runMongoCreateDB(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "database"); err != nil {
		return err
	}
	database := args[0]
	ctx := context.Background()
	cfg := getMongoConfig()
	client, err := openMongo(ctx, cfg)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	exists, err := mongoDBExists(ctx, client, database)
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("Database '" + database + "' already exists.")
		return nil
	}
	// Create DB by creating a collection and inserting one doc
	err = client.Database(database).CreateCollection(ctx, "init_collection")
	if err != nil {
		return err
	}
	_, err = client.Database(database).Collection("init_collection").InsertOne(ctx, bson.M{"initialized": true})
	if err != nil {
		return err
	}
	fmt.Println("Database '" + database + "' created.")
	return nil
}

func runMongoCreateUser(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "username"); err != nil {
		return err
	}
	username := args[0]
	role, database := "readWrite", "admin"
	if len(args) >= 3 {
		role, database = args[1], args[2]
	} else if len(args) >= 2 {
		role = args[1]
	}
	pw := genPassword()
	ctx := context.Background()
	cfg := getMongoConfig()
	client, err := openMongo(ctx, cfg)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	admin := client.Database("admin")
	var u bson.M
	err = admin.RunCommand(ctx, bson.D{{Key: "usersInfo", Value: username}}).Decode(&u)
	if err == nil {
		if users, ok := u["users"].(bson.A); ok && len(users) > 0 {
			fmt.Println("User '" + username + "' already exists.")
			return nil
		}
	}

	cmdDoc := bson.D{
		{Key: "createUser", Value: username},
		{Key: "pwd", Value: pw},
		{Key: "roles", Value: bson.A{bson.D{{Key: "role", Value: role}, {Key: "db", Value: database}}}},
	}
	if err := admin.RunCommand(ctx, cmdDoc).Err(); err != nil {
		return err
	}
	fmt.Println("User:", username)
	fmt.Println("Password:", pw)
	fmt.Println("Save this password.")
	return nil
}

func runMongoDeleteDB(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "database"); err != nil {
		return err
	}
	database := args[0]
	ctx := context.Background()
	cfg := getMongoConfig()
	client, err := openMongo(ctx, cfg)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	exists, err := mongoDBExists(ctx, client, database)
	if err != nil {
		return err
	}
	if !exists {
		fmt.Println("Database '" + database + "' does not exist.")
		return nil
	}
	if err := client.Database(database).Drop(ctx); err != nil {
		return err
	}
	fmt.Println("Database deleted.")
	return nil
}

func runMongoDeleteUser(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "username"); err != nil {
		return err
	}
	username := args[0]
	ctx := context.Background()
	cfg := getMongoConfig()
	client, err := openMongo(ctx, cfg)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	admin := client.Database("admin")
	var u bson.M
	err = admin.RunCommand(ctx, bson.D{{Key: "usersInfo", Value: username}}).Decode(&u)
	if err != nil {
		return err
	}
	if users, ok := u["users"].(bson.A); !ok || len(users) == 0 {
		fmt.Println("User '" + username + "' does not exist.")
		return nil
	}

	if err := admin.RunCommand(ctx, bson.D{{Key: "dropUser", Value: username}}).Err(); err != nil {
		return err
	}
	fmt.Println("User deleted.")
	return nil
}

func runMongoGrant(cmd *cobra.Command, args []string) error {
	if err := requireSafeIdent(args[0], "username"); err != nil {
		return err
	}
	if err := requireSafeIdent(args[2], "database"); err != nil {
		return err
	}
	username, role, database := args[0], args[1], args[2]
	ctx := context.Background()
	cfg := getMongoConfig()
	client, err := openMongo(ctx, cfg)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	admin := client.Database("admin")
	cmdDoc := bson.D{
		{Key: "grantRolesToUser", Value: username},
		{Key: "roles", Value: bson.A{bson.D{{Key: "role", Value: role}, {Key: "db", Value: database}}}},
	}
	if err := admin.RunCommand(ctx, cmdDoc).Err(); err != nil {
		return err
	}
	fmt.Println("Granted.")
	return nil
}

func runMongoClient(cmd *cobra.Command, args []string) error {
	cfg := getMongoConfig()
	argsCli := []string{"--host", cfg.Host, "--port", cfg.Port, "--username", cfg.User}
	if cfg.Password != "" {
		argsCli = append(argsCli, "--password", cfg.Password)
	}
	return exec.RunInherit("mongosh", argsCli...)
}

func runMongoLogin(cmd *cobra.Command, args []string) error {
	cfg := getMongoConfig()
	argsCli := []string{"--host", cfg.Host, "--port", cfg.Port, "--username", args[0], "--password", args[1], "--authenticationDatabase", "admin"}
	return exec.RunInherit("mongosh", argsCli...)
}
