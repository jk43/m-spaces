package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"user/models"

	"github.com/google/uuid"
	"github.com/moly-space/molylibs"
	"github.com/moly-space/molylibs/kafka"
	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (serv *grpcFileClientServer) GetFileRules(ctx context.Context, in *pb.FileRulesRequest) (*pb.FileRulesResponse, error) {
	if in.UserID == "" {
		return nil, fmt.Errorf("UserID is required")
	}
	if in.OrgID == "" {
		return nil, fmt.Errorf("OrgID is required")
	}
	res := &pb.FileRulesResponse{}
	res.S3Dir = in.OrgID + "/users" + "/" + in.UserID + "/" + in.ServiceCtx + "/"
	if in.ServiceCtx == "profile-image" {
		res.MaxFileSize = 10 << 20
		res.AllowedContentTypes = []string{"*"}
		res.AllowedFileTypes = []string{}
	}
	return res, nil
}

func (serv *grpcFileClientServer) SaveFileData(ctx context.Context, in *pb.FileSaveDataRequest) (*pb.FileSaveDataResponse, error) {
	profileImage := service.Metadata{
		"profileImage": in.FileInfos[0].S3Key,
	}
	err := serv.app.DB.UpdateMetadataWithID(in.UserID, profileImage, nil)
	if err != nil {
		return nil, err
	}
	return &pb.FileSaveDataResponse{}, nil
}

