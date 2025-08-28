// Package models contains the data models for the application
package models

// Import necessary packages
import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateEmailRequest represents a request to update an email
type UpdateEmailRequest struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`                  // Unique ID
	UserID      string             `bson:"userId" json:"userId"`           // User ID
	Email       string             `bson:"email" json:"email"`             // New email
	ConfirmCode string             `bson:"confirmCode" json:"confirmCode"` // Confirmation code
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`     // Creation time
}

// InsertUpdateEamilRequest inserts a new UpdateEmailRequest into the database
func (repo *DBRepo) InsertUpdateEamilRequest(req UpdateEmailRequest) error {
	coll := repo.Mongo.Database(database).Collection("update_email_requests") // Get collection
	data, err := bson.Marshal(req)                                            // Convert request to BSON
	if err != nil {
		return err // Return error if any
	}
	_, err = coll.InsertOne(context.TODO(), data) // Insert data into collection
	if err != nil {
		return err // Return error if any
	}
	return nil // Return nil if no error
}

// GetUpdateEmailRequestByUserID retrieves an UpdateEmailRequest from the database by user ID
func (repo *DBRepo) GetUpdateEmailRequestWithUserID(userID string) (UpdateEmailRequest, error) {
	coll := repo.Mongo.Database(database).Collection("update_email_requests") // Get collection
	var req UpdateEmailRequest                                                // Initialize request variable

	opts := options.FindOne().SetSort(bson.D{{"createdAt", -1}})                     // Set find options
	err := coll.FindOne(context.TODO(), bson.M{"userId": userID}, opts).Decode(&req) // Find and decode request
	if err != nil {
		return req, err // Return request and error if any
	}
	return req, nil // Return request and nil if no error
}

// DeleteUpdateEmailRequestByID deletes an UpdateEmailRequest from the database by ID
func (repo *DBRepo) DeleteUpdateEmailRequestWithID(id primitive.ObjectID) error {
	coll := repo.Mongo.Database(database).Collection("update_email_requests") // Get collection
	_, err := coll.DeleteOne(context.TODO(), bson.M{"_id": id})               // Delete request by ID
	if err != nil {
		return err // Return error if any
	}
	return nil // Return nil if no error
}
