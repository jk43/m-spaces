package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Hook struct {
	ID           string `bson:"_id"`
	Service      string `bson:"service"`
	Tag          string `bson:"tag"`
	Priority     int32  `bson:"priority"`
	CallbackType string `bson:"callbackType"`
	CallbackAddr string `bson:"callbackAddr"`
	CallbackFunc string `bson:"callbackFunc"`
}

func (repo *DBRepo) GetHooks(service string) ([]*Hook, error) {
	var hooks []*Hook
	collection := repo.Mongo.Database(database).Collection("hooks")
	// Create a FindOptions instance to specify the sort order
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"priority": 1})

	cursor, err := collection.Find(context.Background(), bson.M{"service": service}, findOptions)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var hook Hook
		err := cursor.Decode(&hook)
		if err != nil {
			return nil, err
		}
		hooks = append(hooks, &hook)
	}
	return hooks, nil
}