func (serv *grpcServer) ProcessOAuthUser(ctx context.Context, in *pb.OAuthUserRequest) (*pb.UserResponse, error) {
	org, err := service.GetOrg(nil, in.OrgID)
	if err != nil {
		return nil, err
	}
	user, err := serv.app.DB.GetUserWithEmailAndOrgID(in.Email, in.OrgID)
	if !org.Info["allowSelfMemberRegistration"].(bool) && user.Email == "" {
		return nil, fmt.Errorf("self member registration is not allowed")
	}
	if user.Email == "" {
		orgID, err := primitive.ObjectIDFromHex(in.OrgID)
		if err != nil {
			return nil, err
		}
		user.OrganizationID = orgID
		user.Email = in.Email
		user.Role = utils.RoleMember
		user.Verified = true
		user.CreatedAt = time.Now()
		user.RegisterMethod = service.RegisterMethod(in.RegisterMethod)
		err = buildUser(org, user)
		res, err := serv.app.DB.InsertUser(user)
		if err != nil {
			return nil, err
		}
		user.ID = res.InsertedID.(primitive.ObjectID)
	} else {
		data := make(utils.MapStringAny)
		data["firstName"] = in.FirstName
		data["lastName"] = in.LastName
		if user.ProfileImage == "" {
			data["profileImage"] = in.ProfileImage
		}
		serv.app.DB.UpdateUserWithIDAndOrgID(user.ID.Hex(), user.OrganizationID.Hex(), data)
	}
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{
		ID:        user.ID.Hex(),
		OrgID:     user.OrganizationID.Hex(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}, nil
}

func (serv *grpcServer) GetUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	var err error
	var user *service.User
	if in.UserID != "" && in.OrgID != "" {
		user, err = serv.app.DB.GetUserWithIDAndOrgID(in.UserID, in.OrgID)
	}
	if in.UserID != "" {
		user, err = serv.app.DB.GetUserWithID(in.UserID)
	}
	if in.Email != "" && in.OrgID != "" {
		user, err = serv.app.DB.GetUserWithEmailAndOrgID(in.Email, in.OrgID)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	res := pb.UserResponse{
		ID:           user.ID.Hex(),
		OrgID:        user.OrganizationID.Hex(),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		ProfileImage: user.ProfileImage,
		Role:         string(user.Role),
		Status:       string(user.Status),
	}
	// metadataSets := make(utils.Permissions, len(in.MetadataSets))
	// for i, p := range in.MetadataSets {
	// 	metadataSets[i].Act = p
	// }
	metadata, err := user.GetMetadata(in.MetadataSets, in.ShareableMetadataOnly)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	res.Metadata = bytes
	res.RegisterMethod = string(user.RegisterMethod)
	return &res, nil
}

func (serv *grpcServer) SetVerified(ctx context.Context, in *pb.EmailVerifedRequest) (*emptypb.Empty, error) {
	utils.TermDebugging(`in`, in)
	err := serv.app.DB.SetVerifiedWithID(in)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func buildUser(org *service.Org, user *service.User) error {
	if user.OrganizationID == primitive.NilObjectID {
		return fmt.Errorf("organizationID is required")
	}
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if user.FirstName == "" {
		return fmt.Errorf("firstName is required")
	}
	if user.LastName == "" {
		return fmt.Errorf("lastName is required")
	}
	user.Metadata = make(map[string]service.Metadata)
	metadata := make(service.Metadata)
	for _, s := range org.Forms[service.FormNameUserMetadata] {
		metadata[s.Key] = s.DefaultValue
	}
	user.Metadata[service.FormNameUserMetadata] = metadata
	user.Status = service.UserStatus(org.Info["initialStatusOfSelfRegisteredMember"].(string))
	user.Store = make(map[string][]service.StoreData)
	return nil
}

func (app *application) RegisterUser(w http.ResponseWriter, r *http.Request) {
	payload := utils.Payload[models.UserRequest]{
		Data: models.UserRequest{},
	}
	utils.ReadJSON(w, r, &payload)
	user := payload.Data

	orgID, err := service.GetOrgID(r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeXHost, w, err, r)
		return
	}
	//userExist := false
	_, err = app.DB.GetUserWithEmailAndOrgID(user.Email, orgID.Hex())
	if err != mongo.ErrNoDocuments {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeNone, w, err, r)
		return
	}
	v := molylibs.NewValidator(nil)

	org, err := service.GetOrg(r, "")

	if org.Settings.Auth.OTP {
		// if OTP is enabled, password is not required and set the password to a random string
		password := uuid.New().String()
		user.Password = password
		user.ConfirmPassword = password
	}

	err, validatorErrors := v.Validate(&user)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: validatorErrors})
		return
	}

	hashedPassword, salt, _ := user.HashPassword()
	user.OrganizationID = orgID
	user.CreatedAt = time.Now()
	user.Role = utils.RoleMember
	user.Status = service.UserStatus(org.Info["initialStatusOfSelfRegisteredMember"].(string))
	utils.TermDebugging(`org.Info`, org.Info)
	user.Store = make(map[string][]service.StoreData)
	user.Verified = false
	//user.SelfRegistered = true
	user.RegisterMethod = service.RegisterMethodSelf
	var userObjectID primitive.ObjectID
	user.Metadata = make(map[string]service.Metadata)
	metadata := make(service.Metadata)

	for _, s := range org.Forms[service.FormNameUserMetadata] {
		metadata[s.Key] = s.DefaultValue
	}
	user.Metadata[service.FormNameUserMetadata] = metadata

	res, err := app.DB.InsertUser(&user.User)

	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, r)
		return
	}
	userObjectID = res.InsertedID.(primitive.ObjectID)

	userID := userObjectID.Hex()
	result := utils.Response{Result: utils.SUCCESS, Data: map[string]string{"id": userID}}

	creds := pb.CredentialsRequest{
		Email:          user.Email,
		HashedPassword: hashedPassword,
		Salt:           salt,
		UserID:         userID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		OrgID:          orgID.Hex(),
		Status:         string(user.Status),
	}
	_, err = app.GRPCClient.SaveCredentials(&creds)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err, r)
		return
	}
	utils.WriteJSON(w, http.StatusOK, result)
	return
}

func (app *application) GetUser(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, r)
	}
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeXHost, w, err, req)
		return
	}
	user, err := app.DB.GetUserWithIDAndOrgID(req.UserID, req.OrgID)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeFailedMongo, w, err, req)
		return
	}
	perm, err := utils.NewCasbin(req.Host, utils.ServiceOrg, utils.CtxOrgForm, "").ParsePermission(req.Role, utils.GeneralPermissionParser)
	//perm, err := utils.NewFormPermissionPackage(host, role)
	//perm, err := utils.NewFormPermissionPackage(req.Host, req.Role)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedCasbin, w, err, req)
		return
	}
	//sets, filters, err := utils.GetPermissions("organization", host, utils.CtxOrgForm, role, utils.GeneralPermissionParser)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedCasbin, w, err, req)
		return
	}
	metadata, err := user.GetMetadata(perm, false)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeGetOrgSettings, w, err, req)
		return
	}
	userOutput := struct {
		Email     string                      `json:"email"`
		FirstName string                      `json:"firstName"`
		LastName  string                      `json:"lastName"`
		Metadata  map[string]service.Metadata `json:"metadata"`
	}{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Metadata:  metadata,
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: userOutput})
}

