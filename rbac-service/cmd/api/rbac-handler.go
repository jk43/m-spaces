package main

import (
	"context"
	"fmt"
	"rbac/models"

	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/utils"
)

type Policy struct {
	Service string
	Type    string `bson:"type"`
	Subject string `bson:"sub"`
	Object  string `bson:"obj"`
	Action  string `bson:"act"`
}

type Model struct {
	Service string
	Model   string `bson:"model"`
}

func (app *grpcServer) GetPolicy(ctx context.Context, in *pb.RBACRequest) (*pb.RBACResponse, error) {
	//policy, err := app.app.DB.GetPolicy(in.Host, in.Service, in.PolicyCtx)
	p, err := app.getPolicy(in.Host, in.Service, in.PolicyCtx)
	if err != nil {
		utils.TermDebugging(`err`, err)
		return nil, err
	}
	return &pb.RBACResponse{
		Policy: p,
	}, nil
}

func (app *grpcServer) GetModel(ctx context.Context, in *pb.RBACRequest) (*pb.ModelResponse, error) {
	//res := app.app.DB.GetModel(in.Host, in.Service, in.ModelCtx)
	m, err := app.getModel(in.Host, in.Service, in.ModelCtx)
	if err != nil {
		utils.TermDebugging(`err`, err)
		return nil, err
	}
	return &pb.ModelResponse{
		Model: m.Model,
	}, nil
}

func (app *grpcServer) GetModelAndPolicy(ctx context.Context, in *pb.RBACRequest) (*pb.ModelAndPolicyResponse, error) {
	// if in.ModelCtx == "" {
	// 	return nil, errors.New("model context is required")
	// }
	// if in.PolicyCtx == "" {
	// 	return nil, errors.New("policy context is required")
	// }
	//p, err := app.app.DB.GetPolicy(in.Host, in.Service, in.PolicyCtx)
	p, err := app.getPolicy(in.Host, in.Service, in.PolicyCtx)
	if err != nil {
		return nil, err
	}
	m, err := app.getModel(in.Host, in.Service, in.ModelCtx)
	if err != nil {
		return nil, err
	}
	return &pb.ModelAndPolicyResponse{
		Model:  m.Model,
		Policy: p,
	}, nil
}

func policiesToStrings(pols []models.Policy) []string {
	var res []string
	for _, pol := range pols {
		res = append(res, fmt.Sprintf("%s,%s,%s,%s", pol.Type, pol.Subject, pol.Object, pol.Action))
	}
	return res
}

func (srv *grpcServer) getPolicy(host, service, ctx string) ([]string, error) {
	pols := []models.Policy{}
	// group policy with host
	gp, err := srv.app.DB.GetPolicies("g", host, service, ctx)
	if err != nil {
		return nil, err
	}
	if len(gp) == 0 {
		// global group policy
		gp, err = srv.app.DB.GetPolicies("g", "", "", "")
		if err != nil {
			return nil, err
		}
	}
	pols = append(pols, gp...)
	// policy with host
	pWithHost, err := srv.app.DB.GetPolicies("p", host, service, ctx)
	if err != nil {
		return nil, err
	}
	pWithoutHost, err := srv.app.DB.GetPolicies("p", "", service, ctx)
	if err != nil {
		return nil, err
	}
	for _, pwoh := range pWithoutHost {
		found := false
		for _, pwh := range pWithHost {
			if pwoh.Object == pwh.Object && pwoh.Action == pwh.Action {
				pols = append(pols, pwh)
				found = true
				continue
			}
		}
		if !found {
			pols = append(pols, pwoh)
		}
	}
	return policiesToStrings(pols), nil
}

func (srv *grpcServer) getModel(host, service, ctx string) (*models.Model, error) {
	// model with service and host
	m := srv.app.DB.GetModel2(host, service, ctx)
	if m == nil {
		// global model
		m = srv.app.DB.GetModel2("", "", "")
	}
	return m, nil
}
