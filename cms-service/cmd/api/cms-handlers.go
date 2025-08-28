package main

import (
	"cms/models"
	"context"
	"fmt"
	"net/http"

	"github.com/moly-space/molylibs/kafka"
	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
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
	if in.ServiceCtx == "posts" {
		res.MaxFileSize = 10 << 20
		res.AllowedContentTypes = []string{"*"}
		res.AllowedFileTypes = []string{}
	}
	return res, nil
}

func (serv *grpcFileClientServer) SaveFileData(ctx context.Context, in *pb.FileSaveDataRequest) (*pb.FileSaveDataResponse, error) {
	// profileImage := service.Metadata{
	// 	"profileImage": in.FileInfos[0].S3Key,
	// }
	// err := serv.app.DB.UpdateMetadataWithID(in.UserID, profileImage, nil)
	// if err != nil {
	// 	return nil, err
	// }
	file := &models.File{
		UserID:           in.UserID,
		FileID:           in.FileInfos[0].UuID,
		ContentType:      in.FileInfos[0].ContentType,
		OriginalFileName: in.FileInfos[0].OriginalFileName,
		S3Path:           in.FileInfos[0].S3Key,
		FormID:           in.FormID,
		Size:             uint64(in.FileInfos[0].Size),
	}
	serv.app.DB.InsertFile(file)
	return &pb.FileSaveDataResponse{}, nil
}

func (app *application) Ping(w http.ResponseWriter, r *http.Request) {
	utils.TermDebugging(`Ping`, nil)
}

func (app *application) GetBoards(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	pn := utils.Pagination{}
	sortBy := req.Payload.Data["sortBy"][0]
	descending := req.Payload.Data["descending"][0] == "true"
	utils.QueryExtractor(&req.Payload.Data, &pn)

	//boards, err := app.DB.GetBoardsByOrgID(req.OrgID, nil, "")
	boards, total, err := app.DB.GetBoardsWithOrgIDAndPagination(req.OrgID, &pn, "", sortBy, descending)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: map[string]any{"rows": boards, "total": total}})
}

func (app *application) PostBoard(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	slug := utils.GenerateSlug(app.DB.GetBoardSlugWithSlug, req.Payload.Data["slug"].(string), req.OrgID)
	board := &models.Board{
		OrganizationID: req.OrgID,
		Name:           req.Payload.Data["name"].(string),
		Slug:           slug,
		Active:         utils.YesOrNo(req.Payload.Data["active"].(string)),
		CreatedBy:      req.UserID,
	}
	id, err := app.DB.InsertBoard(board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	// dependency := map[string]any{
	// 	"id": id,
	// }
	// settings := service.OrgServiceSettings{
	// 	Dependency:        dependency,
	// 	ServiceSettingsFn: BoardServiceSettings,
	// }
	errs := upsertMysqlServiceSettingWithPayload(id, app.Forms.GetForm(req.Host), req.Payload.Data, app.DB.UpsertSettings)
	if len(errs) > 0 {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: errs})
		return
	}

	// settings. ref.

	// if err != nil {
	// 	utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
	// 	return
	// }
	//
	// Org service settings ref.
	//
	// app.Forms.GetForm(req.Host)
	// settings := service.OrgServiceSettings{
	// 	OrgID:     req.OrgID,
	// 	Ctx:       "board",
	// 	Slug:      slug,
	// 	FormElems: app.Forms.GetForm(req.Host),
	// }
	// settings.SavePayload(req.Payload.Data)
	// _ = service.SetServiceSettingsWithPayload(&configs, req.Payload.Data)
	// ref := SettingRef{
	// 	ID: id,
	// }
	//_ = UpsertMysqlServiceSettingWithPayload(BoardServiceSettings, &ref, req.Payload.Data, app.DB.UpsertSettings)
	// if len(errs) > 0 {
	// 	utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: errs})
	// 	return
	// }
	history := kafka.NewMessage(CMSHistory{
		BoardSlug: board.Slug,
		Action:    "PostBoard",
		UserID:    req.UserID,
		OrgID:     req.OrgID,
		IP:        req.IP,
	}, kafka.TopicHistory)
	kafka.SendMessage(history)
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: map[string]any{"slug": slug}})
}

