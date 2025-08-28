package models

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type Rule struct {
	Method string `bson:"method"`
	URL    string `bson:"url"`
	Path   string `bson:"path"`
}

func (repo *DBRepo) GetRules(r map[string]*Rule, rs map[string]*Rule) error {
	collection := repo.Mongo.Database(database).Collection("http_rules")
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	for cur.Next(context.TODO()) {
		route := Rule{}
		cur.Decode(&route)
		ruleKey := route.Method + "|" + route.Path
		if strings.HasSuffix(route.Path, "*") {
			rs[strings.TrimSuffix(ruleKey, "/*")] = &route
		} else {
			r[ruleKey] = &route
		}
	}
	return nil
}
