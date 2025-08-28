package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/moly-space/molylibs/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Hooks struct {
	hooks map[string][]*pb.Hook
}

type HookRequest[T any] struct {
	Request T `json:"request"`
}

func NewHookList(service string) (*Hooks, error) {
	hooks := Hooks{
		hooks: make(map[string][]*pb.Hook),
	}
	err := hooks.GetHooksFromServer(service)
	if err != nil {
		return nil, err
	}
	return &hooks, nil
}

func (h *Hooks) GetHooksFromServer(service string) error {
	fmt.Println("GetHooks called")
	conn, err := grpc.Dial(os.Getenv("HOOK_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return err
	}
	c := pb.NewHookServiceClient(conn)
	in := pb.HookRequest{
		Service: service,
	}
	res, err := c.GetHooks(context.Background(), &in)
	if err != nil {
		return err
	}
	h.buildHooks(res)
	return nil
}

func (h *Hooks) buildHooks(hooks *pb.HookListResponse) error {
	for _, v := range hooks.Hooks {
		if h.hooks[v.Tag] == nil {
			h.hooks[v.Tag] = []*pb.Hook{v}
		} else {
			h.hooks[v.Tag] = append(h.hooks[v.Tag], v)
		}
	}
	return nil
}

func (h *Hooks) GetHooks() map[string][]*pb.Hook {
	return h.hooks
}

func (h *Hooks) GetHook(tag string) []*pb.Hook {
	return h.hooks[tag]
}

// func (h *Hooks) runRestHook(hook *pb.Hook, r *http.Request) (*http.Response, error) {
// 	url := "http://" + hook.CallbackAddr + "/" + hook.CallbackFunc
// 	body, _ := ReadAndRestoreRequestBody(r)
// 	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header = r.Header
// 	client := &http.Client{}
// 	response, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return response, nil
// }

func (h *Hooks) GetClaims(r *http.Request) []byte {
	value := r.Context().Value("claims")
	if value == nil {
		return []byte("")
	}
	claims := r.Context().Value("claims").(JWTClaims)
	claimsByte, err := json.Marshal(claims)
	if err != nil {
		return []byte("")
	}
	return claimsByte
}

func (h *Hooks) runGRPCHook(hook *pb.Hook, in *pb.RunHookRequest) (*pb.RunHookResponse, error) {
	conn, err := grpc.Dial(hook.CallbackAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewHookServiceClient(conn)
	res, err := c.RunHook(context.Background(), in)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (h *Hooks) RunPreHooks(tag string, req *http.Request) error {
	hook := h.hooks["pre/"+tag]
	if hook == nil {
		return nil
	}
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	host, err := GetXHost(req)
	in := pb.RunHookRequest{
		Data:      data,
		JwtClaims: h.GetClaims(req),
		Xhost:     host,
	}
	var hookRes *pb.RunHookResponse
	for i, v := range hook {
		in.CallbackFunc = v.CallbackFunc
		hookRes, err = h.runGRPCHook(v, &in)
		if err != nil {
			return fmt.Errorf("callbackFunc: %s, nth, %d, error: %w", v.CallbackFunc, i, err)
		}
		//error found on response for output error message
		if hookRes.Error == true {
			return &HookHandlerError{
				HandlerError: hookRes.Data,
			}
		}
		in.Data = hookRes.Data
	}
	reader := bytes.NewReader(hookRes.Data)
	readCloser := io.NopCloser(reader)
	req.Body = readCloser
	return nil
}

func (h *Hooks) RunPostHooks(tag string, req *http.Request, res *http.Response) error {
	hook := h.hooks["post/"+tag]
	if hook == nil {
		return nil
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	host, err := GetXHost(req)
	in := pb.RunHookRequest{
		Data:      data,
		JwtClaims: h.GetClaims(req),
		Xhost:     host,
	}
	var hookRes *pb.RunHookResponse
	for i, v := range hook {
		in.CallbackFunc = v.CallbackFunc
		hookRes, err = h.runGRPCHook(v, &in)
		if err != nil {
			return fmt.Errorf("callbackFunc: %s, nth, %d, error: %w", v.CallbackFunc, i, err)
		}
		//error found on response for output error message
		if hookRes.Error == true {
			return &HookHandlerError{
				HandlerError: hookRes.Data,
			}
		}
		in.Data = hookRes.Data
	}
	reader := bytes.NewReader(hookRes.Data)
	readCloser := io.NopCloser(reader)
	res.Body = readCloser
	return nil
}

type hookHandlerFunc func(context.Context, string, *JWTClaims, []byte) ([]byte, error)

type HookHandlerList struct {
	Hooks map[string]hookHandlerFunc
}

func (h *HookHandlerList) AddHook(name string, handler hookHandlerFunc) {
	h.Hooks[name] = handler
}

func (h *HookHandlerList) Setup(fn func()) {
	fn()
}

func NewHookHandlerList() *HookHandlerList {
	return &HookHandlerList{
		Hooks: make(map[string]hookHandlerFunc),
	}
}

// This error is used to send error details.
type HookHandlerError struct {
	Message      string
	Code         int
	HandlerError []byte //json.Marshaled error details.
}

func (he *HookHandlerError) Error() string {
	return fmt.Sprintf(he.Message)
}
