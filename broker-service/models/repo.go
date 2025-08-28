package models

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

var database = os.Getenv("MONGO_DATABASE")

type DBRepo struct {
	Mongo *mongo.Client
}