func (app *application) UpdateUserAccount(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	//orgIDObjectID, _ := service.GetOrgID(r)
	org, err := service.GetOrg(r, "")
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeGetOrgSettings, w, err, r)
		return
	}
	userID := claims.Subject
	orgID := org.ID
	payload := utils.Payload[map[string]utils.FormBuilderField]{
		Data: make(map[string]utils.FormBuilderField),
	}
	utils.ReadJSON(w, r, &payload)
	u, err := app.DB.GetUserWithID(userID)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, userID)
		return
	}
	// error if the user is not self/by admin registered
	if u.RegisterMethod != service.RegisterMethodSelf && u.RegisterMethod != service.RegisterMethodAdmin {
		httpError := utils.NewHttpError("", "", "Unable to update data for oauth user.", utils.ErrorCodeOAuthUserUpdate, http.StatusBadRequest, err)
		httpError.Write(w)
		return
	}
	utils.TermDebugging(`u`, u)
	user := service.User{
		Email:     payload.Data["email"].Value.(string),
		FirstName: payload.Data["firstName"].Value.(string),
		LastName:  payload.Data["lastName"].Value.(string),
	}
	updateRequest := pb.UpdateNameRequest{
		UserID:    userID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	_, err = app.GRPCClient.UpdateName(&updateRequest)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err, updateRequest)
		return
	}
	err = app.DB.UpdateAccountWithIDAndOrgID(userID, orgID, user)
	if err != nil {
		err = app.DB.UpdateAccountWithIDAndOrgID(userID, orgID, user)
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeUserUpdateFailed, w, err, userID, orgID, user)
		return
	}
	if user.Email != claims.Email {
		req := models.UpdateEmailRequest{
			UserID:      userID,
			Email:       user.Email,
			ConfirmCode: utils.GetVerificationCode(),
			CreatedAt:   time.Now(),
		}
		found, _ := app.DB.GetUserWithEmailAndOrgID(user.Email, orgID)
		if found.Email != "" {
			utils.TermDebugging(`found in error`, found)
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeEmailBeingUsed, w, err, user.Email, orgID)
			return
		}
		err = app.DB.InsertUpdateEamilRequest(req)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req)
		}
		//send verification email
		kafka.SendEmail(&kafka.SESTemplate{
			Template: "token-verification",
			TemplateData: map[string]string{
				"name":    claims.FirstName + " " + claims.LastName,
				"orgName": org.Info["name"].(string),
				"url":     org.Info["host"].(string),
				"code":    req.ConfirmCode},
			ToAddresses: []string{claims.Email},
			Source:      org.Info["email"].(string),
		})
		vi := utils.VerificationInstruction{
			URL:     "/user/email",
			KeyName: "code",
			Method:  "POST",
			Message: "A verification code has been sent to your email. Please enter the code to complete the process.",
		}
		utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.VERIFICATION_REQUIRED, Data: vi})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func (app *application) UpdateEmailWithVerificationCode(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	userID := claims.Subject

	type verifycationCode struct {
		Code string `bson:"code" json:"code"`
	}
	payload := utils.Payload[verifycationCode]{
		Data: verifycationCode{},
	}
	utils.ReadJSON(w, r, &payload)
	res, err := app.DB.GetUpdateEmailRequestWithUserID(userID)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, userID)
		return
	}
	if res.ConfirmCode != payload.Data.Code {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeNone, w, err, r)
		return
	}
	updateRequest := pb.UpdateEmailRequest{
		UserID: userID,
		Email:  res.Email,
	}
	_, err = app.GRPCClient.UpdateEmail(&updateRequest)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err, updateRequest)
		return
	}
	err = app.DB.UpdateEmailWithID(userID, res.Email)
	if err != nil {
		_, id := utils.NewErrorDetails("", "", "", utils.ErrorCodeFailedMongo)
		app.Logger.Error().Err(err).Str("eid", string(id)).Send()
	}
	err = app.DB.DeleteUpdateEmailRequestWithID(res.ID)
	if err != nil {
		_, id := utils.NewErrorDetails("", "", "", utils.ErrorCodeFailedMongo)
		app.Logger.Error().Err(err).Str("eid", string(id)).Send()
	}
	data := map[string]string{"title": "Successfully changed email", "message": "You have successfully changed your email. Please close this dialog and you will be able to log in with your new email from the next login."}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: data})
}

