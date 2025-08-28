package models

import (
	"context"
	"fmt"
	"log"

	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRequest struct {
	service.User
	Password        string `bson:"-" json:"password" validate:"required"`
	ConfirmPassword string `bson:"-" json:"confirmPassword" validate:"eqfield=Password""`
}

func (u *UserRequest) HashPassword() (utils.Hashed, utils.Salt, error) {
	hashedPassword, salt, err := utils.HashPassword(u.Password)
	if err != nil {
		return "", "", err
	}
	return hashedPassword, salt, nil
}

func (repo *DBRepo) InsertUser(u *service.User) (*mongo.InsertOneResult, error) {
	coll := repo.Mongo.Database(database).Collection("users")
	er := utils.ErrorDetails{}
	data, err := bson.Marshal(u)
	if err != nil {
		return nil, err
	}
	res, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		log.Panic(err)
		er.Error = "Unable to process request"
		return nil, err
	}
	//to send auth service
	//u.Credentials.ObjectID = res.InsertedID.(primitive.ObjectID)
	return res, nil
}

func (repo *DBRepo) UpdateMetadataWithID(id string, update, delete service.Metadata) error {
	coll := repo.Mongo.Database(database).Collection("users")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", oid}}

	var updateBson, deleteBson utils.MapStringAny
	updateBson = make(utils.MapStringAny)
	deleteBson = make(utils.MapStringAny)
	if update != nil {
		for k, v := range update {
			updateBson[k] = v
		}
		//updateBson = bson.M{"$set": updateBson}
	} else {
		updateBson = bson.M{}
	}
	if delete != nil {
		for k, v := range delete {
			deleteBson[k] = v
		}
		//deleteBson = bson.M{"$unset": deleteBson}
	} else {
		deleteBson = bson.M{}
	}
	query := bson.M{"$set": updateBson, "$unset": deleteBson}
	_, err = coll.UpdateOne(context.TODO(), filter, query)
	if err != nil {
		return err
	}
	return nil
}

// !!
func (repo *DBRepo) GetUserWithID(id string) (*service.User, error) {
	coll := repo.Mongo.Database(database).Collection("users")
	var user service.User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	coll.FindOne(context.TODO(), bson.D{{"_id", oid}}).Decode(&user)
	return &user, nil
}

// !!
func (repo *DBRepo) GetUserWithIDAndOrgID(id, orgID string) (*service.User, error) {
	var user service.User
	//var userMap = make(map[string]any)
	coll := repo.Mongo.Database(database).Collection("users")
	objID, _ := primitive.ObjectIDFromHex(id)
	orgObjID, _ := primitive.ObjectIDFromHex(orgID)
	filter := bson.M{
		"_id":            objID,
		"organizationId": orgObjID,
	}
	coll.FindOne(context.TODO(), filter).Decode(&user)
	// if len(user.Organizations) > 1 {
	// 	return nil, fmt.Errorf("Expected 1 organization but has more. userID: %s orgID; %s", userID, orgID)
	// }
	return &user, nil
}