// type SettingRef struct {
// 	ID uint
// }

func (app *application) PutBoard(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	id := uint(req.Payload.Data["id"].(float64))
	board := &models.Board{
		OrganizationID: req.OrgID,
		Name:           req.Payload.Data["name"].(string),
		ID:             id,
		Active:         utils.YesOrNo(req.Payload.Data["active"].(string)),
	}
	err = app.DB.UpdateBoard(board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	app.Forms.GetForm(req.Host)
	settings := service.OrgServiceSettings{
		OrgID:     req.OrgID,
		Ctx:       "board",
		Slug:      req.Payload.Data["slug"].(string),
		FormElems: app.Forms.GetForm(req.Host),
	}
	errs1 := settings.Save(req.Payload.Data)
	if len(errs1) > 0 {
		utils.TermDebugging(`errs1`, errs1)
	}
	errs := upsertMysqlServiceSettingWithPayload(id, app.Forms.GetForm(req.Host), req.Payload.Data, app.DB.UpsertSettings)
	if len(errs) > 0 {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Response{Result: utils.ERROR, Data: errs})
		return
	}
	history := kafka.NewMessage(CMSHistory{
		BoardSlug: board.Slug,
		Action:    "PutBoard",
		UserID:    req.UserID,
		OrgID:     req.OrgID,
		IP:        req.IP,
	}, kafka.TopicHistory)
	kafka.SendMessage(history)
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: map[string]any{"id": id}})
}

func (app *application) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseRequest[models.Board](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	id := req.Payload.Data.ID
	board := &models.Board{
		OrganizationID: req.OrgID,
		ID:             uint(id),
	}

	err = app.DB.DeleteBoard(board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	history := kafka.NewMessage(CMSHistory{
		BoardSlug: board.Slug,
		Action:    "DeleteBoard",
		UserID:    req.UserID,
		OrgID:     req.OrgID,
		IP:        req.IP,
	}, kafka.TopicHistory)
	kafka.SendMessage(history)
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: map[string]any{"id": id}})
}

func (app *application) GetForm(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	s := app.Forms.GetForm(req.Host)
	utils.TermDebugging(`s`, s)
	utils.TermDebugging(`req.Host`, req.Host)
	// settings := service.OrgServiceSettings{
	// 	ServiceSettingsFn: BoardServiceSettings,
	// }
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: s})
}

func (app *application) GetBoard(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	slug := req.Payload.Data["slug"][0]
	board := &models.Board{
		OrganizationID: req.OrgID,
		Slug:           slug,
	}
	err = app.DB.GetBoard(board)
	err = app.validateBoardAccess("", req.Role, board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeCMSAccessViolation, w, err, req)
		return
	}
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: board})
}

func (app *application) GetPost(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	slug := req.Payload.Data["slug"][0]
	post := &models.Post{
		Slug: slug,
	}
	err = app.DB.GetPost(post)
	board := &models.Board{
		OrganizationID: req.OrgID,
		ID:             post.BoardID,
	}
	err = app.DB.GetBoard(board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	err = app.validateBoardAccess("read", req.Role, board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeCMSAccessViolation, w, err, req)
		return
	}
	output := struct {
		Post     *models.Post     `json:"post"`
		Comments []models.Comment `json:"comments"`
		Board    *models.Board    `json:"board"`
		Files    []models.File    `json:"files"`
	}{
		Post:     post,
		Board:    board,
		Comments: []models.Comment{},
		Files:    []models.File{},
	}
	// comments
	commentLevel := board.GetSetting("comment", "nono")
	if commentLevel != "disabled" {
		// comments, err := app.DB.GetCommentsWithPostID(post.ID)
		// if err != nil {
		// 	utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		// 	return
		// }
		// post.Comments = comments
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: output})
}

