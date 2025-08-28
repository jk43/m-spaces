// Package molylibs provides utility functions for data validation.
package molylibs

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/moly-space/molylibs/utils"
)

// Validator is a struct that holds a Translator and a Validate object.
type Validator struct {
	Translator ut.Translator       // Translator object for translating validation errors
	Validator  *validator.Validate // Validate object for validating struct data
}

// NewValidator creates a new Validator object.
// If a Translator object is provided, it uses that. Otherwise, it creates a new English Translator.
func NewValidator(t ut.Translator) *Validator {
	var v Validator
	if t != nil {
		v.Translator = t
		v.Validator = validator.New()
		en_translations.RegisterDefaultTranslations(v.Validator, v.Translator)
		return &v
	}
	//return validator with english trans
	en := en.New()
	uni := ut.New(en, en)
	v.Translator, _ = uni.GetTranslator("en")
	v.Validator = validator.New()
	en_translations.RegisterDefaultTranslations(v.Validator, v.Translator)
	return &v
}

// Validate validates a struct using the Validator's Validate object.
// It returns an error and a slice of ErrorDetails if validation fails.
func (v *Validator) Validate(s any) (error, []utils.ErrorDetails) {
	vld := v.Validator
	// Register a function to get the name of the struct field from the "json" tag.
	vld.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Validate the struct.
	err := vld.Struct(s)
	if err == nil {
		return nil, nil
	}

	// If validation fails, get the validation errors.
	errs := err.(validator.ValidationErrors)
	var errDetails []utils.ErrorDetails

	// Translate each validation error and add it to the slice of ErrorDetails.
	for _, e := range errs {
		msg := e.Translate(v.Translator)
		errDetails = append(errDetails, utils.ErrorDetails{
			StructField: e.StructNamespace(),
			Field:       e.Field(),
			Code:        0,
			Error:       msg,
		})
	}
	return err, errDetails
}
