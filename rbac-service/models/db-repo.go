package models

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

var database = os.Getenv("MONGO_DATABASE")

type DBRepo struct {
	Mongo *mongo.Client
}

type DatabaseRepo interface {
	GetModel(host, service, ctx string) *string
	//GetPolicy(host, service, ctx string) ([]string, error)
	GetPolicies(pType, host, service, ctx string) ([]Policy, error)
	GetModel2(host, service, ctx string) *Model
}
