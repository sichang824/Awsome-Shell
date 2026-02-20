package db

import "net/url"

// MySQLConfig holds MySQL connection parameters.
// DSN is built via go-sql-driver/mysql.Config in cmd layer (handles password encoding & auth).
type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// PgConfig holds PostgreSQL connection parameters.
type PgConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// DSN returns PostgreSQL connection string for lib/pq.
func (c *PgConfig) DSN() string {
	// postgres://user:password@host:port/dbname?sslmode=disable
	s := "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port
	if c.Database != "" {
		s += "/" + c.Database
	} else {
		s += "/postgres"
	}
	return s + "?sslmode=disable"
}

// MongoConfig holds MongoDB connection parameters.
type MongoConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

// URI returns MongoDB connection URI for mongo-driver.
// authSource=admin is required for root user (MONGO_INITDB_ROOT_*).
func (c *MongoConfig) URI() string {
	u := &url.URL{Scheme: "mongodb", Host: c.Host + ":" + c.Port, Path: "/", RawQuery: "authSource=admin"}
	if c.User != "" || c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	}
	return u.String()
}
