package models

import (
	"context"

	"github.com/moly-space/molylibs/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Label               string             `bson:"label" json:"label"`
	Icon                string             `bson:"icon" json:"icon"`
	To                  string             `bson:"to" json:"to"`
	Order               int                `bson:"order" json:"order"`
	InternalDescription string             `bson:"internalDescription" json:"-"`
	Description         string             `bson:"description" json:"description"`
	Options             utils.MapStringAny `bson:"options" json:"options"`
	SubMenu             []Item             `bson:"subMenu" json:"subMenu"`
	When                []map[string]bool  `bson:"when" json:"when"`
	ElementClass        string             `bson:"elementClass" json:"elementClass"`
	ElementID           string             `bson:"elementId" json:"elementId"`
}

type Option struct {
	Key   string `json:"key"`
	Value string `json:"Value"`
}

func (repo *DBRepo) GetItem(id string) *Item {
	col := repo.Mongo.Database(database).Collection("items")
	doc := &Item{}
	objId, _ := primitive.ObjectIDFromHex(id)
	col.FindOne(context.TODO(), bson.D{{"_id", objId}}).Decode(doc)
	return doc
}
