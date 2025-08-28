// Package molylibs provides utility functions for logging and database connections.
package molylibs

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Ping is a simple function that prints "Got ping, out Pong" when called.
func Ping() {
	fmt.Println("Got ping, out Pong")
}

// Logger sets up a global logger with debug level and custom caller marshal function.
// It returns a pointer to the configured logger.
func Logger() *zerolog.Logger {
	// Set global log level to Debug
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// Customize the caller marshal function to only display the file name and line number
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	// Create a logger with service name and caller information
	logger := log.With().
		Str("service", os.Getenv("SERVICE_NAME")).
		Caller().
		Logger()
	return &logger
}

// Mongo sets up a MongoDB client with the URI and credentials from environment variables.
// It returns a pointer to the client and any error encountered.
func Mongo() (*mongo.Client, error) {
	// Set client options with MongoDB URI from environment variable
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	// Set authentication options with username from environment variable
	clientOptions.SetAuth(options.Credential{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Mysql sets up a MySQL database connection using the host, port, password, user, and database name from environment variables.
// It returns a pointer to the database and any error encountered.
func Mysql() (*gorm.DB, error) {
	// Get MySQL connection details from environment variables
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlDB := os.Getenv("MYSQL_DB")

	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDB)

	i := 0
	var db *gorm.DB
	var err error

	// Try to establish a connection to the MySQL database
	// If the connection fails, it will retry up to 3 times, waiting 3 seconds between each attempt
	for {
		i++
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			//Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		time.Sleep(time.Second * 3)
		if i >= 3 {
			break
		}
	}

	// Return the database connection and any error encountered
	return db, err
}
