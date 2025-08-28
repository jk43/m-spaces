package forms

import (
	"github.com/moly-space/molylibs/service"
)

func (a *CMSFormElems) Base(e FormElems) {
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
			LabelValue: []service.OptionLabelValue{
				{Value: "guest", Label: "Guest"},
				{Value: "member", Label: "Member"},
				{Value: "admin", Label: "Admin"},
			},
		},
		Order: 2,
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
			LabelValue: []service.OptionLabelValue{
				{Value: "guest", Label: "Guest"},
				{Value: "member", Label: "Member"},
				{Value: "admin", Label: "Admin"},
			},
		},
		Order: 2,
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
			LabelValue: []service.OptionLabelValue{
				{Value: "guest", Label: "Guest"},
				{Value: "member", Label: "Member"},
				{Value: "admin", Label: "Admin"},
			},
		},
		Order: 1,
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
		Order: 1,
	}
	e["delete"] = &service.FormElem{
		Key:         "delete",
		Value:       "oo",
		ValueType:   service.ServiceSettingTypeString,
		Name:        "Delete",
		Description: "",
		Type:        "Select",
		Editable:    true,
		Rules:       []string{"Required"},
		Options: service.FormElemOptions{
			LabelValue: []service.OptionLabelValue{
				{Value: "oo", Label: "Author Only"},
				{Value: "ao", Label: "Admin Only"},
				{Value: "oa", Label: "Author and Admin"},
			},
		},
		Order: 1,
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
		Order: 3,
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
		Order: 4,
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
		Order: 5,
	}
}
