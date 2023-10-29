package validator

import (
	"{{{mytemplate}}}/model"
	"{{{mytemplate}}}/translator"

	"github.com/go-playground/validator"
)

type Validator interface {
	//  internal validation of struct for request body
	validateStruct(i interface{}) []model.FieldError
	// validate any agains  particular tag
	ValidateAny(i interface{}, tag string) []model.FieldError
	// validate request body and translate  the error if exist according to specify language
	ValidateRequest(lang string, i interface{}) []model.RequestError
}

type validatorWithTrans struct {
	trans     translator.Translator
	validator *validator.Validate
}

func (v *validatorWithTrans) ValidateRequest(lang string, i interface{}) []model.RequestError {
	var return_err []model.RequestError
	errs := v.validateStruct(i)
	if errs != nil {
		for _, e := range errs {
			errMsg := v.trans.TranslateRequest(lang, e)
			err := model.RequestError{Field: e.Field(), Msg: errMsg}
			return_err = append(return_err, err)

		}
		return return_err
	}
	return nil

}

func (v *validatorWithTrans) validateStruct(i interface{}) []model.FieldError {
	validate_errs := v.validator.Struct(i)

	var errs []model.FieldError
	if validate_errs != nil {
		for _, e := range validate_errs.(validator.ValidationErrors) {
			errs = append(errs, e)
		}
		return errs
	}

	return nil
}

func (v *validatorWithTrans) ValidateAny(i interface{}, tag string) []model.FieldError {
	validate_errs := v.validator.Var(i, tag)

	var errs []model.FieldError
	if validate_errs != nil {
		for _, e := range validate_errs.(validator.ValidationErrors) {
			errs = append(errs, e)
		}
		return errs
	}

	return nil
}

func New(trans translator.Translator) Validator {
	return &validatorWithTrans{
		validator: validator.New(),
		trans:     trans,
	}
}
