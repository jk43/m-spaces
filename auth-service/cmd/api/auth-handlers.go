package main

import (
	"auth/models"
	svc "auth/service"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/moly-space/molylibs/kafka"
	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (ser *grpcServer) SaveCredentials(ctx context.Context, in *pb.CredentialsRequest) (*pb.CredentialsResponse, error) {
	if in.Email == "" {
		return nil, errors.New("Email cannot be empty")
	}
	if in.HashedPassword == "" {
		return nil, errors.New("HashedPassword cannot be empty")
	}
	if in.UserID == "" {
		return nil, errors.New("UserID cannot be empty")
	}
	if in.OrgID == "" {
		return nil, errors.New("OrgId cannot be empty")
	}
	token := uuid.New().String()
	vToken := models.VerificationTokens{
		Token: sql.NullString{String: token, Valid: true},
	}
	cred := models.Credentials{
		UserID:             in.UserID,
		Email:              in.Email,
		Password:           in.HashedPassword,
		Salt:               in.Salt,
		OrganizationID:     in.OrgID,
		FirstName:          in.FirstName,
		LastName:           in.LastName,
		VerificationTokens: vToken,
		Status:             service.UserStatus(in.Status),
	}
	if in.Invite {
		cred.VerificationTokens.ResetPassword = utils.Yes
	}
	id, err := ser.app.DB.InsertCredentials(&cred)
	if err != nil {
		return nil, errors.New("User already exists")
	}
	org, err := service.GetOrg(nil, in.OrgID)
	if err != nil {
		return nil, err
	}
	// if OTP is enabled, return without sending verification email. Except for invite
	if org.Settings.Auth.OTP && !in.Invite {
		return &pb.CredentialsResponse{
			Id: int32(id),
		}, err
	}
	var name string
	if in.FirstName == "" && in.LastName == "" {
		name = "There"
	} else {
		name = in.FirstName + " " + in.LastName
	}
	var template string
	templateData := map[string]string{
		"name":    name,
		"orgName": org.Info["name"].(string),
		"url":     org.Info["host"].(string),
		"token":   token}
	if in.Invite {
		template = "invite-email-verification"
	} else {
		template = "register-email-verification"
		templateData["senderName"] = in.SenderName
	}
	//send verification email
	kafka.SendEmail(&kafka.SESTemplate{
		Template:     template,
		TemplateData: templateData,
		Source:       org.Info["email"].(string),
		ToAddresses:  []string{in.Email},
	})
	return &pb.CredentialsResponse{
		Id: int32(id),
	}, err
}

func (g *grpcServer) UpdateCredentials(ctx context.Context, in *pb.CredentialsRequest) (*emptypb.Empty, error) {
	creds, err := g.app.DB.GetCredentialsWithUserID(in.UserID)
	if err != nil {
		return nil, err
	}
	if in.Email != "" {
		creds.Email = in.Email
	}
	if in.HashedPassword != "" {
		creds.Password = in.HashedPassword
	}
	if in.Salt != "" {
		creds.Salt = in.Salt
	}
	if in.FirstName != "" {
		creds.FirstName = in.FirstName
	}
	if in.LastName != "" {
		creds.LastName = in.LastName
	}
	if in.Status != "" {
		creds.Status = service.UserStatus(in.Status)
	}
	err = g.app.DB.UpdateCredentialsWithCredentials(creds)
	return &emptypb.Empty{}, err
}

func (g *grpcServer) DeleteCredentials(ctx context.Context, in *pb.DeleteCredentialsRequest) (*emptypb.Empty, error) {
	err := g.app.DB.DeleteCredentialsWithUserID(in.UserID)
	if err != nil {
		return nil, err
	}
	err = g.app.DB.DeleteResetPasswordTokenWithUserID(in.UserID)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, err
}

func (ser *grpcServer) VerifyPassword(ctx context.Context, in *pb.VerifyPasswordRequest) (*pb.VerifyPasswordResponse, error) {
	creds, err := ser.app.DB.GetCredentialsWithUserID(in.UserID)
	res := &pb.VerifyPasswordResponse{
		IsVerified: false,
	}
	if err != nil {
		return res, errors.New("user not found")
	}
	ok := utils.CheckPasswordHash(in.Password, creds.Password, creds.Salt)
	if !ok {
		return res, errors.New("password not match")
	}
	res.IsVerified = true
	return res, nil
}

func (ser *grpcServer) UpdatePassword(ctx context.Context, in *pb.PasswordUpdateRequest) (*emptypb.Empty, error) {
	err := ser.app.DB.UpdatePasswordWithUserID(in.UserID, in.Password, in.Salt)
	return &emptypb.Empty{}, err
}

func (ser *grpcServer) UpdateName(ctx context.Context, in *pb.UpdateNameRequest) (*emptypb.Empty, error) {
	err := ser.app.DB.UpdateNameWithUserID(in.UserID, in.FirstName, in.LastName)
	return &emptypb.Empty{}, err
}

func (ser *grpcServer) UpdateEmail(ctx context.Context, in *pb.UpdateEmailRequest) (*emptypb.Empty, error) {
	err := ser.app.DB.UpdateEmailWithUserID(in.UserID, in.Email)
	return &emptypb.Empty{}, err
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
	}
	if err != nil {
		fmt.Println("")
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedJSON, w, err)
		return
	}
	email := req.Payload.Data["email"].(string)
	password := req.Payload.Data["password"].(string)
	if email == "" {
		utils.NewHttpError("", "", "Email can't be empty", utils.ErrorCodeNone, http.StatusNotAcceptable, err).Log(nil, 3).Write(w)
		return
	}
	orgID, err := service.GetOrgID(r)
	if err != nil {
		utils.NewHttpError("", "", "Server was unable to process request", utils.ErrorCodeNone, http.StatusInternalServerError, err).Log(nil, 1).Write(w)
		return
	}
	creds, err := app.DB.GetCredentialsWithEmailAndOrgID(email, orgID.Hex())
	if creds == nil || creds.Email == "" {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeNone, w, err, nil)
		return
	}
	//creds := user.Credentials[0]
	veri := creds.VerificationTokens
	if veri.Verified == utils.No {
		orgSettings, err := service.GetOrg(r, "")
		//send verification email
		kafka.SendEmail(&kafka.SESTemplate{
			Template: "register-email-verification",
			TemplateData: map[string]string{
				"name":    creds.FirstName + " " + creds.LastName,
				"orgName": orgSettings.Info["name"].(string),
				"url":     orgSettings.Info["host"].(string),
				"token":   creds.VerificationTokens.Token.String},
			ToAddresses: []string{creds.Email},
			Source:      orgSettings.Info["email"].(string),
		})
		utils.NewHttpError("", "", "Your email address needs to be verified. We have just sent you an email for the verification process. Please check it and follow the next steps.", utils.ErrorCodeEmailHasNotVerified, http.StatusNotAcceptable, err).Log(nil, 3).Write(w)
		return
	}

	ok := utils.CheckPasswordHash(password, creds.Password, creds.Salt)
	if !ok {
		utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeNone, w, errors.New("password not match"))
		return
	}
	//user from user-service
	user, err := service.GetUserWithRole(creds.UserID, "", "", req.Host, utils.ServiceOrg, utils.CtxOrgForm, utils.RoleMember, true, true)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err)
		return
	}
	if user.Status != service.UserStatusActive {
		utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeNone, w, errors.New("inactive user"))
		return
	}
	orgSettings, err := service.GetOrg(r, "")
	utils.TermDebugging(`orgSettings`, orgSettings)
	if orgSettings.Auth.Mfa {
		token, err := app.sendMFACode(orgSettings, creds)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
			return
		}

		res := utils.Response{
			Result: utils.SUCCESS,
			Data:   map[string]any{"mfa": true, "mfaToken": token},
		}
		utils.WriteJSON(w, http.StatusOK, res)
		return
	}
	if err != nil {
		utils.NewHttpError("", "", "Unable to authorize", utils.ErrorCodeGetOrgSettings, http.StatusInternalServerError, err).Log(nil, 1, r).Write(w)
		return
	}
	/////
	// tokens, claim, err := svc.GenerateTokenPair(user, &orgSettings.Settings, r)
	// if err != nil {
	// 	utils.NewHttpError("", "", "Unable to authorize", utils.ErrorCodeFailedJWT, http.StatusInternalServerError, err).Log(nil, 1, user, orgSettings).Write(w)
	// 	return
	// }

	t, err := app.postLogin(user, orgSettings, w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	res := utils.Response{
		Result: utils.SUCCESS,
		Data:   t,
	}
	utils.WriteJSON(w, http.StatusOK, res)
	return
}

