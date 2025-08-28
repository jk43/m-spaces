package main

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/utils"
)

func (ser *grpcServer) GetHooks(ctx context.Context, in *pb.HookRequest) (*pb.HookListResponse, error) {
	docs, err := ser.app.DB.GetHooks(in.Service)
	if err != nil {
		return nil, err
	}
	hooks := make([]*pb.Hook, 0)
	for _, v := range docs {
		hooks = append(hooks, &pb.Hook{
			ID:           v.ID,
			Tag:          v.Tag,
			Priority:     v.Priority,
			CallbackType: v.CallbackType,
			CallbackAddr: v.CallbackAddr,
			CallbackFunc: v.CallbackFunc})
	}
	output := &pb.HookListResponse{
		Hooks: hooks,
	}
	return output, nil
}

func (ser *grpcServer) RunHook(ctx context.Context, in *pb.RunHookRequest) (*pb.RunHookResponse, error) {
	callbackFunc := ser.app.HookHandlers.Hooks[in.CallbackFunc]
	if callbackFunc == nil {
		err := fmt.Errorf("callbackFunc: %s, not found", in.CallbackFunc)
		ser.app.Logger.Error().Err(err).Msg("")
		return nil, err
	}
	claims := utils.JWTClaims{}
	err := json.Unmarshal(in.JwtClaims, &claims)
	if err != nil {
		ser.app.Logger.Error().Err(err).Msg("")
		return nil, err
	}
	out := pb.RunHookResponse{}
	res, err := callbackFunc(ctx, in.Xhost, &claims, in.Data)
	out.Data = res
	if err != nil {
		hookErr, ok := err.(*utils.HookHandlerError)
		if ok {
			out.Error = true
			out.Data = hookErr.HandlerError
			//to send error details as []byte, error must be nil and send the error details as []byte in data.
			return &out, nil
		}
		return nil, err
	}
	out.Error = false
	return &out, nil
}

func (app *application) Ping(w http.ResponseWriter, r *http.Request) {
	app.Logger.Info().Msg("ping")
}
