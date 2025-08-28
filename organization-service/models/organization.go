package models

import (
	"context"
	"strings"

	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *DBRepo) GetOrgWithHostAddr(host string) *service.Organization {
	col := repo.Mongo.Database(database).Collection("organizations")
	doc := &service.Organization{}
	col.FindOne(context.TODO(), bson.M{"info.host": host}).Decode(doc)
	return doc
}

func (repo *DBRepo) GetOrgWithID(id string) (*service.Organization, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	col := repo.Mongo.Database(database).Collection("organizations")
	doc := &service.Organization{}
	col.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(doc)
	return doc, nil
}

func (repo *DBRepo) UpdateOrgWithHostAddr(host string, update utils.MapStringAny) error {
	col := repo.Mongo.Database(database).Collection("organizations")
	_, err := col.UpdateOne(context.TODO(), bson.D{{"info.host", host}}, bson.D{{"$set", update}})
	return err
}

func (repo *DBRepo) UpdateFormOrder(id string, data utils.MapStringAny) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	col := repo.Mongo.Database(database).Collection("organizations")
	updates := bson.D{}
	arrayFilters := options.ArrayFilters{Filters: make([]interface{}, 0, len(data))}
	for k, v := range data {
		replacedK := k
		if strings.HasPrefix(k, "_") {
			replacedK = strings.ReplaceAll(k, "_", "")
		}
		key := replacedK + "Elem"
		updates = append(updates, bson.E{"settings.forms.userMetadata.$[" + key + "].order", v})
		arrayFilters.Filters = append(arrayFilters.Filters, bson.M{key + ".key": k})
	}
	updateOpts := options.Update().SetArrayFilters(arrayFilters)
	_, err = col.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.D{{Key: "$set", Value: updates}}, updateOpts)
	return err
}

func (repo *DBRepo) UpdateFormInput(id, form string, data service.FormInput) (*mongo.UpdateResult, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	col := repo.Mongo.Database(database).Collection("organizations")
	filter := bson.M{
		"_id": bson.M{
			"$eq": objID,
		},
		"settings.forms." + form + ".key": bson.M{
			"$eq": data.Key,
		},
	}
	res, err := col.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"settings.forms.userMetadata.$": data}}) // Modify bson.D to bson.M
	return res, nil
}

/*
db.organizations.findOne(

	{ "_id": ObjectId("64540b79cb7c89d64e01ff45"), "settings.forms.userMetadata.key": "newsletter" },
	{ "settings.forms.userMetadata.$": 1 }

)
*/
func (repo *DBRepo) GetFormInputWithKey(id, form, key string) (*service.FormInput, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	//res := &utils.MapStringAny{}
	// col := repo.Mongo.Database(database).Collection("organizations")
	formField := "settings.forms." + form
	doc := &service.Organization{}
	col := repo.Mongo.Database(database).Collection("organizations")
	projection := bson.D{{formField + ".$", 1}}
	res := col.FindOne(context.TODO(), bson.M{"_id": objID, formField + ".key": key}, options.FindOne().SetProjection(projection)).Decode(doc)
	if res != nil {
		return &service.FormInput{}, nil
	}
	return &doc.Settings.Forms[form][0], nil
}

func (repo *DBRepo) InsertFormInput(id, form string, data service.FormInput) (*mongo.UpdateResult, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	col := repo.Mongo.Database(database).Collection("organizations")
	formField := "settings.forms." + form
	update := bson.M{"$push": bson.M{formField: data}}                      // Change type from bson.D to bson.M
	res, err := col.UpdateOne(context.TODO(), bson.M{"_id": objID}, update) // Modify the arguments to match the expected types
	return res, err
}

func (repo *DBRepo) DeleteFormInput(id, form, key string) (*mongo.UpdateResult, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	col := repo.Mongo.Database(database).Collection("organizations")
	formField := "settings.forms." + form
	delete := bson.M{"$pull": bson.M{formField: bson.M{"key": key}}}        // Change type from bson.D to bson.M
	res, err := col.UpdateOne(context.TODO(), bson.M{"_id": objID}, delete) // Modify the arguments to match the expected types
	return res, err
}

func (repo *DBRepo) GetFormWithFormName(id, form string) (*service.Organization, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	col := repo.Mongo.Database(database).Collection("organizations")
	doc := &service.Organization{}
	findOptions := options.FindOne().SetProjection(bson.M{"settings.forms." + form: 1})
	col.FindOne(context.TODO(), bson.M{"_id": objID}, findOptions).Decode(doc)
	return doc, nil
}

func (repo *DBRepo) InsertServiceSettings(id, service, key string, data []*pb.OrgServiceSetting) (*mongo.UpdateResult, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	col := repo.Mongo.Database(database).Collection("organizations")
	path := "settings." + service + "." + key
	//delete old settings
	delete := bson.M{"$unset": bson.M{path: ""}}
	_, err = col.UpdateOne(context.TODO(), filter, delete)
	if err != nil {
		return nil, err
	}
	opts := options.Update().SetUpsert(true)
	update := bson.M{"$set": bson.M{path: data}}                    // Use $set instead of $push for setting array data
	res, err := col.UpdateOne(context.TODO(), filter, update, opts) // Use filter instead of id string
	return res, err
}

/*
db.collection.updateOne(
    { "_id": ObjectId("6632dae84d14cabf2e0bf339") },
    {
        "$set": {
            "options.$[phoneElem].order": 4,
            "options.$[addressElem].order": 3
        }
    },
    {
        "arrayFilters": [
            { "phoneElem.key": "phone" },
            { "addressElem.key": "address" }
        ]
    }
);
*/