func (app *application) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token := ""
	host, err := utils.GetXHost(r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeXHost, w, err)
	}
	config := utils.GetJWTConfig(host)
	orgSettings, err := service.GetOrg(r, "")
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeGetOrgSettings, w, err)
	}
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeGetOrgSettings, w, err, r)
		return
	}
	for _, cookie := range r.Cookies() {
		if cookie.Name == config.CookieName() {
			token = cookie.Value
		}
	}
	if token == "" {
		utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeJWTTokenIsMissingInCookie, w, nil)
		return
	}
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.RefreshSecret()), nil
	})
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedJWT, w, err)
		return
	}
	user, err := service.GetUserWithRole(claims["sub"].(string), "", "", host, utils.ServiceOrg, utils.CtxOrgForm, utils.RoleMember, true, true)
	// //user.ProfileImage, err = utils.GetCloudFrontSignedURL(user.ProfileImage)
	// cookies, err := utils.GetCloudFrontSignedCookie(user.ProfileImage)
	// utils.TermDebugging(`cookies`, cookies)
	// if err != nil {
	// 	utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err, claims["sub"].(string), []string{service.FormNameUserMetadata}, []string{"_"})
	// 	return
	// }
	// for _, cookie := range cookies {
	// 	utils.TermDebugging(`cookie`, cookie)
	// 	http.SetCookie(w, cookie)
	// }
	t, err := app.postLogin(user, orgSettings, w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	res := utils.Response{
		Result: utils.SUCCESS,
		Data:   t,
	}
	utils.WriteJSON(w, http.StatusOK, res)
	return
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	// delete auth cookie
	host, _ := utils.GetXHost(r)
	config := utils.GetJWTConfig(host)
	cookie := http.Cookie{
		Name:     config.CookieName(),
		Path:     "/",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Domain:   host,
		HttpOnly: true,
		Secure:   false,
	}
	http.SetCookie(w, &cookie)

	// delete cloudfront cookie
	awsCookie, _ := utils.GetCloudFrontSignedCookie()
	awsCookieDomain := os.Getenv("AWS_CLOUDFRONT_SIGNED_COOKIES_DOMAIN")
	for _, cookie := range awsCookie {
		cookie.Domain = awsCookieDomain
		cookie.Secure = true
		cookie.MaxAge = -1
		cookie.Expires = time.Unix(0, 0)
		cookie.Path = "/"
		cookie.SameSite = http.SameSiteStrictMode
		cookie.HttpOnly = false
		http.SetCookie(w, cookie)
	}

	res := utils.Response{
		Result: "ok",
	}

	utils.WriteJSON(w, http.StatusOK, res)
	return
}