func (app *application) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	userID := claims.Subject

	type passwords struct {
		CurrentPassword utils.FormBuilderFieldRequired `json:"currentPassword" validate:"required"`
		Password        utils.FormBuilderFieldRequired `json:"password" validate:"required"`
		ConfirmPassword utils.FormBuilderFieldRequired `json:"confirmPassword" validate:"required"`
	}

	payload := utils.Payload[passwords]{
		Data: passwords{},
	}

	v := molylibs.NewValidator(nil)

	utils.ReadJSON(w, r, &payload)
	data := payload.GetData()

	err, validatorErrors := v.Validate(&payload.Data)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: validatorErrors})
		return
	}

	in := &pb.VerifyPasswordRequest{
		UserID:   userID,
		Password: data.CurrentPassword.Value.(string),
	}

	_, err = app.GRPCClient.VerifyPassword(in)
	if err != nil {
		ed, _ := utils.NewErrorDetails("", "", "Cannot verify password", utils.ErrorCodePasswordNotMatch)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Response{Result: utils.ERROR, Data: ed})
		return
	}

	if data.Password.Value.(string) != data.ConfirmPassword.Value.(string) {
		ed, _ := utils.NewErrorDetails("", "", "Passwords do not match", utils.ErrorCodePasswordNotMatch)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: ed})
		return
	}

	hashedPassword, salt, err := utils.HashPassword(data.Password.Value.(string))
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedHashing, w, err)
		return
	}

	req := models.UpdatePasswordRequest{
		UserID:      userID,
		Password:    hashedPassword,
		Salt:        salt,
		ConfirmCode: utils.GetVerificationCode(),
		CreatedAt:   time.Now(),
	}

	err = app.DB.InsertUpdatePasswordRequest(req)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedMongo, w, err)
		return
	}

	org, err := service.GetOrg(r, "")
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeGetOrgSettings, w, err, r)
		return
	}

	//send verification email
	kafka.SendEmail(&kafka.SESTemplate{
		Template: "token-verification",
		TemplateData: map[string]string{
			"name":    claims.FirstName + " " + claims.LastName,
			"orgName": org.Info["name"].(string),
			"url":     org.Info["host"].(string),
			"code":    req.ConfirmCode},
		ToAddresses: []string{claims.Email},
		Source:      org.Info["email"].(string),
	})
	vi := utils.VerificationInstruction{
		URL:     "/user/password",
		KeyName: "code",
		Method:  "POST",
		Message: "A verification code has been sent to your email. Please enter the code to complete the process.",
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.VERIFICATION_REQUIRED, Data: vi})
}

func (app *application) UpdatePasswordWithVerificationCode(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	userID := claims.Subject

	type verifycationCode struct {
		Code string `bson:"code" json:"code"`
	}
	payload := utils.Payload[verifycationCode]{
		Data: verifycationCode{},
	}
	utils.ReadJSON(w, r, &payload)
	res, err := app.DB.GetUpdatePasswordRequestWithUserID(userID)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, userID)
		return
	}
	if res.ConfirmCode != payload.Data.Code {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeNone, w, err, userID)
		return
	}
	in := &pb.PasswordUpdateRequest{
		UserID:   res.UserID,
		Password: res.Password,
		Salt:     res.Salt,
	}
	_, err = app.GRPCClient.UpdatePassword(in)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedGRPC, w, err, in)
		return
	}
	err = app.DB.DeleteUpdatePasswordRequestWithID(res.ID)
	if err != nil {
		_, id := utils.NewErrorDetails("", "", "", utils.ErrorCodeFailedMongo)
		app.Logger.Error().Err(err).Str("eid", string(id)).Send()
	}
	data := map[string]string{"title": "Successfully changed password", "message": "You have successfully changed your password. Please close this dialog and you will be able to log in with your new password from the next login."}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: data})
}

