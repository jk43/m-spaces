package hook

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moly-space/molylibs/utils"
)

func (hdr *HookHandler) LocalhostPostPutUser(ctx context.Context, tag string, claims *utils.JWTClaims, data []byte) ([]byte, error) {
	out := make(map[string]any)
	out["server"] = "localhost"
	out["data"] = "hooked"
	output := utils.Response{
		Data: out,
	}
	realOut, _ := json.Marshal(output)

	return realOut, nil
}

func (hdr *HookHandler) LocalhostPrePutUser(ctx context.Context, tag string, claims *utils.JWTClaims, data []byte) ([]byte, error) {
	type KeyValue struct {
		Key   string `json:"key"`
		Value any    `json:"value"`
	}
	type Payload struct {
		Data []KeyValue `json:"data"`
	}
	payload := Payload{}
	err := json.Unmarshal(data, &payload)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	validations := make([]utils.ErrorDetails, 0)
	for _, v := range payload.Data {
		switch v.Value.(type) {
		case string:
			if v.Value == "" {
				validations = append(validations, utils.ErrorDetails{
					Field: v.Key,
					Error: v.Key + " cannot be empty",
				})
			}
		case int:
			if v.Value == nil {
				validations = append(validations, utils.ErrorDetails{
					Field: v.Key,
					Error: v.Key + " cannot be empty",
				})
			}
		}
	}
	if len(validations) > 0 {
		res := utils.Response{
			Result: utils.ERROR,
			Data:   validations,
		}
		out, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		return nil, &utils.HookHandlerError{
			Message:      "Validation Error",
			HandlerError: out,
		}
	}

	out, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	return out, nil
}
