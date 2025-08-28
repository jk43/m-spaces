package main

import (
	"context"
	"encoding/json"
	"fmt"
	"organization/models"
	"os"
	"reflect"
	"strconv"
	"strings"

	"net/http"

	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
)

func (srv *grpcServer) GetOrgID(ctx context.Context, in *pb.OrgRequest) (*pb.OrgIDResponse, error) {
	org := srv.app.DB.GetOrgWithHostAddr(in.Host)
	if org.Info["host"] == "" {
		return nil, fmt.Errorf("Organization " + in.Host + "  not found")
	}
	return &pb.OrgIDResponse{
		ID: org.ID,
	}, nil
}

func (srv *grpcServer) GetOrgSettings(ctx context.Context, in *pb.OrgSettingRequest) (*pb.OrgResponse, error) {
	var org *service.Organization
	var err error
	if in.Host != "" {
		org = srv.app.DB.GetOrgWithHostAddr(in.Host)
	}
	if in.OrgID != "" {
		org, err = srv.app.DB.GetOrgWithID(in.OrgID)
		if err != nil {
			return nil, err
		}
	}
	settings, err := json.Marshal(org.Settings)
	if err != nil {
		return nil, err
	}
	info, err := json.Marshal(org.Info)
	if err != nil {
		return nil, err
	}
	output := &pb.OrgResponse{
		ID:       org.ID,
		Settings: settings,
		Info:     info,
	}
	return output, nil
}

func (srv *grpcServer) SetOrgServiceSettings(ctx context.Context, in *pb.OrgServiceSettings) (*pb.OrgServiceSettings, error) {
	utils.TermDebugging(`configs`, nil)
	output := &pb.OrgServiceSettings{
		OrgID:   in.OrgID,
		Service: in.Service,
		Key:     in.Key,
	}
	utils.TermDebugging(`output`, output)
	res, err := srv.app.DB.InsertServiceSettings(in.OrgID, in.Service, in.Key, in.Settings)
	if err != nil {
		utils.TermDebugging(`err`, err)
	}
	utils.TermDebugging(`res`, res)
	return output, nil
}

func (srv *grpcServer) GetOrgServiceSettingsWithKey(ctx context.Context, in *pb.OrgServiceSettings) (*pb.OrgServiceSettings, error) {
	utils.TermDebugging(`GetConfigs`, in)
	return nil, nil
}

// Return accessible menus
// This function should be called after logging in, therefore it cannot be invoked from GetSettingsForEndUser.
func (app *application) GetItems(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeXHost, w, err, r)
		return
	}
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	perm, err := utils.NewCasbin(req.Host, utils.ServiceOrg, utils.CtxOrgMenu, "").ParsePermission(req.Role, utils.PermissionParserAsJoin)
	//perm, err := utils.GetPermissions(utils.ServiceOrg, req.Host, utils.CtxOrgMenu, utils.UserRole(req.Role), utils.GeneralPermissionParser)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeXHost, w, err, r)
		return
	}
	menus := make(map[string][]models.Item)
	for _, p := range perm {
		pol := strings.Split(p, ":")
		item := app.DB.GetItem(pol[1])
		menus[pol[2]] = append(menus[pol[2]], *item)
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: menus})
}

// Return settings for end user(RoleMember) to store in the Pinia store.
func (app *application) GetSettingsForEndUser(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, r)
		return
	}
	cache := molylibs.Redis[service.Organization]{
		DBNumber:   molylibs.RedisDBShortLife,
		Expiration: 10,
	}
	cache.SetKey(req.Host, "organization-service", "GetSettingsForEndUser")
	org := app.DB.GetOrgWithHostAddr(req.Host)
	forms, err := getForms(req.Role, req.Host, org)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedCasbin, w, err, req)
		return
	}
	// if req.Role != utils.RoleMember {
	// 	formsAsMember, err := getForms(utils.RoleMember, req.Host, org)
	// 	if err != nil {
	// 		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedCasbin, w, err, req)
	// 		return
	// 	}
	// 	//override forms as member
	// 	for k, v := range formsAsMember {
	// 		forms[k] = v
	// 	}
	// }
	output := service.GetDataWithFormInputs(org.Info, org.Settings.Forms[service.FormNameOrgInfo], true)
	if org.Info["cdnAddr"] == "" {
		output["cdnAddr"] = "https://" + os.Getenv("AWS_CLOUDFRONT_DOMAIN")
	}
	output["otp"] = org.Settings.Auth.OTP
	output["settings"] = utils.MapStringAny{}
	output["settings"].(utils.MapStringAny)["forms"] = forms
	perm, err := utils.NewCasbin(req.Host, utils.ServiceOrg, utils.CtxRBACRole, "").ParsePermission(req.Role, utils.GeneralPermissionParser)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedCasbin, w, err, req)
		return
	}
	roles := make([]utils.MapStringAny, 0)
	for _, p := range perm {
		o, _ := strconv.Atoi(p)
		r := utils.MapStringAny{
			"order": o,
			"role":  p,
		}
		roles = append(roles, r)
	}
	output["settings"].(utils.MapStringAny)["roles"] = roles

	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: output})
}

