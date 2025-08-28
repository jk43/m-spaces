package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FormName = string

const (
	FormNameUserMetadata FormName = "userMetadata"
	FormNameOrgInfo      FormName = "adminOrganizationInformation"
)

// type InfoAttr struct {
// 	Name      string            `bson:"name" json:"name"`
// 	Value     any               `bson:"value" json:"value"`
// 	FormInput InfoAttrFormInput `bson:"formInput" json:"formInput"`
// }

// type InfoAttrFormInput struct {
// 	Type         string             `bson:"type" json:"type"`
// 	Description  string             `bson:"description" json:"description"`
// 	Rules        any                `bson:"rules" json:"rules"`
// 	Options      utils.MapStringAny `bson:"options" json:"options"`
// 	DefaultValue any                `bson:"defaultValue" json:"defaultValue"`
// 	Editable     bool               `bson:"editable" json:"editable"`   // If the value is true, the filed will be appears on the form
// 	Shareable    bool               `bson:"shareable" json:"shareable"` // If the value is true, value will be share with front-end(vue store/JWT token)
// }

type Organization struct {
	ID       string             `bson:"_id"`
	Info     utils.MapStringAny `bson:"info" json:"info"`
	Settings Settings           `bson:"settings" json:"settings"`
}
type Settings struct {
	Forms map[FormName][]FormInput `bson:"forms" json:"forms"`
	Auth  AuthSetting              `bson:"auth" json:"auth"`
}

// FormInput struct represents the structure of a form input.
type FormInput struct {
	Key          string         `bson:"key" json:"key"`                   // The unique identifier for the form input
	Name         string         `bson:"name" json:"name"`                 // The name of the form input
	Type         string         `bson:"type" json:"type"`                 // The type of the form input (e.g., text, number, date)
	Description  string         `bson:"description" json:"description"`   // A description of the form input
	Rules        []string       `bson:"rules" json:"rules"`               // Validation rules for the form input
	Options      Options        `bson:"options" json:"options"`           // Options for the form input (used for select inputs)
	DefaultValue any            `bson:"defaultValue" json:"defaultValue"` // The default value for the form input
	Shareable    bool           `bson:"shareable" json:"shareable"`       // If true, the value will be shared with the front-end (vue store/JWT token)
	Order        int            `bson:"order" json:"order"`               // The order in which the form input appears
	Edit         utils.UserRole `bson:"edit" json:"edit"`                 // The user role that can edit this form input
	View         utils.UserRole `bson:"view" json:"view"`                 // The user role that can view this form input
	Editable     bool           `bson:"-" json:"editable"`                // If false, the front-end will just output the value
	OnTable      bool           `bson:"onTable" json:"onTable"`           // If false, the front-end will not show this field on the table
	Slug         string         `bson:"slug" json:"slug"`                 // The slug of the form input
}

// Options struct represents the options for a select form input.
type Options struct {
	LableValue []OptionLabelValue `bson:"labelValue" json:"labelValue"` // The label-value pairs for the options
	Options    []string           `bson:"options" json:"options"`       // The options for the select input
}

// OptionLabelValue struct represents a label-value pair for a select form input option.
type OptionLabelValue struct {
	Value any    `bson:"value" json:"value"` // The value of the option
	Label string `bson:"label" json:"label"` // The label of the option
}

// AuthSetting struct represents the authentication settings.
type AuthSetting struct {
	Type      string                 `bson:"type" json:"type"`   // The type of authentication
	Mfa       bool                   `bson:"mfa" json:"mfa"`     // If true, multi-factor authentication is enabled
	OTP       bool                   `bson:"otp" json:"otp"`     // If true, one-time password is enabled
	MfaSource []string               `bson:"mfaSource" json:"-"` // The sources for multi-factor authentication
	OAuth     map[string]OAuthConfig `bson:"oauth" json:"oauth"` // The OAuth settings
}