func (app *application) PasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[utils.FormBuilderFields](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
	}
	orgObjID, err := service.GetOrgID(r)
	orgID := orgObjID.Hex()
	email := req.Payload.Data["email"].Value.(string)
	row, err := app.DB.GetCredentialsWithEmailAndOrgID(email, orgID)
	if err != nil {
		// fmt.Println("not good")
		// utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}

	rpt := models.ResetPasswordTokens{
		UserID:            row.UserID,
		VerificationToken: uuid.New().String(),
	}

	_, err = app.DB.InsertResetPasswordToken(&rpt)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}

	orgSettings, err := service.GetOrg(r, "")
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeGetOrgSettings, w, err, req)
		return
	}

	//send verification email
	kafka.SendEmail(&kafka.SESTemplate{
		Template: "reset-password",
		TemplateData: map[string]string{
			"title":   "Forgot Your Password?",
			"name":    row.FirstName + " " + row.LastName,
			"orgName": orgSettings.Info["name"].(string),
			"url":     "https://" + orgSettings.Info["host"].(string) + "/auth/reset-password",
			"token":   rpt.VerificationToken},
		ToAddresses: []string{row.Email},
		Source:      orgSettings.Info["email"].(string),
	})
	res := utils.Response{
		Result: utils.SUCCESS,
		Data:   nil,
	}
	utils.WriteJSON(w, http.StatusOK, res)
}

func (app *application) UpdatePasswordWithVerificationCode(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[utils.FormBuilderFields](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
	}
	token := req.Payload.Data["token"].Value.(string)
	pw := req.Payload.Data["password"].Value.(string)
	app.DB.CleanUpResetPassword(5)
	row, err := app.DB.GetResetPasswordTokenWithToken(token)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
		return
	}
	password, salt, err := utils.HashPassword(pw)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeEmptyPassword, w, err)
		return
	}
	err = app.DB.UpdatePasswordWithUserID(row.UserID, password, salt)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
		return
	}
	app.DB.DeleteResetPasswordToken(row.Model.ID)
	res := utils.Response{
		Result: utils.SUCCESS,
		Data:   nil,
	}
	utils.WriteJSON(w, http.StatusOK, res)
}