func (app *application) GetInfo(w http.ResponseWriter, r *http.Request) {
	utils.TermDebugging("GetInfo", nil)
	host, err := utils.GetXHost(r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeXHost, w, err, r)
		return
	}
	org := app.DB.GetOrgWithHostAddr(host)
	if org.Info["cdnAddr"] == "" {
		org.Info["cdnAddr"] = "https://" + os.Getenv("AWS_CLOUDFRONT_DOMAIN")
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: org.Info})
}

func (app *application) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[utils.FormBuilderFields](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, r)
		return
	}
	org := app.DB.GetOrgWithHostAddr(req.Host)

	getDocField := func(formName, key string) string {
		return `info.` + key
	}
	casbin := utils.NewCasbin(req.Host, utils.ServiceOrg, utils.CtxOrgFormEditView, utils.CtxOrgFormEditView)
	utils.TermDebugging(`req.Role`, req.Role)
	_, toUpdate, errors := service.ParseDataToUpdate(req.Role, casbin, service.FormNameOrgInfo, org.Info, req.Payload.Data, &org.Settings, getDocField)
	if errors != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: errors})
		return
	}
	err = app.DB.UpdateOrgWithHostAddr(req.Host, toUpdate)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedMongo, w, err, req, toUpdate)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func (app *application) UpdateFormInputOrder(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, r)
		return
	}
	_ = app.DB.UpdateFormOrder(req.OrgID, req.Payload.Data)
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func (app *application) DeleteFormInput(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[service.FormInput](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, r)
		return
	}
	_, err = app.DB.DeleteFormInput(req.OrgID, service.FormNameUserMetadata, req.Payload.Data.Key)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func Contains[T any](arr []T, item T) bool {
	for _, a := range arr {
		if reflect.DeepEqual(a, item) {
			return true
		}
	}
	return false
}

func (app *application) GetForm(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r, "form")
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, r)
	}
	form := req.Payload.Data["form"]
	casbin := utils.NewCasbin(req.Host, utils.ServiceOrg, utils.CtxOrgForm, "")
	perms, err := casbin.ParsePermission(req.Role, utils.GeneralPermissionParser)
	//perm, err := utils.NewFormPermissionPackage(req.Host, req.Role)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedCasbin, w, err, req)
	}
	if !Contains(perms, form[0]) {
		utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeUnauthorized, w, err, req)
		return
	}
	org, err := app.DB.GetFormWithFormName(req.OrgID, form[0])
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req, form)
		return
	}
	forms, err := getForms(req.Role, req.Host, org)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req, form)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: forms})
}

func (app *application) UpdateFormInput(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[service.FormInput](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, r)
	}
	_, err = app.DB.UpdateFormInput(req.OrgID, service.FormNameUserMetadata, req.Payload.Data)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func (app *application) CreateFormInput(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[service.FormInput](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, r)
	}
	utils.TermDebugging(`req`, req.Payload.Data)
	found, err := app.DB.GetFormInputWithKey(req.OrgID, service.FormNameUserMetadata, req.Payload.Data.Key)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req.OrgID, service.FormNameUserMetadata, req)
		return
	}
	if found.Key == req.Payload.Data.Key {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: []utils.ErrorDetails{{Field: "key", Error: "Key already exists"}}})
		return
	}
	_, err = app.DB.InsertFormInput(req.OrgID, service.FormNameUserMetadata, req.Payload.Data)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func getForms(role utils.UserRole, host string, data *service.Organization) (map[string][]service.FormInput, error) {
	casbin := utils.NewCasbin(host, utils.ServiceOrg, utils.CtxOrgFormEditView, utils.CtxOrgFormEditView)
	perm, err := utils.NewCasbin(host, utils.ServiceOrg, utils.CtxOrgForm, "").ParsePermission(role, utils.GeneralPermissionParser)
	//perm, err := utils.NewFormPermissionPackage(host, role)
	if err != nil {
		return nil, err
	}
	//flag editable based on role
	forms := make(map[string][]service.FormInput)
	for _, f := range perm {
		form := data.Settings.Forms[f]
		for _, i := range form {
			edit, _ := casbin.HasLink(role, i.Edit)
			if edit {
				i.Editable = true
			}
			// Drop the input if not viewable
			view, _ := casbin.HasLink(role, i.View)
			if !view {
				continue
			}
			forms[f] = append(forms[f], i)
		}
	}
	return forms, nil
}
