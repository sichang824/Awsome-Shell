package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sichang824/awesome-shell/internal/config"
	"github.com/sichang824/awesome-shell/internal/db"
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
	// Read from env first so direnv/shell exports (and special chars like + in password) are not overwritten
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	user := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	pw := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
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
	if user == "" {
		user = mongoUser
	}
	if user == "" {
		user = config.GetEnv("MONGO_INITDB_ROOT_USERNAME", "root")
	}
	if pw == "" {
		pw = mongoPassword
	}
	if pw == "" {
		pw = config.GetEnv("MONGO_INITDB_ROOT_PASSWORD", "")
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
	mongoDbsCmd = &cobra.Command{
		Use:   "dbs",
		Short: "List databases",
		Args:  cobra.NoArgs,
		RunE:  runMongoDbs,
	}
	mongoUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "List users",
		Args:  cobra.NoArgs,
		RunE:  runMongoUsers,
	}
	mongoCollectionsCmd = &cobra.Command{
		Use:   "collections [database]",
		Short: "List collections in a database",
		Args:  cobra.ExactArgs(1),
		RunE:  runMongoCollections,
	}
	mongoClientCmd = &cobra.Command{
		Use:   "client",
		Short: "Connect to MongoDB (interactive, Go driver REPL)",
		RunE:  runMongoClient,
	}
	mongoLoginCmd = &cobra.Command{
		Use:   "login [username] [password]",
		Short: "Connect to MongoDB as user (args or --user/--password/env)",
		Args:  cobra.RangeArgs(0, 2),
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

func runMongoDbs(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	cfg := getMongoConfig()
	client, err := openMongo(ctx, cfg)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	list, err := client.ListDatabases(ctx, bson.M{})
	if err != nil {
		return err
	}
	for _, d := range list.Databases {
		fmt.Println(d.Name)
	}
	return nil
}

func runMongoUsers(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	cfg := getMongoConfig()
	client, err := openMongo(ctx, cfg)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	admin := client.Database("admin")
	var result struct {
		Users []struct {
			User string `bson:"user"`
		} `bson:"users"`
	}
	if err := admin.RunCommand(ctx, bson.D{{Key: "usersInfo", Value: 1}}).Decode(&result); err != nil {
		return err
	}
	for _, u := range result.Users {
		fmt.Println(u.User)
	}
	return nil
}

func runMongoCollections(cmd *cobra.Command, args []string) error {
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
	cursor, err := client.Database(database).ListCollections(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var doc struct {
			Name string `bson:"name"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return err
		}
		fmt.Println(doc.Name)
	}
	return cursor.Err()
}

func runMongoREPL(ctx context.Context, client *mongo.Client) error {
	scanner := bufio.NewScanner(os.Stdin)
	var currentDB string
	fmt.Fprintln(os.Stderr, "Go driver REPL (use <db>, show dbs, show collections, find <coll> [limit], \\q to quit)")
	for {
		if currentDB != "" {
			fmt.Fprintf(os.Stderr, "mongo:%s> ", currentDB)
		} else {
			fmt.Fprint(os.Stderr, "mongo> ")
		}
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "\\q" || strings.EqualFold(line, "quit") || strings.EqualFold(line, "exit") {
			break
		}
		parts := strings.Fields(line)
		cmd := strings.ToLower(parts[0])
		switch cmd {
		case "use":
			if len(parts) < 2 {
				fmt.Fprintln(os.Stderr, "usage: use <database>")
				continue
			}
			currentDB = parts[1]
			fmt.Fprintln(os.Stderr, "switched to db", currentDB)
		case "show":
			if len(parts) < 2 {
				fmt.Fprintln(os.Stderr, "usage: show dbs | show collections")
				continue
			}
			switch strings.ToLower(parts[1]) {
			case "dbs":
				list, err := client.ListDatabases(ctx, bson.M{})
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR:", err)
					continue
				}
				for _, d := range list.Databases {
					fmt.Println(d.Name)
				}
			case "collections":
				if currentDB == "" {
					fmt.Fprintln(os.Stderr, "ERROR: no database selected (use <db> first)")
					continue
				}
				cursor, err := client.Database(currentDB).ListCollections(ctx, bson.M{})
				if err != nil {
					fmt.Fprintln(os.Stderr, "ERROR:", err)
					continue
				}
				for cursor.Next(ctx) {
					var doc struct {
						Name string `bson:"name"`
					}
					if err := cursor.Decode(&doc); err != nil {
						fmt.Fprintln(os.Stderr, "ERROR:", err)
						break
					}
					fmt.Println(doc.Name)
				}
				cursor.Close(ctx)
			default:
				fmt.Fprintln(os.Stderr, "usage: show dbs | show collections")
			}
		case "find":
			if len(parts) < 2 {
				fmt.Fprintln(os.Stderr, "usage: find <collection> [limit]")
				continue
			}
			if currentDB == "" {
				fmt.Fprintln(os.Stderr, "ERROR: no database selected (use <db> first)")
				continue
			}
			collName := parts[1]
			limit := 20
			if len(parts) >= 3 {
				if n, err := strconv.Atoi(parts[2]); err == nil && n > 0 {
					limit = n
				}
			}
			cursor, err := client.Database(currentDB).Collection(collName).Find(ctx, bson.M{}, options.Find().SetLimit(int64(limit)))
			if err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				continue
			}
			for cursor.Next(ctx) {
				var doc bson.M
				if err := cursor.Decode(&doc); err != nil {
					fmt.Fprintln(os.Stderr, "ERROR:", err)
					break
				}
				raw, _ := json.Marshal(doc)
				fmt.Println(string(raw))
			}
			cursor.Close(ctx)
		default:
			fmt.Fprintln(os.Stderr, "unknown command. use: use <db>, show dbs, show collections, find <coll> [limit], \\q")
		}
	}
	return scanner.Err()
}

func runMongoClient(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	cfg := getMongoConfig()
	client, err := openMongo(ctx, cfg)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	return runMongoREPL(ctx, client)
}

func runMongoLogin(cmd *cobra.Command, args []string) error {
	cfg := getMongoConfig()
	if len(args) >= 1 {
		cfg.User = args[0]
	}
	if len(args) >= 2 {
		cfg.Password = args[1]
	}
	if cfg.User == "" {
		return fmt.Errorf("username required: pass [username] [password] or use --user and --password (or set MONGO_INITDB_ROOT_USERNAME / MONGO_INITDB_ROOT_PASSWORD)")
	}
	ctx := context.Background()
	client, err := openMongo(ctx, cfg)
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	return runMongoREPL(ctx, client)
}
