package forms

import (
	"github.com/moly-space/molylibs/service"
	"github.com/moly-space/molylibs/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (a *CMSFormElems) LocalHotdevCom(e FormElems) {
	host := utils.Host("local.hotdev.com")
	efc, err := a.Casbin.GetEnforcer(host)
	if err != nil {
		a.Logger.Error().Err(err).Msg("")
	}
	roles, err := efc.GetAllRoles()
	if err != nil {
		a.Logger.Error().Err(err).Msg("")
	}

	roleOptionValues := make([]service.OptionLabelValue, 0)
	for _, role := range roles {
		roleOptionValues = append(roleOptionValues, service.OptionLabelValue{
			Value: role,
			Label: cases.Title(language.English).String(role),
		})
	}
	e["list"] = &service.FormElem{
		Key:         "list",
		Value:       "Member",
		ValueType:   service.ServiceSettingTypeString,
		Name:        "List",
		Description: "",
		Type:        "Select",
		Editable:    true,
		Rules:       []string{"Required"},
		Options: service.FormElemOptions{
			LabelValue: roleOptionValues,
		},
		Order: 1,
	}
	e["read"] = &service.FormElem{
		Key:         "read",
		Value:       "Member",
		ValueType:   service.ServiceSettingTypeString,
		Name:        "Read",
		Description: "",
		Type:        "Select",
		Editable:    true,
		Rules:       []string{"Required"},
		Options: service.FormElemOptions{
			LabelValue: roleOptionValues,
		},
		Order: 2,
	}
	e["write"] = &service.FormElem{
		Key:         "write",
		Value:       "Member",
		ValueType:   service.ServiceSettingTypeString,
		Name:        "Write",
		Description: "",
		Type:        "Select",
		Editable:    true,
		Rules:       []string{"Required"},
		Options: service.FormElemOptions{
			LabelValue: roleOptionValues,
		},
		Order: 3,
	}
	e["update"] = &service.FormElem{
		Key:         "update",
		Value:       "au",
		ValueType:   service.ServiceSettingTypeString,
		Name:        "Update",
		Description: "",
		Type:        "Select",
		Editable:    true,
		Rules:       []string{"Required"},
		Options: service.FormElemOptions{
			LabelValue: []service.OptionLabelValue{
				{Value: "au", Label: "Author Only"},
				{Value: "ao", Label: "Admin Only"},
				{Value: "oa", Label: "Author and Admin"},
			},
		},
		Order: 4,
	}
	e["delete"] = &service.FormElem{
		Key:         "delete",
		Value:       "au",
		ValueType:   service.ServiceSettingTypeString,
		Name:        "Delete",
		Description: "",
		Type:        "Select",
		Editable:    true,
		Rules:       []string{"Required"},
		Options: service.FormElemOptions{
			LabelValue: []service.OptionLabelValue{
				{Value: "au", Label: "Author Only"},
				{Value: "ao", Label: "Admin Only"},
				{Value: "oa", Label: "Author and Admin"},
			},
		},
		Order: 5,
	}

	e["editor"] = &service.FormElem{
		Key:         "editor",
		Value:       "tiptap",
		ValueType:   service.ServiceSettingTypeString,
		Name:        "Editor",
		Description: "",
		Type:        "Select",
		Editable:    true,
		Rules:       []string{"Required"},
		Options: service.FormElemOptions{
			LabelValue: []service.OptionLabelValue{
				{Value: "textarea", Label: "textarea"},
				{Value: "editor", Label: "editor"},
				{Value: "tiptap", Label: "tiptap"},
			},
		},
		Order: 6,
	}

	e["upload"] = &service.FormElem{
		Key:         "upload",
		Value:       "false",
		ValueType:   service.ServiceSettingTypeBool,
		Name:        "upload",
		Description: "",
		Type:        "Select",
		Editable:    true,
		Rules:       []string{"Required"},
		Options: service.FormElemOptions{
			LabelValue: []service.OptionLabelValue{
				{Value: "false", Label: "Disabled"},
				{Value: "true", Label: "Enabled"},
			},
		},
		Order: 7,
	}

	e["comment"] = &service.FormElem{
		Key:         "comment",
		Value:       "disabled",
		ValueType:   service.ServiceSettingTypeString,
		Name:        "comment",
		Description: "",
		Type:        "Select",
		Editable:    true,
		Rules:       []string{"Required"},
		Options: service.FormElemOptions{
			LabelValue: []service.OptionLabelValue{
				{Value: "disabled", Label: "Disabled"},
				{Value: "single", Label: "One Level"},
				{Value: "multi", Label: "Multi Level"},
			},
		},
		Order: 8,
	}

	e["autoDragDrop"] = &service.FormElem{
		Key:         "autoDragDrop",
		Value:       "false",
		ValueType:   service.ServiceSettingTypeBool,
		Name:        "Auto Drad adn drop files to editor",
		Description: "",
		Type:        "Select",
		Editable:    true,
		Rules:       []string{"Required"},
		Options: service.FormElemOptions{
			LabelValue: []service.OptionLabelValue{
				{Value: "false", Label: "Disabled"},
				{Value: "true", Label: "Enabled"},
			},
		},
		Order: 9,
	}
}
