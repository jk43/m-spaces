package models

import (
	"errors"
	"os"

	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var database = os.Getenv("MONGO_DATABASE")

type DBRepo struct {
	Mongo *mongo.Client
}

type DatabaseRepo interface {
	InsertUser(u *service.User) (*mongo.InsertOneResult, error)
	GetUserWithID(id string) (*service.User, error)
	GetUserWithIDAndOrgID(id, orgID string) (*service.User, error)
	GetUserWithEmailAndOrgID(email, orgObjectID string) (*service.User, error)
	UpdateAccountWithIDAndOrgID(id, orgID string, user service.User) error
	UpdateMetadataWithID(id string, update, delete service.Metadata) error
	InsertUpdatePasswordRequest(req UpdatePasswordRequest) error
	GetUpdatePasswordRequestWithUserID(userID string) (UpdatePasswordRequest, error)
	DeleteUpdatePasswordRequestWithID(id primitive.ObjectID) error
	InsertUpdateEamilRequest(req UpdateEmailRequest) error
	GetUpdateEmailRequestWithUserID(userID string) (UpdateEmailRequest, error)
	DeleteUpdateEmailRequestWithID(id primitive.ObjectID) error
	UpdateEmailWithID(id string, email string) error
	GetUsersWithOrgID(orgID string, search utils.MapStringSlice[[]string], pn *utils.Pagination) ([]service.User, int64, error)
	UpdateUserWithIDAndOrgID(id, orgID string, data utils.MapStringAny) (*mongo.UpdateResult, error)
	AddStore(id, ctx string, data service.StoreData) (*mongo.UpdateResult, error)
	GetStore(id, ctx, name string) ([]service.StoreData, error)
	DeleteStore(id, ctx, key string) (*mongo.UpdateResult, error)
	DeleteUserWithID(id string) error
	SetVerifiedWithID(in *pb.EmailVerifedRequest) error
}

type TestDBRepo struct{}

func (repo *TestDBRepo) GetUserWithEmail(email string) service.User {
	user := service.User{}
	if email == "exist@test.com" {
		user.Email = "xxxx@xxxx.com"
		return user
	}
	return user
}

func (repo *TestDBRepo) InsertUser(u service.User) (*mongo.InsertOneResult, error) {
	if u.Email == "mongoerror@test.com" {
		return nil, errors.New("Error")
	}
	objectID, _ := primitive.ObjectIDFromHex("63a5e4739837ae51dc1d574c")
	ir := mongo.InsertOneResult{
		InsertedID: objectID,
	}
	return &ir, nil
}

func (repo *TestDBRepo) GetUerWithObjectID(oid string) service.User {
	return service.User{}
}