func (app *application) UpdateUserMetadata(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[utils.FormBuilderFields](w, r)
	user, err := app.DB.GetUserWithID(req.UserID)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeUserNotFound, w, err, req.UserID)
		return
	}
	metadata, err := user.GetMetadata([]string{service.FormNameUserMetadata}, false)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeGetOrgSettings, w, err, user)
		return
	}
	org, _ := service.GetOrg(r, "")
	getDocField := func(formName, key string) string {
		return `metadata.` + formName + `.` + key
	}
	current := metadata[service.FormNameUserMetadata]
	casbin := utils.NewCasbin(req.Host, utils.ServiceOrg, utils.CtxOrgFormEditView, utils.CtxOrgFormEditView)
	toDelete, toUpdate, errors := service.ParseDataToUpdate(req.Role, casbin, service.FormNameUserMetadata, current, req.Payload.Data, &org.Settings, getDocField)
	if errors != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: errors})
		return
	}
	err = app.DB.UpdateMetadataWithID(req.UserID, toUpdate, toDelete)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedMongo, w, err, req.UserID, toUpdate, toDelete)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
	return
}

func (app *application) GetOrgUsers(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r)
	pn := utils.Pagination{}
	utils.QueryExtractor(&req.Payload.Data, &pn)

	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
	}
	users, total, err := app.DB.GetUsersWithOrgID(req.OrgID, req.Payload.Data, &pn)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedMongo, w, err, req)
		return
	}
	perm, err := utils.NewCasbin(req.Host, utils.ServiceOrg, utils.CtxOrgForm, "").ParsePermission(req.Role, utils.GeneralPermissionParser)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedCasbin, w, err)
		return
	}
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedCasbin, w, err)
		return
	}
	for i, u := range users {
		metadata, err := u.GetMetadata(perm, false)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeGetOrgSettings, w, err, u)
			return
		}
		users[i].Metadata = metadata
	}
	res := utils.MapStringAny{
		"total": total,
		"data":  users,
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: &res})
}