func (app *application) GetPosts(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	slug := req.Payload.Data["slug"][0]
	pn := utils.Pagination{}
	sortBy := req.Payload.Data["sortBy"][0]
	descending := req.Payload.Data["descending"][0] == "true"
	utils.QueryExtractor(&req.Payload.Data, &pn)
	board := &models.Board{
		OrganizationID: req.OrgID,
		Slug:           slug,
	}
	err = app.DB.GetBoard(board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	err = app.validateBoardAccess("list", req.Role, board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeCMSAccessViolation, w, err, req)
		return
	}
	posts, total, err := app.DB.GetPostsWithBoardIDAndPagination(board.ID, &pn, "", sortBy, descending)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
	}
	for _, post := range posts {
		utils.TermDebugging(`post`, post)
		user, err := service.GetUser(post.PostAuthor.UserID, "", req.OrgID, false, false, []string{})
		if err != nil {
			user = &service.User{
				FirstName: "",
				LastName:  "",
			}
		}
		totalComments, _ := app.DB.GetTotalCommentsWithID(post.ID)
		post.TotalComments = totalComments

		post.UserName = &models.UserName{
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: map[string]any{"rows": posts, "total": total}})
}

func (app *application) PostPosts(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	if req.UserID == "" && req.Payload.Data["code"] == nil {
		org, err := service.GetOrg(r, "")
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeGetOrgSettings, w, err, r)
			return
		}
		gvc := &models.GuestVerificationCode{
			Email:  req.Payload.Data["email"].(string),
			PostID: nil,
			Method: models.MethodPost,
		}
		err = app.DB.InsertGuestVerificationCode(gvc)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
			return
		}
		kafka.SendEmail(&kafka.SESTemplate{
			Template: "token-verification",
			TemplateData: map[string]string{
				"name":    req.Payload.Data["name"].(string),
				"orgName": org.Info["name"].(string),
				"url":     org.Info["host"].(string),
				"code":    gvc.Code},
			ToAddresses: []string{req.Payload.Data["email"].(string)},
			Source:      org.Info["email"].(string),
		})
		vi := utils.VerificationInstruction{
			URL:     "/board/posts",
			KeyName: "code",
			Method:  "POST",
			Message: "A verification code has been sent to your email. Please enter the code to complete the process.",
		}
		utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.VERIFICATION_REQUIRED, Data: vi})
		return
	}

	if req.Payload.Data["code"] != nil {
		gvc := &models.GuestVerificationCode{
			Email:  req.Payload.Data["email"].(string),
			Method: models.MethodPost,
		}
		err = app.DB.GetGuestVerificationCode(gvc)
		utils.TermDebugging(`gvc`, gvc)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
			return
		}
		if gvc.Code != req.Payload.Data["code"].(string) {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeVerificationCodeNotMatch, w, err, req)
			return
		}
	}

	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	slug := utils.GenerateSlug(app.DB.GetPostSlugWithSlugAndParams, req.Payload.Data["slug"].(string))
	board := &models.Board{
		OrganizationID: req.OrgID,
		Slug:           req.Payload.Data["boardSlug"].(string),
	}
	err = app.DB.GetBoard(board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	post := &models.Post{
		BoardID:  board.ID,
		Slug:     slug,
		ParentID: nil,
		Title:    req.Payload.Data["title"].(string),
		Text:     req.Payload.Data["text"].(string),
	}
	parentSlug := req.Payload.Data["parentSlug"].(string)
	if parentSlug != "" {
		parent := &models.Post{
			BoardID: board.ID,
			Slug:    parentSlug,
		}
		err = app.DB.GetPost(parent)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
			return
		}
		post.RootID = parent.RootID
		post.ParentID = &parent.ID
	}
	id, err := app.DB.InsertPost(post)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	postAuthor := &models.PostAuthor{}
	postAuthor.PostID = id
	if req.UserID != "" {
		postAuthor.UserID = req.UserID
	} else {
		postAuthor.Name = req.Payload.Data["name"].(string)
		postAuthor.Email = req.Payload.Data["email"].(string)
	}
	utils.TermDebugging(`postAuthor`, postAuthor)
	err = app.DB.InsertPostAuthor(postAuthor)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	// attachments
	if formID, ok := req.Payload.Data["formId"]; ok {
		err = app.DB.UpdateFilePostIDWithFormID(formID.(string), id)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
			return
		}
	}
	//root post
	if parentSlug == "" {
		postsWithParent := &models.Post{
			ID:       id,
			ParentID: &id,
			RootID:   &id,
		}
		err = app.DB.UpdatePostParent(postsWithParent)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
			return
		}
	}
	history := kafka.NewMessage(CMSHistory{
		BoardSlug: board.Slug,
		PostSlug:  post.Slug,
		Action:    "PostPosts",
		UserID:    req.UserID,
		OrgID:     req.OrgID,
		IP:        req.IP,
	}, kafka.TopicHistory)
	kafka.SendMessage(history)
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: map[string]any{"id": id}})
}

