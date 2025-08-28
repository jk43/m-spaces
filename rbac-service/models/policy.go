package models

import (
	"context"
	"fmt"

	"github.com/moly-space/molylibs/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Policy struct {
	Service string
	Type    string `bson:"type"`
	Subject string `bson:"sub"`
	Object  string `bson:"obj"`
	Action  string `bson:"act"`
}

// func (repo *DBRepo) GetPolicy(host, service, ctx string) ([]string, error) {
// 	var pols []Policy
// 	collection := repo.Mongo.Database(database).Collection("policies")
// 	// service must be not empty
// 	filters := bson.D{
// 		bson.E{Key: "service", Value: service},
// 	}
// 	if ctx != "" {
// 		d := bson.E{Key: "ctx", Value: ctx}
// 		filters = append(filters, d)
// 	}
// 	withHost := bson.E{Key: "host", Value: host}
// 	filtersWithHost := append(filters, withHost)
// 	cur, err := collection.Find(context.Background(), filtersWithHost)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// check if there are documents
// 	hasDocuments := cur.TryNext(context.Background())
// 	if !hasDocuments {
// 		// try without host
// 		withoutHost := bson.E{Key: "host", Value: bson.D{bson.E{Key: "$in", Value: []any{"", nil}}}}
// 		filtersWithoutHost := append(filters, withoutHost)
// 		cur, err = collection.Find(context.Background(), filtersWithoutHost)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	var res []string
// 	for cur.Next(context.TODO()) {
// 		pol := Policy{}
// 		cur.Decode(&pol)
// 		pols = append(pols, pol)
// 		if pol.Action == "" {
// 			res = append(res, fmt.Sprintf("%s,%s,%s", pol.Type, pol.Subject, pol.Object))
// 		} else {
// 			res = append(res, fmt.Sprintf("%s,%s,%s,%s", pol.Type, pol.Subject, pol.Object, pol.Action))
// 		}
// 	}
// 	return res, nil
// }

func (repo *DBRepo) GetPolicies(pType, host, service, ctx string) ([]Policy, error) {
	var pols []Policy
	collection := repo.Mongo.Database(database).Collection("policies")
	// service must be not empty
	filters := bson.D{
		bson.E{Key: "type", Value: pType},
	}
	// service
	serviceFilter := bson.E{Key: "service", Value: service}
	if service == "" {
		serviceFilter = bson.E{Key: "service", Value: bson.D{bson.E{Key: "$in", Value: []any{"", nil}}}}
	}
	filters = append(filters, serviceFilter)
	// host
	var hostFilter bson.E
	if host != "" {
		hostFilter = bson.E{Key: "host", Value: host}
	} else {
		hostFilter = bson.E{Key: "host", Value: bson.D{bson.E{Key: "$in", Value: []any{"", nil}}}}
	}
	if ctx != "" {
		d := bson.E{Key: "ctx", Value: ctx}
		filters = append(filters, d)
	}
	filters = append(filters, hostFilter)
	//utils.TermDebugging(`filters`, filters)

	// withHost := bson.E{Key: "host", Value: host}
	// filtersWithHost := append(filters, withHost)
	cur, err := collection.Find(context.Background(), filters)
	// pol := Policy{}
	// //first document
	// cur.Decode(&pol)
	// pols = append(pols, pol)

	//utils.TermDebugging(`cur`, cur)
	if err != nil {
		return nil, err
	}
	// check if there are documents
	//hasDocuments := cur.TryNext(context.Background())
	// if !hasDocuments {
	// 	return nil, nil
	// }
	// if !hasDocuments {
	// 	// try without host
	// 	withoutHost := bson.E{Key: "host", Value: bson.D{bson.E{Key: "$in", Value: []any{"", nil}}}}
	// 	filtersWithoutHost := append(filters, withoutHost)
	// 	cur, err = collection.Find(context.Background(), filtersWithoutHost)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	for cur.Next(context.TODO()) {
		pol := Policy{}
		cur.Decode(&pol)
		pols = append(pols, pol)
	}
	return pols, nil
}

func (repo *DBRepo) GetGroups(host, service, ctx string) ([]string, error) {
	var pols []Policy
	collection := repo.Mongo.Database(database).Collection("policies")
	// service must be not empty
	var filters bson.D
	if service != "" {
		filters = bson.D{
			bson.E{Key: "service", Value: service},
		}
	}
	//filters = append(filters, filters)
	if ctx != "" {
		d := bson.E{Key: "ctx", Value: ctx}
		filters = append(filters, d)
	}
	withHost := bson.E{Key: "host", Value: host}
	filtersWithHost := append(filters, withHost)
	utils.TermDebugging(`filtersWithHost`, filtersWithHost)
	cur, err := collection.Find(context.Background(), filtersWithHost)
	if err != nil {
		return nil, err
	}
	// check if there are documents
	hasDocuments := cur.TryNext(context.Background())
	if !hasDocuments {
		// try without host
		withoutHost := bson.E{Key: "host", Value: bson.D{bson.E{Key: "$in", Value: []any{"", nil}}}}
		filtersWithoutHost := append(filters, withoutHost)
		cur, err = collection.Find(context.Background(), filtersWithoutHost)
		if err != nil {
			return nil, err
		}
	}
	var res []string
	for cur.Next(context.TODO()) {
		pol := Policy{}
		cur.Decode(&pol)
		pols = append(pols, pol)
		if pol.Action == "" {
			res = append(res, fmt.Sprintf("%s,%s,%s", pol.Type, pol.Subject, pol.Object))
		} else {
			res = append(res, fmt.Sprintf("%s,%s,%s,%s", pol.Type, pol.Subject, pol.Object, pol.Action))
		}
	}
	return res, nil
}
