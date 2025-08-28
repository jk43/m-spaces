package service

import (
	"github.com/moly-space/molylibs/utils"
)

type FieldsToDelete utils.MapStringAny
type FieldsToUpdate utils.MapStringAny

type KeyvValuePairPayload map[string]utils.KeyValuePair

// This function returns the data necessary to update MongoDB from the Payload passed from the FormBuilder.
// Note that the function returns an ErrorDetails slice, not an error
//
//	formName: The parameter represents the name of the form to be updated.
//	current: The parameter represents the current data in the form.
//	payload: The parameter represents the data to be updated.
//	settings: The parameter represents the settings for the form.
//	fieldFunc: The parameter represents the function that will be used to generate the field name. This function will be called with the form name and the field name(Key in FormBuilderField) as parameters.
func ParseDataToUpdate(role utils.UserRole, casbin *utils.Casbin, formName FormName, oldData utils.MapStringAny, newData utils.FormBuilderFields, settings *Settings, fieldFunc func(string, string) string) (FieldsToUpdate, FieldsToDelete, []utils.ErrorDetails) {
	toDelete := make(utils.MapStringAny)
	toUpdate := make(utils.MapStringAny)
	// find fields to delete
	for k, _ := range oldData {
		fieldFound := false
		for _, f := range settings.Forms[formName] {
			if k == f.Key {
				fieldFound = true
				break
			}
		}
		if !fieldFound {
			toDelete[fieldFunc(formName, k)] = ""
		}
	}
	var errDetails []utils.ErrorDetails
	// find fields to add or update
	for k, v := range newData {
		for _, f := range settings.Forms[formName] {
			if k == f.Key {
				edit, _ := casbin.HasLink(role, f.Edit)
				if !edit {
					errDetails = append(errDetails, utils.ErrorDetails{
						StructField: "",
						Field:       k,
						Code:        0,
						Error:       f.Name + " is not editable",
					})
					continue
				}
				field := fieldFunc(formName, k)
				if f.DefaultValue != "" && v.Value == "" {
					toUpdate[field] = f.DefaultValue
				} else {
					toUpdate[field] = v.Value
				}
				break
			}
		}
	}
	if len(errDetails) > 0 {
		return nil, nil, errDetails
	}
	return toDelete, toUpdate, nil
}

// Parse the data from the form inputs and return the data that is shareable.
func GetDataWithFormInputs(data utils.MapStringAny, formInputs []FormInput, shareable bool) utils.MapStringAny {
	output := utils.MapStringAny{}
	for k, v := range data {
		//For service.FormNameOrgInfo is not accessible by RoleMember, can't use getForms(RoleMember... function
		for _, f := range formInputs {
			if f.Key == k {
				if !shareable {
					output[k] = v
					continue
				}
				if f.Shareable {
					output[k] = v
				}
			}
		}
	}
	return output
}