func (app *application) PutPosts(w http.ResponseWriter, r *http.Request) {
	// due to circle dependency, we need to use service.ParseRequest instead of utils.ParseRequest
	req, err := service.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	board := &models.Board{
		OrganizationID: req.OrgID,
		Slug:           req.Payload.Data["boardSlug"].(string),
	}
	err = app.DB.GetBoard(board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	post := &models.Post{
		BoardID: board.ID,
		Slug:    req.Payload.Data["slug"].(string),
	}
	err = app.DB.GetPost(post)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	// validate
	setting := getSetting("update", board)
	switch setting {
	case "au": // author only
		if post.PostAuthor.UserID != req.UserID {
			utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeCMSAccessViolation, w, fmt.Errorf("you are not the owner of this post"), req)
			return
		}
	case "ao": // admin only
		if req.Role != utils.RoleAdmin {
			utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeCMSAccessViolation, w, fmt.Errorf("you are not the admin"), req)
			return
		}
	case "oa": // author and admin
		enforcer, _ := app.CasbinManager.GetEnforcer(req.Host)
		ok, _ := enforcer.GetRoleManager().HasLink(string(req.Role), string(utils.RoleAdmin))
		if !ok && post.PostAuthor.UserID != req.UserID {
			utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeCMSAccessViolation, w, fmt.Errorf("you are not the owner of this post"), req)
			return
		}
	default:
		utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeCMSAccessViolation, w, fmt.Errorf("Could not validate access"), req)
		return
	}

	// if post.CreatedBy != req.UserID && req.Role != utils.RoleAdmin {
	// 	utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, fmt.Errorf("you are not the owner of this post"), req)
	// 	return
	// }
	// if err != nil {
	// 	utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
	// 	return
	// }
	// utils.TermDebugging(`setting`, setting)

	post.Title = req.Payload.Data["title"].(string)
	post.Text = req.Payload.Data["text"].(string)
	err = app.DB.UpdatePost(post)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	} // attachments
	if formID, ok := req.Payload.Data["formId"]; ok {
		err = app.DB.UpdateFilePostIDWithFormID(formID.(string), post.ID)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
			return
		}
	}
	history := kafka.NewMessage(CMSHistory{
		BoardSlug: board.Slug,
		PostSlug:  post.Slug,
		Action:    "PutPosts",
		UserID:    req.UserID,
		OrgID:     req.OrgID,
		IP:        req.IP,
	}, kafka.TopicHistory)
	kafka.SendMessage(history)
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: map[string]any{"id": post.ID}})
}

