package validator

import (
	"{{{mytemplate}}}/model"
	"{{{mytemplate}}}/translator"

	"github.com/go-playground/validator"
)

type Validator interface {
	//  internal validation of struct for request body
	validateStruct(i interface{}) []model.FieldError
	// validate any agains particular tag
	ValidateAndTranslateAny(lang string, i interface{}, tag string) []model.FieldErrorResponse
	// validate request body and translate the error if exist according to specify locale
	ValidateRequestAndTranslate(lang string, i interface{}) []model.FieldErrorResponse
}

// validator that support i18n
type validatorWithTrans struct {
	trans     translator.Translator
	validator *validator.Validate
}

func (v *validatorWithTrans) validateStruct(i interface{}) []model.FieldError {
	// validate struct against criteria  tag
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

func (v *validatorWithTrans) ValidateRequestAndTranslate(lang string, i interface{}) []model.FieldErrorResponse {
	var return_err []model.FieldErrorResponse
	errs := v.validateStruct(i)
	if errs != nil {
		for _, e := range errs {
			msg, err := v.trans.TranslateFieldError(lang, e)
			var errMsg model.FieldErrorResponse

			if err != nil {

				log.Err(err).Msgf("Couldn't translate error when validation into %s", lang)
				errMsg = model.FieldErrorResponse{Field: e.Field(), Msg: "[error]"}
			} else {
				errMsg = model.FieldErrorResponse{Field: e.Field(), Msg: msg}
			}
			return_err = append(return_err, errMsg)

		}
		return return_err
	}
	return nil

}

func (v *validatorWithTrans) validateAny(i interface{}, tag string) []model.FieldError {
	// validate struct against criteria  tag
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

// helper to manuallty validate interface against any tag
func (v *validatorWithTrans) ValidateAndTranslateAny(lang string, i interface{}, tag string) []model.FieldErrorResponse {
	var return_err []model.FieldErrorResponse
	errs := v.validateAny(i, tag)

	if errs != nil {
		for _, e := range errs {
			msg, err := v.trans.TranslateFieldError(lang, e)

			var errMsg model.FieldErrorResponse
			if err != nil {
				log.Err(err).Msgf("Couldn't translate error when validation into %s", lang)

				errMsg = model.FieldErrorResponse{Field: e.Field(), Msg: "[error]"}
			} else {
				errMsg = model.FieldErrorResponse{Field: e.Field(), Msg: msg}
			}
			return_err = append(return_err, errMsg)

		}
		return return_err
	}
	return nil
}

func New(trans translator.Translator) Validator {
	return &validatorWithTrans{
		validator: validator.New(),
		trans:     trans,
	}
}
