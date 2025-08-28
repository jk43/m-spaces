package service

// Due to the import cycle, utils.ParseRequest cannot use GetOrgID()
// This functions are used to parse the request from the service

import (
	"net/http"

	"github.com/moly-space/molylibs/utils"
)

var orgIDs = make(map[string]string)

func ParseRequest[T any](w http.ResponseWriter, r *http.Request, p ...string) (*utils.Request[T], error) {
	req, err := utils.ParseRequest[T](w, r, p...)
	if err != nil {
		return nil, err
	}
	if req.OrgID == "" {
		orgID, ok := orgIDs[req.Host]
		if !ok {
			id, err := GetOrgID(r)
			if err != nil {
				return nil, err
			}
			orgIDs[req.Host] = id.Hex()
			req.OrgID = id.Hex()
		} else {
			req.OrgID = orgID
		}
	}
	return req, nil
}

func ParseGetRequest(w http.ResponseWriter, r *http.Request, p ...utils.ChiParam) (*utils.Request[utils.MapStringSlice[[]string]], error) {
	req, err := utils.ParseGetRequest(w, r, p...)
	if err != nil {
		return nil, err
	}
	if req.OrgID == "" {
		orgID, ok := orgIDs[req.Host]
		if !ok {
			id, err := GetOrgID(r)
			if err != nil {
				return nil, err
			}
			orgIDs[req.Host] = id.Hex()
			req.OrgID = id.Hex()
		} else {
			req.OrgID = orgID
		}
	}
	return req, nil
}