func (app *application) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	creds, err := app.DB.GetCredentialsWithVeriToken(req.Payload.Data["token"].(string))
	conn, err := grpc.Dial(os.Getenv("USER_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err, req)
		return
	}
	c := pb.NewUserServiceClient(conn)
	_, err = c.SetVerified(context.Background(), &pb.EmailVerifedRequest{
		ID: creds.UserID,
	})
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err, req)
		return
	}

	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeNoRecordFound, w, err, req)
		return
	}
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}

	output := utils.MapStringAny{
		"email":         creds.Email,
		"firstName":     creds.FirstName,
		"lastName":      creds.LastName,
		"resetPassword": creds.VerificationTokens.ResetPassword,
		"status":        creds.Status,
	}
	if creds.VerificationTokens.ResetPassword == utils.Yes {
		utils.WriteJSON(w, http.StatusAccepted, utils.Response{Result: utils.SUCCESS, Data: output})
		return
	}
	vt := models.VerificationTokens{
		Token:    sql.NullString{String: "", Valid: true},
		Verified: utils.Yes,
	}
	err = app.DB.UpdateVerificationToken(req.Payload.Data["token"].(string), &vt)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: output})
}

func (app *application) SetPassword(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	org, err := service.GetOrg(r, "")
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeGetOrgSettings, w, err)
		return
	}
	var password string
	if org.Auth.OTP {
		uuid := uuid.New().String()
		password = uuid
	} else {
		password, ok := req.Payload.Data["password"].(string)
		if !ok || password == "" {
			ed, _ := utils.NewErrorDetails("", "password", "Invalid password", utils.ErrorCodeEmptyPassword)
			utils.WriteJSON(w, http.StatusNotAcceptable, utils.Response{Result: utils.ERROR, Data: []utils.ErrorDetails{ed}})
			return
		}
		confirmPassword, ok := req.Payload.Data["confirmPassword"].(string)
		if !ok || password != confirmPassword {
			ed, _ := utils.NewErrorDetails("", "confirmPassword", "Password does not match", utils.ErrorCodeEmptyPassword)
			utils.WriteJSON(w, http.StatusNotAcceptable, utils.Response{Result: utils.ERROR, Data: []utils.ErrorDetails{ed}})
			return
		}
	}
	token := req.Payload.Data["token"].(string)
	firstName := req.Payload.Data["firstName"].(string)
	lastName := req.Payload.Data["lastName"].(string)

	creds, err := app.DB.GetCredentialsWithVeriToken(token)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
		return
	}
	hashed, salt, err := utils.HashPassword(password)
	creds.FirstName = firstName
	creds.LastName = lastName
	creds.Status = service.UserStatusActive
	creds.Salt = salt
	creds.Password = hashed
	// if firstName != "" && lastName != "" {
	// 	err = app.DB.UpdateNameWithUserID(creds.UserID, firstName, lastName)
	// 	if err != nil {
	// 		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
	// 		return
	// 	}
	// }
	// if err != nil {
	// 	utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedHashing, w, err)
	// 	return
	// }
	err = app.DB.UpdateCredentialsWithCredentials(creds)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedHashing, w, err)
		return
	}
	vt := models.VerificationTokens{
		Token:         sql.NullString{String: "", Valid: true},
		Verified:      utils.Yes,
		ResetPassword: utils.No,
	}
	conn, err := grpc.Dial(os.Getenv("USER_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err, req)
		return
	}
	c := pb.NewUserServiceClient(conn)
	_, err = c.SetVerified(context.Background(), &pb.EmailVerifedRequest{
		ID:        creds.UserID,
		FirstName: firstName,
		LastName:  lastName,
		Status:    string(service.UserStatusActive),
	})
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err, req)
		return
	}
	err = app.DB.UpdateVerificationToken(token, &vt)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func (app *application) OAuthLogin(w http.ResponseWriter, r *http.Request) {
	handler, err := getOAuthHandler(r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err)
		return
	}
	url := handler.PreOAuth(w, r)
	output := utils.MapStringAny{
		"redirect": url,
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: output})
}

