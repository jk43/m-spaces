package models

import (
	"os"

	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

var database = os.Getenv("MONGO_DATABASE")

type DBRepo struct {
	Mongo *mongo.Client
}

type DatabaseRepo interface {
	GetOrgWithHostAddr(host string) *service.Organization
	GetOrgWithID(id string) (*service.Organization, error)
	GetItem(id string) *Item
	UpdateOrgWithHostAddr(host string, update utils.MapStringAny) error
	UpdateFormOrder(id string, data utils.MapStringAny) error
	UpdateFormInput(id, form string, data service.FormInput) (*mongo.UpdateResult, error)
	GetFormInputWithKey(id, form, key string) (*service.FormInput, error)
	InsertFormInput(id, form string, data service.FormInput) (*mongo.UpdateResult, error)
	DeleteFormInput(id, form, key string) (*mongo.UpdateResult, error)
	GetFormWithFormName(id, form string) (*service.Organization, error)
	InsertServiceSettings(id, service, key string, data []*pb.OrgServiceSetting) (*mongo.UpdateResult, error)
}

type TestDBRepo struct{}
