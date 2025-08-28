package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"tree/models"
	"tree/service"

	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/utils"
)

func (serv *grpcServer) GetTree(ctx context.Context, in *pb.TreeRequest) (*pb.TreeResponse, error) {
	s := service.TreeRequest{
		Slug:  in.Slug,
		OrgID: in.OrgID,
	}
	s.SetDB(serv.app.DB)
	res, err := s.GetAllNodes()
	if err != nil {
		return nil, err
	}
	j, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	out := &pb.TreeResponse{
		Data: j,
	}
	return out, nil
}

func (serv *grpcServer) GetTreeAsArray(ctx context.Context, in *pb.TreeRequest) (*pb.TreeResponse, error) {
	utils.TermDebugging(`GetTreeAsArray end called`, in)
	s := service.TreeRequest{
		Slug:  in.Slug,
		OrgID: in.OrgID,
	}
	s.SetDB(serv.app.DB)
	res, err := s.GetAllNodesAsArray()
	if err != nil {
		return nil, err
	}
	j, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	out := &pb.TreeResponse{
		Data: j,
	}
	return out, nil
}

func (serv *grpcServer) IsEnd(ctx context.Context, in *pb.TreeRequest) (*pb.BoolResponse, error) {
	utils.TermDebugging(`is end called`, in)
	s := service.TreeRequest{
		Slug:  in.Slug,
		OrgID: in.OrgID,
	}
	s.SetDB(serv.app.DB)
	res, err := s.IsEnd()
	if err != nil {
		return nil, err
	}
	utils.TermDebugging(`res`, res)
	out := &pb.BoolResponse{
		Data: res,
	}
	return out, nil
}

func (serv *grpcServer) GetAncestors(ctx context.Context, in *pb.TreeRequest) (*pb.TreeResponse, error) {
	s := service.TreeRequest{
		Slug:  in.Slug,
		OrgID: in.OrgID,
	}
	s.SetDB(serv.app.DB)
	res, err := s.GetAncestors()
	if err != nil {
		return nil, err
	}
	j, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	out := &pb.TreeResponse{
		Data: j,
	}
	return out, nil
}

func (app *application) GetTree(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r, "slug")
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, r)
	}
	slug := req.Payload.Data["slug"][0]
	s := service.TreeRequest{
		Slug:  slug,
		OrgID: req.OrgID,
	}
	s.SetDB(app.DB)
	res, err := s.GetAllNodes()
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: res})
}

func (app *application) GetAdminTree(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, r)
	}
	slug := req.Payload.Data["slug"][0]
	mode, ok := req.Payload.Data["mode"]
	if !ok {
		mode = make([]string, 1)
		mode[0] = service.Single
	}
	svc := &service.TreeRequest{
		OrgID: req.OrgID,
		Slug:  slug,
	}
	svc.SetDB(app.DB)
	var res any
	switch mode[0] {
	case service.Single:
		res, err = svc.GetAllNodes()
	case service.All:
		res, err = svc.GetAllNodes()
	}
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: res})
}

func (app *application) PutAdminTreeReorder(w http.ResponseWriter, r *http.Request) {
	fmt.Println("calling")
	req, err := utils.ParseRequest[[]models.Tree](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	for i, t := range req.Payload.Data {
		// if t.Attributes.OrganizationID != req.OrgID {
		// 	utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeSecurityBreach, w, err, req)
		// 	return
		// }
		utils.TermDebugging(`t `, t.Attributes.Slug)
		_, err := app.DB.GetTreeWithSlug(req.OrgID, t.Attributes.Slug)
		if err != nil {
			utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeSecurityBreach, w, err, req)
			return
		}
		err = app.DB.UpdateOrder(&t, uint(i))
		if err != nil {
			utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedMysql, w, err, req)
			return
		}
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func (app *application) GetAdminTrees(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseGetRequest(w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, r)
	}
	t, err := app.DB.GetRootsWithOrgID(req.OrgID)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.InternalServerErr, utils.ErrorCodeFailedMysql, w, err, req)
		return
	}
	utils.TermDebugging(`t`, t)
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: t})
}

func (app *application) PostAdminTree(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[service.TreeRequest](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	utils.TermDebugging(`PostAdminTree`, req.Payload.Data)
	service := &req.Payload.Data
	service.OrgID = req.OrgID
	service.UserID = req.UserID
	service.SetDB(app.DB)
	id, slug, err := service.Insert()
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeServiceHandler, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: map[string]any{"id": id, "slug": slug}})
}

func (app *application) PutAdminTree(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[models.TreeAttribute](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	err = app.DB.UpdateTreeAttribute(req.OrgID, &req.Payload.Data)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnauthorizedRequest, utils.ErrorCodeSecurityBreach, w, err, req)
		return
	}
	utils.TermDebugging(`PutAdminTree`, req.Payload)
}

func (app *application) DeleteAdminTree(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[models.Tree](w, r)
	utils.TermDebugging(`req.Payload.Data`, req.Payload.Data)
	s := service.TreeRequest{
		Slug:  req.Payload.Data.Attributes.Slug,
		OrgID: req.OrgID,
	}
	s.SetDB(app.DB)
	err = s.DeleteTree()
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Response{Result: utils.SUCCESS, Data: nil})
}

func (app *application) PostAdminTrees(w http.ResponseWriter, r *http.Request) {
	req, err := utils.ParseRequest[service.BatchTreeRequest](w, r)
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeFailedParseRequest, w, err, req)
		return
	}
	if (req.Payload.Data.Label == "") || (req.Payload.Data.View == "") {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeServiceHandler, w, err, req)

		return
	}
	s := service.TreeRequest{}
	s.Label = req.Payload.Data.Label
	s.Description = req.Payload.Data.Description
	if req.Payload.Data.View == "" {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeServiceHandler, w, err, req)
	}
	s.View = req.Payload.Data.View
	s.OrgID = req.OrgID
	s.UserID = req.UserID
	s.SetDB(app.DB)
	_, slug, err := s.Insert()
	if err != nil {
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeServiceHandler, w, err, req)
		return
	}
	err = s.BatchProcess(&req.Payload.Data, slug)
	if err != nil {
		treeToDel := service.TreeRequest{
			Slug:  slug,
			OrgID: req.OrgID,
		}
		treeToDel.SetDB(app.DB)
		treeToDel.DeleteTree()
		utils.NewPreDefinedHttpError(utils.UnableToProcessRequest, utils.ErrorCodeServiceHandler, w, err, req)
		return
	}
	utils.TermDebugging(`req`, slug)
}
