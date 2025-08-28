// Package molylibs provides utility functions for Redis database connections.
package molylibs

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisDBNumber is a type alias for an integer representing a Redis database number.
type RedisDBNumber = int

// Constants representing different Redis databases.
const (
	RedisDBLongLife  RedisDBNumber = 0 // Database for long-lived data
	RedisDBShortLife RedisDBNumber = 1 // Database for short-lived data
)

// RedisClients is a map from Redis database numbers to Redis clients.
var RedisClients = make(map[RedisDBNumber]*redis.Client)

// NewRedisClient creates a new Redis client connected to the specified database.
// It returns the client and any error encountered.
func NewRedisClient(db RedisDBNumber) (*redis.Client, error) {
	// Get Redis connection details from environment variables
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")

	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Ping the Redis server to check the connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	// Return the client
	return client, nil
}

// Redis is a generic struct for Redis data.
// It includes the database number, expiration time, key, and data.
type Redis[T any] struct {
	DBNumber   RedisDBNumber // The Redis database number
	Expiration int           // The expiration time of the data in seconds
	key        *string       // The key of the data
	data       *T            // The data
}

// SetKey sets the key of the Redis data.
// It takes a variable number of string arguments and concatenates them to form the key.
func (r *Redis[T]) SetKey(args ...string) error {
	if r.key != nil {
		return errors.New("key already set")
	}
	key := strings.Join(args, ":")
	r.key = &key
	return nil
}

// GetClient retrieves the Redis client for the database number of the Redis data.
// If the client does not exist, it creates a new one.
func (r *Redis[T]) GetClient() *redis.Client {
	if RedisClients[r.DBNumber] == nil {
		RedisClients[r.DBNumber], _ = NewRedisClient(RedisDBNumber(r.DBNumber))
	}
	return RedisClients[r.DBNumber]
}

// Set sets the data of the Redis data in the Redis database.
// It serializes the data using the gob package and sets it in the database with the key and expiration time of the Redis data.
func (r *Redis[T]) Set(ctx context.Context, data *T) error {
	if r.key == nil {
		return errors.New("key not set")
	}
	r.data = data
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(r.data)
	if err != nil {
		return err
	}
	if r.Expiration == 0 {
		r.Expiration = 10
	}
	r.GetClient().Set(ctx, string(*r.key), buf.Bytes(), time.Duration(r.Expiration)*time.Minute)
	return nil
}

// Get retrieves the data of the Redis data from the Redis database.
// It gets the data from the database using the key of the Redis data and deserializes it using the gob package.
func (r *Redis[T]) Get(ctx context.Context) (*T, error) {
	if r.key == nil {
		return nil, errors.New("key not set")
	}
	data, err := r.GetClient().Get(ctx, *r.key).Bytes()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&r.data)
	if err != nil {
		return nil, err
	}
	return r.data, nil
}

// GetRedisClientInfo prints the information of the Redis client with the specified ID.
// It only prints the information of the client connected to the first database.
func GetRedisClientInfo(id string) {
	for i, v := range RedisClients {
		if i == 0 {
			info := v.Info(context.Background(), id)
			println(info)
		}
	}
}
