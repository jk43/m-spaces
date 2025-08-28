package main

import (
	"context"
	"errors"
	"os"

	"github.com/moly-space/molylibs/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCClient struct {
	app *application
}

type GRPCClientInterface interface {
	SaveCredentials(in *pb.CredentialsRequest) (*pb.CredentialsResponse, error)
	VerifyPassword(in *pb.VerifyPasswordRequest) (*pb.VerifyPasswordResponse, error)
	GetPolicy(in *pb.RBACRequest) (*pb.RBACResponse, error)
	GetModel(in *pb.RBACRequest) (*pb.ModelResponse, error)
	UpdatePassword(in *pb.PasswordUpdateRequest) (*emptypb.Empty, error)
	UpdateName(in *pb.UpdateNameRequest) (*emptypb.Empty, error)
	UpdateEmail(in *pb.UpdateEmailRequest) (*emptypb.Empty, error)
	DeleteCredentials(in *pb.DeleteCredentialsRequest) (*emptypb.Empty, error)
	UpdateCredentials(in *pb.CredentialsRequest) (*emptypb.Empty, error)
}

func (gc *GRPCClient) SaveCredentials(in *pb.CredentialsRequest) (*pb.CredentialsResponse, error) {
	conn, err := grpc.Dial(os.Getenv("AUTH_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewAuthServiceClient(conn)
	return c.SaveCredentials(context.Background(), in)
}

func (gc *GRPCClient) VerifyPassword(in *pb.VerifyPasswordRequest) (*pb.VerifyPasswordResponse, error) {
	conn, err := grpc.Dial(os.Getenv("AUTH_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewAuthServiceClient(conn)
	return c.VerifyPassword(context.Background(), in)
}

func (gc *GRPCClient) UpdatePassword(in *pb.PasswordUpdateRequest) (*emptypb.Empty, error) {
	conn, err := grpc.Dial(os.Getenv("AUTH_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewAuthServiceClient(conn)
	return c.UpdatePassword(context.Background(), in)
}

func (gc *GRPCClient) UpdateEmail(in *pb.UpdateEmailRequest) (*emptypb.Empty, error) {
	conn, err := grpc.Dial(os.Getenv("AUTH_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewAuthServiceClient(conn)
	return c.UpdateEmail(context.Background(), in)
}

func (gc *GRPCClient) UpdateName(in *pb.UpdateNameRequest) (*emptypb.Empty, error) {
	conn, err := grpc.Dial(os.Getenv("AUTH_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewAuthServiceClient(conn)
	return c.UpdateName(context.Background(), in)
}

func (gc *GRPCClient) GetPolicy(in *pb.RBACRequest) (*pb.RBACResponse, error) {
	conn, err := grpc.Dial(os.Getenv("RBAC_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewRBACServiceClient(conn)
	return c.GetPolicy(context.Background(), in)
}

func (gc *GRPCClient) GetModel(in *pb.RBACRequest) (*pb.ModelResponse, error) {
	conn, err := grpc.Dial(os.Getenv("RBAC_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewRBACServiceClient(conn)
	return c.GetModel(context.Background(), in)
}
func (gc *GRPCClient) DeleteCredentials(in *pb.DeleteCredentialsRequest) (*emptypb.Empty, error) {
	conn, err := grpc.Dial(os.Getenv("AUTH_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewAuthServiceClient(conn)
	return c.DeleteCredentials(context.Background(), in)
}

func (gc *GRPCClient) UpdateCredentials(in *pb.CredentialsRequest) (*emptypb.Empty, error) {
	conn, err := grpc.Dial(os.Getenv("AUTH_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewAuthServiceClient(conn)
	return c.UpdateCredentials(context.Background(), in)
}

type GRPCTestClient struct{}

func (grpc *GRPCTestClient) SaveCredentials(in *pb.CredentialsRequest) (*pb.CredentialsResponse, error) {
	if in.Email == "grpcerror@test.com" {
		return nil, errors.New("GRPC error")
	}
	return nil, nil
}