func (repo *DBRepo) UpdateAccountWithIDAndOrgID(id, orgID string, user service.User) error {
	coll := repo.Mongo.Database(database).Collection("users")

	objID, _ := primitive.ObjectIDFromHex(id)
	orgObjID, _ := primitive.ObjectIDFromHex(orgID)

	filter := bson.M{
		"_id":            objID,
		"organizationId": orgObjID,
	}
	update := bson.M{"$set": bson.M{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *DBRepo) UpdateEmailWithID(id string, email string) error {
	coll := repo.Mongo.Database(database).Collection("users")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", oid}}
	update := bson.D{{"$set", bson.D{{"email", email}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *DBRepo) GetUserWithEmailAndOrgID(email, orgObjectID string) (*service.User, error) {
	coll := repo.Mongo.Database(database).Collection("users")
	var user service.User
	orgID, _ := primitive.ObjectIDFromHex(orgObjectID)
	filter := bson.M{
		"email":          email,
		"organizationId": orgID,
	}
	err := coll.FindOne(context.TODO(), filter).Decode(&user)
	return &user, err
}

func (repo *DBRepo) GetUsersWithOrgID(orgID string, search utils.MapStringSlice[[]string], pn *utils.Pagination) ([]service.User, int64, error) {
	coll := repo.Mongo.Database(database).Collection("users")
	var users []service.User
	orgObjID, _ := primitive.ObjectIDFromHex(orgID)
	utils.TermDebugging(`search`, len(search))
	filter := bson.M{
		"organizationId": orgObjID,
	}
	if len(search) > 0 {
		filter["$and"] = []bson.M{}
		for key, values := range search {
			for _, value := range values {
				filter["$and"] = append(filter["$and"].([]bson.M), bson.M{key: bson.M{"$regex": primitive.Regex{Pattern: ".*" + value + ".*", Options: "i"}}})
			}
		}
	}
	total, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}
	cursor, err := coll.Find(context.TODO(), filter, pn.GetMongoOptions())
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (repo *DBRepo) UpdateUserWithIDAndOrgID(id, orgID string, data utils.MapStringAny) (*mongo.UpdateResult, error) {
	coll := repo.Mongo.Database(database).Collection("users")
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	oOrgID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", oID}, {"organizationId", oOrgID}}
	query := bson.M{"$set": data}
	return coll.UpdateOne(context.TODO(), filter, query)
}

func (repo *DBRepo) AddStore(id, ctx string, data service.StoreData) (*mongo.UpdateResult, error) {
	coll := repo.Mongo.Database(database).Collection("users")
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", oID}}
	query := bson.M{"$push": bson.M{"store." + ctx: data}}
	return coll.UpdateOne(context.TODO(), filter, query)
}

func (repo *DBRepo) GetStore(id, ctx, key string) ([]service.StoreData, error) {
	coll := repo.Mongo.Database(database).Collection("users")
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var filter bson.D
	var opts *options.FindOneOptions
	if key != "" {
		filter = bson.D{{"_id", oID}, {"store." + ctx, bson.M{"$elemMatch": bson.M{"key": key}}}}
		opts = options.FindOne().SetProjection(bson.M{"store." + ctx + ".$": 1})
	} else {
		filter = bson.D{{"_id", oID}}
		opts = options.FindOne().SetProjection(bson.M{"store." + ctx: 1})
	}
	var user = service.User{}
	err = coll.FindOne(context.TODO(), filter, opts).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user.Store[ctx], nil
}

func (repo *DBRepo) DeleteStore(id, ctx, key string) (*mongo.UpdateResult, error) {
	coll := repo.Mongo.Database(database).Collection("users")
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", oID}}
	query := bson.M{"$pull": bson.M{"store." + ctx: bson.M{"key": key}}}
	return coll.UpdateOne(context.TODO(), filter, query)
}

func (repo *DBRepo) DeleteUserWithID(id string) error {
	coll := repo.Mongo.Database(database).Collection("users")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", oid}}
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (repo *DBRepo) SetVerifiedWithID(in *pb.EmailVerifedRequest) error {
	fmt.Println("connected")
	coll := repo.Mongo.Database(database).Collection("users")
	oid, err := primitive.ObjectIDFromHex(in.ID)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", oid}}
	fields := bson.D{{"verified", true}}
	if in.Status != "" {
		fields = append(fields, bson.E{"status", in.Status})
	}
	if in.FirstName != "" {
		fields = append(fields, bson.E{"firstName", in.FirstName})
	}
	if in.LastName != "" {
		fields = append(fields, bson.E{"lastName", in.LastName})
	}
	update := bson.D{{"$set", fields}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

/*
db.users.aggregate([

	{
	    $match: {
	        "_id": ObjectID("650b4efdf04d408b0538309f") // Filter the corresponding document
	    }
	},
	{
	    $set: {
	        "organizations": {
	            $arrayElemAt: [
	                {
	                    $filter: {
	                        input: "$organizations",
	                        as: "org",
	                        cond: { $eq: ["$$org._id", ObjectID("64540b79cb7c89d64e01ff45")] } // Filter only elements with specific _id
	                    }
	                },
	                0 // Select the first element of the array
	            ]
	        }
	    }
	},
	{
	    $set: {
	        "organizations.metadata": {
	            $cond: {
	                if: { $isArray: "$organizations.metadata" },
	                then: {
	                    $filter: {
	                        input: "$organizations.metadata",
	                        as: "meta",
	                        cond: {
	                            $regexMatch: {
	                                input: "$$meta.key",
	                                regex: /^n/,
	                                //regex: , // Select keys that don't start with "_"
	                                options: "" // Regex options (e.g., "i" - case insensitive)
	                            }
	                        }
	                    }
	                },
	                else: "$organizations.metadata"
	            }
	        }
	    }
	}

]);
*/

// func (repo *DBRepo) GetUserWithOptionalMetadataKeyByUserIDAndOrgID(userObjectID, orgObjectID, metadataKey string) (*service.Metadata, error) {
// 	coll := repo.Mongo.Database(database).Collection("users")
// 	userID, _ := primitive.ObjectIDFromHex(userObjectID)
// 	orgID, _ := primitive.ObjectIDFromHex(orgObjectID)
// 	var user []service.User

// 	pipeline := mongo.Pipeline{
// 		{{"$match", bson.D{{"_id", userID}}}},
// 		{{"$set", bson.D{
// 			{"organizations", bson.D{
// 				{"$arrayElemAt", bson.A{
// 					bson.D{
// 						{"$filter", bson.D{
// 							{"input", "$organizations"},
// 							{"as", "org"},
// 							{"cond", bson.D{{"$eq", bson.A{"$$org._id", orgID}}}},
// 						}},
// 					},
// 					0, // Select the first element of the array
// 				}},
// 			}},
// 		}}},
// 		{{"$set", bson.D{
// 			{"organizations.metadata", bson.D{
// 				{"$cond", bson.D{
// 					{"if", bson.D{{"$isArray", "$organizations.metadata"}}},
// 					{"then", bson.D{
// 						{"$filter", bson.D{
// 							{"input", "$organizations.metadata"},
// 							{"as", "meta"},
// 							{"cond", bson.D{
// 								{"$regexMatch", bson.D{
// 									{"input", "$$meta.key"},
// 									{"regex", primitive.Regex{Pattern: metadataKey, Options: ""}}, // Select keys that start with "_"
// 								}},
// 							}},
// 						}},
// 					}},
// 					{"else", "$organizations.metadata"},
// 				}},
// 			}},
// 		}}},
// 	}
// 	cursor, err := coll.Aggregate(context.TODO(), pipeline)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cursor.Close(context.TODO())
// 	//var results []bson.M
// 	if err = cursor.All(context.TODO(), &user); err != nil {
// 		log.Fatal(err)
// 	}
// 	return nil, nil
// }

// !!

/*
db.users.updateOne(

	{
	    "_id": ObjectID("650b4efdf04d408b0538309f"),
	    "organizations._id": ObjectID("64540b79cb7c89d64e01ff45")
	},
	{
	    $push: {
	        "organizations.$.metadata": {
	            $each: [
	                { "key": "addme", "value": "alex kim" },
	                { "key": "addme2", "value": 1 }
	            ]
	        }
	    }
	}

);
*/
// func (repo *DBRepo) AddMetadataByUserIDAndOrgID(userObjectID, orgObjectID string, metadata []service.Metadata) error {
// 	collection := repo.Mongo.Database(database).Collection("users")
// 	userID, err := primitive.ObjectIDFromHex(userObjectID)
// 	if err != nil {
// 		return err
// 	}
// 	orgID, err := primitive.ObjectIDFromHex(orgObjectID)
// 	if err != nil {
// 		return err
// 	}
// 	filter := bson.M{
// 		"_id":               userID,
// 		"organizations._id": orgID,
// 	}

// 	update := bson.M{
// 		"$push": bson.M{
// 			"organizations.$.metadata": bson.M{
// 				"$each": metadata,
// 			},
// 		},
// 	}
// 	_, err = collection.UpdateOne(context.TODO(), filter, update)
// 	return err
// }

/*
db.users.updateOne(

	{
	    "_id": ObjectID("650b4efdf04d408b0538309f"),
	    "organizations._id": ObjectID("64540b79cb7c89d64e01ff45")
	},
	{
	    $pull: {
	        "organizations.$.metadata": {
	            "key": { $in: ["newsletter", "notification"] }
	        }
	    }
	}

);
*/
// func (repo *DBRepo) DeleteMetadataByUserIDOrgIDMetadataKeys(userObjectID, orgObjectID string, metadataKeys []string) error {
// 	collection := repo.Mongo.Database(database).Collection("users")
// 	userID, err := primitive.ObjectIDFromHex(userObjectID)
// 	if err != nil {
// 		return err
// 	}
// 	orgID, err := primitive.ObjectIDFromHex(orgObjectID)
// 	if err != nil {
// 		return err
// 	}
// 	filter := bson.M{
// 		"_id":               userID,
// 		"organizations._id": orgID,
// 	}
// 	update := bson.M{
// 		"$pull": bson.M{
// 			"organizations.$.metadata": bson.M{
// 				"key": bson.M{"$in": metadataKeys},
// 			},
// 		},
// 	}
// 	_, err = collection.UpdateOne(context.TODO(), filter, update)
// 	return err
// }

/*
db.collection.update(
   { "_id": ObjectID("650b4efdf04d408b0538309f"), "organizations._id": ObjectID("64540b79cb7c89d64e01ff45") },
   {
      $set: {
         "organizations.$.metadata.$[element].value": 111,
         "organizations.$.metadata.$[genderElement].value": "female"
      }
   },
   {
      arrayFilters: [
         { "element.key": "zipcode" },
         { "genderElement.key": "gender" }
      ]
   }
)
*/

// func (repo *DBRepo) UpdateMetadataByUserIDAndOrgID(userObjectID, orgObjectID string, metadata []Metadata) error {
// 	collection := repo.Mongo.Database(database).Collection("users")

// 	userID, _ := primitive.ObjectIDFromHex(userObjectID)
// 	orgID, _ := primitive.ObjectIDFromHex(orgObjectID)

// 	filter := bson.M{
// 		"_id":                    userID,
// 		"organizations._id":      orgID,
// 		"organizations.metadata": bson.M{"$exists": true},
// 	}

// 	arrayFilters := make([]any, len(metadata))
// 	updateBson := bson.M{}

// 	for i, v := range metadata {
// 		arrayFilterKey := fmt.Sprintf("%s.key", v.Key)
// 		updateKey := fmt.Sprintf("organizations.$.metadata.$[%s].value", v.Key)
// 		arrayFilters[i] = bson.D{{Key: arrayFilterKey, Value: v.Key}}
// 		updateBson[updateKey] = v.Value
// 	}

// 	update := bson.M{
// 		"$set": updateBson,
// 	}

// 	opts := options.Update().SetArrayFilters(options.ArrayFilters{
// 		Filters: arrayFilters})

// 	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

/*
db.users.update(
  { "organizations._id": ObjectID("64540b79cb7c89d64e01ff45") },
  {
    $push: {
      "organizations.$.metadata": {
        $each: [
          { "key": "age", "value": 20 },
          { "key": "gender", "value": "male" }
        ]
      }
    }
  }
);
*/

/*
db.users.aggregate([
    {
        $match: {
            "_id": ObjectID("650b4efdf04d408b0538309f") // Filter the corresponding document
        }
    },
    {
        $set: {
            "organizations": {
                $map: {
                    input: {
                        $filter: {
                            input: "$organizations",
                            as: "org",
                            cond: { $eq: ["$$org._id", ObjectID("64540b79cb7c89d64e01ff45")] } // Filter only elements with specific _id
                        }
                    },
                    as: "org",
                    in: {
                        $mergeObjects: [
                            "$$org",
                            {
                                metadata: {
                                    $cond: {
                                        if: { $isArray: "$$org.metadata" },
                                        then: {
                                            $filter: {
                                                input: "$$org.metadata",
                                                as: "meta",
                                                cond: { $eq: [{ $indexOfBytes: ["$$meta.key", "_"] }, -1] } // Filter items where key doesn't start with "_"
                                            }
                                        },
                                        else: "$$org.metadata"
                                    }
                                }
                            }
                        ]
                    }
                }
            }
        }
    }
]);

*/

/*
db.collection.update(

	{
	  "_id": ObjectID("650b4efdf04d408b0538309f"),
	  "organizations._id": ObjectID("64540b79cb7c89d64e01ff45")
	},
	{
	  $set: {
	    "organizations.$.firstName": "james"
	  }
	}

)
*/
//!!
