package models

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

type Model struct {
	Service string
	Model   string `bson:"model"`
	Ctx     string `bson:"ctx"`
}

// org
// func (repo *DBRepo) GetModel(host, service, ctx string) string {
// 	model := &Model{}
// 	collection := repo.Mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("models")
// 	filter := bson.D{
// 		bson.E{Key: "service", Value: service},
// 		bson.E{Key: "host", Value: host},
// 	}
// 	if ctx != "" {
// 		filter = append(filter, bson.E{Key: "ctx", Value: ctx})
// 	}
// 	collection.FindOne(context.Background(), filter).Decode(model)
// 	return model.Model
// }

func (repo *DBRepo) GetModel(host, service, ctx string) *string {
	model := &Model{}
	collection := repo.Mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("models")

	filter := bson.D{
		{Key: "service", Value: service},
		{Key: "host", Value: host},
	}
	if ctx != "" {
		filter = append(filter, bson.E{Key: "ctx", Value: ctx})
	}

	err := collection.FindOne(context.Background(), filter).Decode(model)
	if err == nil {
		return &model.Model
	}

	filter = bson.D{
		{Key: "service", Value: service},
		{Key: "host", Value: bson.D{{Key: "$in", Value: bson.A{"", nil}}}},
	}
	if ctx != "" {
		filter = append(filter, bson.E{Key: "ctx", Value: ctx})
	}

	collection.FindOne(context.Background(), filter).Decode(model)
	return &model.Model
}

func (repo *DBRepo) GetModel2(host, service, ctx string) *Model {
	model := &Model{}
	collection := repo.Mongo.Database(os.Getenv("MONGO_DATABASE")).Collection("models")

	filter := bson.D{}
	if host != "" {
		filter = append(filter, bson.E{Key: "host", Value: host})
	} else {
		filter = append(filter, bson.E{Key: "host", Value: bson.D{{Key: "$in", Value: bson.A{"", nil}}}})
	}
	if service != "" {
		filter = append(filter, bson.E{Key: "service", Value: service})
	} else {
		filter = append(filter, bson.E{Key: "service", Value: bson.D{{Key: "$in", Value: bson.A{"", nil}}}})
	}
	if ctx != "" {
		filter = append(filter, bson.E{Key: "ctx", Value: ctx})
	} else {
		filter = append(filter, bson.E{Key: "ctx", Value: bson.D{{Key: "$in", Value: bson.A{"", nil}}}})
	}

	err := collection.FindOne(context.Background(), filter).Decode(model)
	if err == nil {
		return model
	}
	return nil
}