func (app *application) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	handler, err := getOAuthHandler(r)
	user, err := handler.PostOAuth(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeOAuth, w, err)
		return
	}
	orgID, err := service.GetOrgID(r)
	if err != nil {
		utils.NewHttpError("", "", "Server was unable to process request", utils.ErrorCodeNone, http.StatusInternalServerError, err).Log(nil, 1).Write(w)
		return
	}
	user.OrganizationID = orgID
	userRes, err := service.ProcessOAuthUser(user)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err)
		return
	}
	user, err = service.GetUserWithRole(userRes.ID, "", "", req.Host, utils.ServiceOrg, utils.CtxOrgForm, utils.RoleMember, true, true)
	orgSettings, err := service.GetOrg(r, "")
	_, err = app.DB.GetCredentialsWithUserID(user.ID.Hex())
	if err != nil {
		uuid := uuid.New().String()
		hashedPassword, salt, err := utils.HashPassword(uuid)
		creds := models.Credentials{
			OrganizationID: orgID.Hex(),
			UserID:         user.ID.Hex(),
			Email:          user.Email,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Status:         service.UserStatusActive,
			Password:       hashedPassword,
			Salt:           salt,
		}
		_, err = app.DB.InsertCredentials(&creds)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
			return
		}
	}
	t, err := app.postLogin(user, orgSettings, w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorParsingLoginTokens, w, err, t)
		return
	}
	t.Redirect = "http://localhost:8080/user/dashboard"
	res := utils.Response{
		Result: utils.SUCCESS,
		Data:   t,
	}
	utils.WriteJSON(w, http.StatusOK, res)
	return
}

func (app *application) ResendMFACode(w http.ResponseWriter, r *http.Request) {
	app.DB.DeleteMFATokenWithInterval(3)
	req, err := utils.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	org, err := service.GetOrg(r, "")
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeGetOrgSettings, w, err)
		return
	}
	token := req.Payload.Data["token"].(string)
	utils.TermDebugging(`token`, token)
	if token == "" {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, nil)
		return
	}
	mfaToken, err := app.DB.GetMFATokenWithToken(token)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
		return
	}
	creds, err := app.DB.GetCredentialsWithUserID(mfaToken.UserID)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
		return
	}
	app.DB.DeleteMFAToken(token)
	token, err = app.sendMFACode(org, creds)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	utils.TermDebugging(`req`, req.Payload.Data["token"].(string))
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeMFA, w, err)
		return
	}
	res := utils.Response{
		Result: utils.SUCCESS,
		Data:   map[string]any{"mfa": true, "mfaToken": token},
	}
	utils.WriteJSON(w, http.StatusOK, res)
}

func (app *application) VerifyMFACode(w http.ResponseWriter, r *http.Request) {
	app.DB.DeleteMFATokenWithInterval(3)
	req, err := utils.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	token := req.Payload.Data["token"].(string)
	code := req.Payload.Data["code"].(string)
	mfaToken, err := app.DB.GetMFATokenWithTokenAndCode(token, code)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
		return
	}
	user, err := service.GetUserWithRole(mfaToken.UserID, "", "", req.Host, utils.ServiceOrg, utils.CtxOrgForm, utils.RoleMember, true, true)
	org, err := service.GetOrg(r, "")
	t, err := app.postLogin(user, org, w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	res := utils.Response{
		Result: utils.SUCCESS,
		Data:   t,
	}
	utils.WriteJSON(w, http.StatusOK, res)
}

func (app *application) SendOTP(w http.ResponseWriter, r *http.Request) {
	app.DB.DeleteOTPTokenWithInterval(3)
	req, err := utils.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	email, ok := req.Payload.Data["email"].(string)
	emptyEmail := false
	if !ok {
		emptyEmail = true
	}
	token, ok := req.Payload.Data["token"].(string)
	if !ok {
		token = ""
	}
	if emptyEmail {
		otpToken, err := app.DB.GetOTPTokenWithToken(token)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
			return
		}
		creds, err := app.DB.GetCredentialsWithUserID(otpToken.UserID)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
			return
		}
		app.DB.DeleteOTPToken(token)
		email = creds.Email
	}
	org, err := service.GetOrg(r, "")
	creds, err := app.DB.GetCredentialsWithEmailAndOrgID(email, org.ID)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
		return
	}
	newToken := uuid.New().String()
	randNumber := rand.Intn(900000) + 100000
	code := strconv.Itoa(randNumber)
	otpToken := models.OTPToken{
		Token:  sql.NullString{String: newToken, Valid: true},
		UserID: creds.UserID,
		Code:   code,
	}
	_, err = app.DB.InsertOTPToken(&otpToken)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
		return
	}
	kafka.SendEmail(&kafka.SESTemplate{
		Template: "otp-code",
		TemplateData: map[string]string{
			"title":   "Your One Time Password Code",
			"name":    creds.FirstName + " " + creds.LastName,
			"orgName": org.Info["name"].(string),
			"code":    code},
		ToAddresses: []string{creds.Email},
		Source:      org.Info["email"].(string),
	})
	res := utils.Response{
		Result: utils.SUCCESS,
		Data:   map[string]any{"token": newToken},
	}
	utils.WriteJSON(w, http.StatusOK, res)
}

