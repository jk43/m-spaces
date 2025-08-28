package main

import (
	"cms/models"
	"fmt"

	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
)

func getSetting(k string, m *models.Board) string {
	for _, v := range m.Settings {
		if v.K == k {
			return v.V
		}
	}
	return ""
}

func (app *application) validateBoardAccess(action string, role utils.UserRole, board *models.Board) error {
	enforcer, _ := app.CasbinManager.GetEnforcer(board.OrganizationID)
	// admin access
	// The admin should have access to inactive boards
	ok, _ := enforcer.GetRoleManager().HasLink(string(role), string(utils.RoleAdmin))
	if ok {
		return nil
	}
	if role != utils.RoleAdmin && board.Active == utils.No {
		return fmt.Errorf("board is not active")
	}
	if action == "" {
		return nil
	}
	setting := getSetting(action, board)
	ok, _ = enforcer.GetRoleManager().HasLink(string(role), string(setting))
	if !ok {
		return fmt.Errorf("you are not allowed to read this board")
	}
	return nil
}

func upsertMysqlServiceSettingWithPayload(id uint, fe service.FormElems, data utils.MapStringAny, f func(uint, service.FormElems) error) []utils.ErrorDetails {
	errs := []utils.ErrorDetails{}
	for k, v := range fe {
		if _, ok := data[k]; ok {
			err := v.SetValue(data[k].(string))
			if err != nil {
				errDetail, _ := utils.NewErrorDetails("", k, err.Error(), utils.ErrorCodeNone)
				errs = append(errs, errDetail)
			}
			v.Key = k
		}
	}
	if len(errs) > 0 {
		return errs
	}
	err := f(id, fe)
	if err != nil {
		errDetail, _ := utils.NewErrorDetails("", "", err.Error(), utils.ErrorCodeNone)
		errs = append(errs, errDetail)
	}
	return errs
}