func (app *application) DeletePosts(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseRequest[utils.MapStringAny](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	board := &models.Board{
		OrganizationID: req.OrgID,
		Slug:           req.Payload.Data["boardSlug"].(string),
	}
	err = app.DB.GetBoard(board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	post := &models.Post{
		BoardID: board.ID,
		Slug:    req.Payload.Data["slug"].(string),
	}
	err = app.DB.GetPost(post)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	// validate
	setting := getSetting("update", board)
	switch setting {
	case "au": // author only
		if post.PostAuthor.UserID != req.UserID {
			utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeCMSAccessViolation, w, fmt.Errorf("you are not the owner of this post"), req)
			return
		}
	case "ao": // admin only
		if req.Role != utils.RoleAdmin {
			utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeCMSAccessViolation, w, fmt.Errorf("you are not the admin"), req)
			return
		}
	case "oa": // author and admin
		enforcer, _ := app.CasbinManager.GetEnforcer(req.Host)
		ok, _ := enforcer.GetRoleManager().HasLink(string(req.Role), string(utils.RoleAdmin))
		if !ok && post.PostAuthor.UserID != req.UserID {
			utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeCMSAccessViolation, w, fmt.Errorf("you are not the owner of this post"), req)
			return
		}
	default:
		utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeCMSAccessViolation, w, fmt.Errorf("Could not validate access"), req)
		return
	}
	err = app.DB.DeletePost(post)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	history := kafka.NewMessage(CMSHistory{
		BoardSlug: board.Slug,
		PostSlug:  post.Slug,
		Action:    "DeletePosts",
		UserID:    req.UserID,
		OrgID:     req.OrgID,
		IP:        req.IP,
	}, kafka.TopicHistory)
	kafka.SendMessage(history)
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: map[string]any{"id": post.ID}})
}

func (app *application) GetComments(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	slug := req.Payload.Data["slug"][0]
	pn := utils.Pagination{}
	sortBy := req.Payload.Data["sortBy"][0]
	descending := req.Payload.Data["descending"][0] == "true"
	utils.QueryExtractor(&req.Payload.Data, &pn)
	// board := &models.Board{
	// 	OrganizationID: req.OrgID,
	// 	Slug:           boardSlug,
	// }
	// err = app.DB.GetBoard(board)
	// if err != nil {
	// 	utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
	// 	return
	// }
	post := &models.Post{
		Slug: slug,
	}
	err = app.DB.GetPost(post)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
	}
	board := &models.Board{
		OrganizationID: req.OrgID,
		ID:             post.BoardID,
	}
	err = app.DB.GetBoard(board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	err = app.validateBoardAccess("read", req.Role, board)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeCMSAccessViolation, w, err, req)
		return
	}
	posts, total, err := app.DB.GetCommentsWithIDAndPagination(post.ID, &pn, "", sortBy, descending)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
	}
	// for _, post := range posts {
	// 	user, err := service.GetUser(post.PostAuthor.UserID, "", req.OrgID, false, false, []string{})
	// 	if err != nil {
	// 		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
	// 		return
	// 	}

	// 	//totalComments, _ := app.DB.GetTotalCommentsWithID(post.ID)
	// 	//post.TotalComments = totalComments

	// 	post.UserName = &models.UserName{
	// 		FirstName: user.FirstName,
	// 		LastName:  user.LastName,
	// 	}
	// }
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: map[string]any{"rows": posts, "total": total}})
}

func (app *application) GetFiles(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	var formId string
	var slug string
	files := []*models.File{}
	_formID, ok := req.Payload.Data["formId"]
	if ok {
		formId = _formID[0]
	}
	_slug, ok := req.Payload.Data["slug"]
	if ok {
		slug = _slug[0]
	}
	// output orphan files
	if formId != "" {
		app.DB.UpdateFileFormIDWithUserID(req.UserID, req.Payload.Data["formId"][0])
		files, err = app.DB.GetOrphanFilesWithUserID(req.UserID)
	} else if slug != "" {
		post := &models.Post{
			Slug: slug,
		}
		app.DB.GetPost(post)
		files, err = app.DB.GetFilesWithPostID(post.ID)
		// utils.TermDebugging(`req.Payload.Data["formId"][0]`, req.Payload.Data["formId"][0])
	}
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: files})
}

func (app *application) DeleteFile(w http.ResponseWriter, r *http.Request) {
	req, err := service.ParseRequest[utils.MapStringAny](w, r)
	id := req.Payload.Data["id"].(string)
	file := &models.File{
		FileID: id,
		UserID: req.UserID,
	}
	err = app.DB.GetFile(file)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	_, err = service.DeleteFile(req.UserID, req.OrgID, file.FileID)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedGRPC, w, err, req)
	}
	err = app.DB.DeleteFile(file)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: map[string]any{"id": req.Payload.Data["id"]}})
}