func (app *application) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	app.DB.DeleteOTPTokenWithInterval(3)
	req, err := utils.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	token := req.Payload.Data["token"].(string)
	code := req.Payload.Data["code"].(string)
	otpToken, err := app.DB.GetOTPTokenWithTokenAndCode(token, code)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
		return
	}
	user, err := service.GetUserWithRole(otpToken.UserID, "", "", req.Host, utils.ServiceOrg, utils.CtxOrgForm, utils.RoleMember, true, true)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err)
		return
	}
	app.DB.DeleteOTPToken(token)
	org, err := service.GetOrg(r, "")
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	t, err := app.postLogin(user, org, w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err)
		return
	}
	res := utils.Response{
		Result: utils.SUCCESS,
		Data:   t,
	}
	utils.WriteJSON(w, http.StatusOK, res)
}

//################################################################################################################################################

type TokenResponse struct {
	*svc.TokenPairs
	Info     utils.MapStringAny `json:"info"`
	Redirect string             `json:"redirect"`
}

func (app *application) postLogin(user *service.User, org *service.Org, w http.ResponseWriter, r *http.Request) (*TokenResponse, error) {
	tokens, claim, err := svc.GenerateTokenPair(user, &org.Settings, r)
	if err != nil {
		utils.NewHttpError("", "", "Unable to authorize", utils.ErrorCodeFailedJWT, http.StatusInternalServerError, err).Log(nil, 1, user, org).Write(w)
		return nil, err
	}
	data := TokenResponse{&tokens, claim, ""}
	//set cookie to refresh the token
	host, err := utils.GetXHost(r)
	if err != nil {
		utils.NewHttpError("", "", "Unable to get host", utils.ErrorCodeFailedJWT, http.StatusInternalServerError, err).Log(nil, 1, user, org).Write(w)
		return nil, err
	}
	config := utils.GetJWTConfig(host)
	cookie := http.Cookie{
		Name:     config.CookieName(),
		Path:     "/",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(config.RefreshTokenExpiry()),
		MaxAge:   int(config.RefreshTokenExpiry().Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   host,
		HttpOnly: false,
		Secure:   true,
	}
	// cloudfront cookie
	http.SetCookie(w, &cookie)
	awsCookie, _ := utils.GetCloudFrontSignedCookie()
	awsCookieDomain := os.Getenv("AWS_CLOUDFRONT_SIGNED_COOKIES_DOMAIN")
	for _, cookie := range awsCookie {
		cookie.Domain = awsCookieDomain
		cookie.Secure = true
		cookie.HttpOnly = false
		cookie.Path = "/"
		http.SetCookie(w, cookie)
	}
	history := kafka.NewMessage(user, kafka.TopicHistory)
	kafka.SendMessage(history)
	return &data, nil
}

func (app *application) sendMFACode(org *service.Org, creds *models.Credentials) (string, error) {
	rand.Seed(time.Now().UnixNano())
	randNumber := rand.Intn(900000) + 100000
	code := strconv.Itoa(randNumber)
	mfaToken := models.MFAToken{
		UserID: creds.UserID,
		Token:  sql.NullString{String: uuid.New().String(), Valid: true},
		Code:   code,
	}
	_, err := app.DB.InsertMFAToken(&mfaToken)
	if err != nil {
		return "", err
	}
	//send verification email
	kafka.SendEmail(&kafka.SESTemplate{
		Template: "mfa-code",
		TemplateData: map[string]string{
			"title":   "Your Two-Factor Authentication Code",
			"name":    creds.FirstName + " " + creds.LastName,
			"orgName": org.Info["name"].(string),
			"code":    code},
		ToAddresses: []string{creds.Email},
		Source:      org.Info["email"].(string),
	})
	return mfaToken.Token.String, nil
}