type OAuthConfig struct {
	ID       string            `bson:"id" json:"id"`             // The ID of the OAuth provider
	Secret   string            `bson:"secret" json:"secret"`     // The secret of the OAuth provider
	Redirect string            `bson:"redirect" json:"redirect"` // The redirect URL of the OAuth provider
	Scopes   []string          `bson:"scopes" json:"scopes"`     // The scopes of the OAuth provider
	Options  map[string]string `bson:"options" json:"options"`   // The options of the OAuth provider
}

// The struct return over RPC
// Org struct represents an organization.
type Org struct {
	ID       string
	Settings                    // The settings for the organization
	Info     utils.MapStringAny `bson:"info" json:"info"` // Additional information about the organization
}

// GetOrg function retrieves an organization's settings.
func GetOrg(r *http.Request, orgID string) (*Org, error) {
	var xHost string
	var err error
	if r != nil {
		xHost, err = utils.GetXHost(r)
		if err != nil {
			return nil, err
		}
	}
	cache := molylibs.Redis[Org]{
		DBNumber:   molylibs.RedisDBShortLife,
		Expiration: 10,
	}
	cache.SetKey(xHost, "molylib", "serivce", "GetOrgSettings")
	//cached, err := cache.Get(context.Background())
	if err == nil {
		fmt.Println("Cache is disabled")
		//return cached, nil
	}
	conn, err := grpc.Dial(os.Getenv("ORGANIZATION_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewOrganizationServiceClient(conn)
	in := pb.OrgSettingRequest{
		Host:  xHost,
		OrgID: orgID,
	}
	res, err := c.GetOrgSettings(context.Background(), &in)
	if err != nil {
		return nil, err
	}
	org := Org{}
	org.ID = res.ID
	err = json.Unmarshal(res.Settings, &org.Settings)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res.Info, &org.Info)
	if err != nil {
		return nil, err
	}

	err = cache.Set(context.Background(), &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// GetOrgID function retrieves an organization's ID.
func GetOrgID(r *http.Request) (primitive.ObjectID, error) {
	xHost, err := utils.GetXHost(r)
	if err != nil {
		return primitive.NilObjectID, err
	}
	cache := molylibs.Redis[primitive.ObjectID]{
		DBNumber:   molylibs.RedisDBLongLife,
		Expiration: 60 * 24,
	}
	cache.SetKey(xHost, "molylib", "serivce", "GetOrgID")
	cached, err := cache.Get(context.Background())
	if err == nil {
		return *cached, nil
	}
	conn, err := grpc.Dial(os.Getenv("ORGANIZATION_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return primitive.NilObjectID, err
	}
	c := pb.NewOrganizationServiceClient(conn)
	in := pb.OrgRequest{
		Host: xHost,
	}
	res, err := c.GetOrgID(context.Background(), &in)
	if err != nil {
		return primitive.NilObjectID, err
	}
	objID, err := primitive.ObjectIDFromHex(res.ID)
	cache.Set(context.Background(), &objID)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return objID, nil
}

// Service Settings

// Service Settings

type ServiceSettingValueType string

const (
	ServiceSettingTypeString ServiceSettingValueType = "string"
	ServiceSettingTypeInt    ServiceSettingValueType = "int"
	ServiceSettingTypeFloat  ServiceSettingValueType = "float"
	ServiceSettingTypeBool   ServiceSettingValueType = "boolean"
)

type FormElemOptions struct {
	LabelValue []OptionLabelValue `json:"labelValue"` // The label-value pairs for the options
	Options    []string           `json:"options"`    // The options for the select input
}

type FormElemLabelValue struct {
	Label string `json:"label"` // The label of the option
	Value any    `json:"value"` // The value of the option
}

type FormElems map[string]*FormElem
type FormElem struct {
	Key          string                  `json:"key"`
	Value        string                  `json:"value"`
	ValueType    ServiceSettingValueType `json:"-"`
	Name         string                  `json:"name"`
	Description  string                  `json:"description"`
	Type         string                  `json:"type"`
	Editable     bool                    `json:"editable"`
	Rules        []string                `json:"rules"`
	Options      FormElemOptions         `json:"options"`
	Order        int                     `json:"order"`
	DropFromBase bool                    `json:"-"` // If true, the form element will be dropped from the base form
}

//type ServiceSettings map[string]*FormElem

func (c *FormElem) IntValue() (int, error) {
	if c.ValueType != ServiceSettingTypeInt {
		return 0, errors.New("value type is not int")
	}
	return strconv.Atoi(c.Value)
}

func (c *FormElem) BoolValue() (bool, error) {
	if c.ValueType != ServiceSettingTypeBool {
		return false, errors.New("value type is not bool")
	}
	return strconv.ParseBool(c.Value)
}

func (c *FormElem) FloatValue() (float64, error) {
	if c.ValueType != ServiceSettingTypeFloat {
		return 0, errors.New("value type is not float")
	}
	return strconv.ParseFloat(c.Value, 64)
}

func (c *FormElem) StringValue() string {
	return c.Value
}

func (c *FormElem) SetValue(value string) error {
	switch c.ValueType {
	case ServiceSettingTypeInt:
		if _, err := strconv.Atoi(value); err != nil {
			return errors.New("value is not int")
		}
	case ServiceSettingTypeFloat:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return errors.New("value is not float")
		}
	case ServiceSettingTypeBool:
		if _, err := strconv.ParseBool(value); err != nil {
			return errors.New("value is not bool")
		}
	}
	c.Value = value
	return nil
}

type FormElemsFunc = func(FormElems)

func NewForms(fn func(map[utils.Host]FormElemsFunc)) *Forms {
	f := &Forms{
		forms: make(map[utils.Host]FormElems),
	}
	fMap := make(map[utils.Host]func(FormElems))
	fn(fMap)
	for k, v := range fMap {
		f.forms[k] = make(FormElems)
		v(f.forms[k])
	}
	return f
}

type FormFuncMap func(FormElems) FormElems
type HostFormMap map[utils.Host]FormElems
type Forms struct {
	FormFuncs func(map[string]func(FormElems) FormElems) func(FormElems) FormElems
	forms     HostFormMap
}

func (f *Forms) GetForm(host utils.Host) FormElems {
	return f.forms[host]
}

type OrgServiceSettings struct {
	OrgID string
	Ctx   string
	Slug  string // key of settings
	//Dependency map[string]any // The dependency of the service settings
	FormElems FormElems
	//settings          ServiceSettings
}

func (s *OrgServiceSettings) Save(data utils.MapStringAny) []utils.ErrorDetails {
	// settings.SetServiceSettings()
	errs := []utils.ErrorDetails{}
	for k, v := range s.FormElems {
		if _, ok := data[k]; ok {
			err := v.SetValue(data[k].(string))
			if err != nil {
				errDetail, _ := utils.NewErrorDetails("", k, err.Error(), utils.ErrorCodeNone)
				errs = append(errs, errDetail)
			}
			v.Key = k
		}
	}
	if len(errs) > 0 {
		return errs
	}
	conn, err := grpc.Dial(os.Getenv("ORGANIZATION_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errDetail, _ := utils.NewErrorDetails("", "", string(utils.UnableToProcessRequest), utils.ErrorCodeFailedGRPC)
		errs = append(errs, errDetail)
		return errs
	}
	defer conn.Close()
	c := pb.NewOrganizationServiceClient(conn)
	in := pb.OrgServiceSettings{
		OrgID:   s.OrgID,
		Service: s.Ctx,
		Key:     s.Slug,
	}
	for _, config := range s.FormElems {
		in.Settings = append(in.Settings, &pb.OrgServiceSetting{
			Key:       config.Key,
			Value:     config.Value,
			ValueType: string(config.ValueType),
		})
	}
	_, err = c.SetOrgServiceSettings(context.Background(), &in)
	if err != nil {
		errDetail, _ := utils.NewErrorDetails("", "", string(utils.UnableToProcessRequest), utils.ErrorCodeFailedGRPC)
		errs = append(errs, errDetail)
		return errs
	}
	return nil
}

func UpsertMysqlServiceSettingWithPayload(settings *OrgServiceSettings, data utils.MapStringAny, f func(FormElem, ...any) error, p ...any) error {
	return nil
}