func (app *application) CreateOrgUser(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[models.UserRequest](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	//set random password
	randPW := utils.GetRandStrings(16)
	req.Payload.Data.ConfirmPassword = randPW
	req.Payload.Data.Password = randPW
	objOrgId, err := primitive.ObjectIDFromHex(req.OrgID)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req)
		return
	}
	user := req.Payload.Data.User
	user.OrganizationID = objOrgId
	user.Status = service.UserStatusWaiting
	hashedPassword, salt, _ := req.Payload.Data.HashPassword()
	//check if user already exists
	_, err = app.DB.GetUserWithEmailAndOrgID(req.Payload.Data.Email, req.OrgID)
	if err != mongo.ErrNoDocuments {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: []utils.ErrorDetails{{Field: "email", Error: "User exists"}}})
		return
	}
	user.Store = make(map[string][]service.StoreData)
	user.Verified = false
	//user.SelfRegistered = false
	user.RegisterMethod = service.RegisterMethodAdmin
	user.CreatedAt = time.Now()
	utils.TermDebugging(`user`, user)
	res, err := app.DB.InsertUser(&user)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req)
		return
	}
	utils.TermDebugging(`req.Payload.Data`, req.Payload.Data)
	creds := pb.CredentialsRequest{
		Email:          req.Payload.Data.Email,
		HashedPassword: hashedPassword,
		Salt:           salt,
		UserID:         res.InsertedID.(primitive.ObjectID).Hex(),
		FirstName:      req.Payload.Data.FirstName,
		LastName:       req.Payload.Data.LastName,
		OrgID:          req.OrgID,
		Invite:         true,
		SenderName:     req.Claims.FirstName + " " + req.Claims.LastName,
		Status:         string(service.UserStatusWaiting),
	}
	_, err = app.GRPCClient.SaveCredentials(&creds)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err, creds)
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func (app *application) UpdateOrgUser(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[utils.FormBuilderFields](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	fieldsToRemove := []string{"id", "confirmPassword", "password"}
	data := utils.MapStringAny{}
	var userID string
	for k, v := range req.Payload.Data {
		inRemove := false
		if k == "id" {
			userID = v.Value.(string)
		}
		for _, f := range fieldsToRemove {
			if k == f {
				inRemove = true
				break
			}
		}
		if !inRemove {
			data[k] = v.Value
		}
	}
	user, err := app.DB.GetUserWithID(userID)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedMongo, w, err, req)
		return
	}
	if user.RegisterMethod == service.RegisterMethodSelf && user.Status == service.UserStatusInactive && data["status"] == string(service.UserStatusActive) {
		//send verification email
		org, _ := service.GetOrg(r, "")
		kafka.SendEmail(&kafka.SESTemplate{
			Template: "account-activation",
			TemplateData: map[string]string{
				"name":        user.FirstName + " " + user.LastName,
				"orgName":     org.Info["name"].(string),
				"url":         org.Info["host"].(string),
				"sender_name": req.Claims.FirstName + " " + req.Claims.LastName,
			},
			ToAddresses: []string{user.Email},
			Source:      org.Info["email"].(string),
		})
	}
	res, err := app.DB.UpdateUserWithIDAndOrgID(userID, req.OrgID, data)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedMongo, w, err, req)
		return
	}
	if res.MatchedCount != 1 {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeMongoZeroMatchedCount, w, err, req)
		return
	}

	var hashedPassword, salt string
	// process password at the end for security
	if req.Payload.Data["password"].Value != nil {
		if req.Payload.Data["confirmPassword"].Value == nil {
			ed, _ := utils.NewErrorDetails("", "confirmPassword", "Passwords do not match", utils.ErrorCodePasswordNotMatch)
			utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: []utils.ErrorDetails{ed}})
			return
		}
		password := req.Payload.Data["password"].Value.(string)
		confirmpassword := req.Payload.Data["confirmPassword"].Value.(string)
		if password != confirmpassword {
			ed, _ := utils.NewErrorDetails("", "", "Passwords do not match", utils.ErrorCodePasswordNotMatch)
			utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: ed})
			return
		}
		hashedPassword, salt, err = utils.HashPassword(password)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedHashing, w, err)
			return
		}
	}
	update := pb.CredentialsRequest{
		UserID:         userID,
		FirstName:      req.Payload.Data["firstName"].Value.(string),
		LastName:       req.Payload.Data["lastName"].Value.(string),
		Email:          req.Payload.Data["email"].Value.(string),
		HashedPassword: hashedPassword,
		Salt:           salt,
	}
	app.GRPCClient.UpdateCredentials(&update)
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func (app *application) DeleteOrgUser(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	id := req.Payload.Data["id"].(string)
	err = app.DB.DeleteUserWithID(id)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req)
		return
	}
	dReq := pb.DeleteCredentialsRequest{
		UserID: id,
	}
	_, err = app.GRPCClient.DeleteCredentials(&dReq)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err, dReq)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func (app *application) GetStore(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r, "ctx")
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
	}
	res, err := app.DB.GetStore(req.UserID, req.Payload.Data["ctx"][0], "")
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req)
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: res})
}

func (app *application) DeleteStore(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	_, err = app.DB.DeleteStore(req.UserID, req.Payload.Data["ctx"].(string), req.Payload.Data["key"].(string))
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req)
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func (app *application) CreateStore(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	ctx := req.Payload.Data["ctx"].(string)
	name := req.Payload.Data["name"].(string)
	key := req.Payload.Data["key"].(string)
	value := req.Payload.Data["value"].(utils.MapStringAny)
	data := service.StoreData{
		Key:   key,
		Label: name,
		Value: value,
	}
	_, err = app.DB.GetStore(req.UserID, ctx, key)
	if err == nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: []utils.ErrorDetails{{Field: "name", Error: "Name already exists"}}})
		return
	}
	_, err = app.DB.AddStore(req.UserID, ctx, data)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMongo, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func (app *application) UpdateStore(w http.ResponseWriter, r *http.Request) {

}
