package validator

import (
	"{{{mytemplate}}}/model"
	"{{{mytemplate}}}/translator"

	"github.com/go-playground/validator"
)

type Validator interface {
	//  internal validation of struct for request body
	validateStruct(i interface{}) []model.FieldError
	// validate i against particular tag
	ValidateAndTranslateAny(lang string, i interface{}, tag string) []model.ErrorResponse
	// validate request body and translate the error if exist according to specify locale
	ValidateRequestAndTranslate(lang string, i interface{}) []model.ErrorResponse
}

// validator that support i18n
type validatorWithTrans struct {
	logger    *zerolog.Logger
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

func (v *validatorWithTrans) ValidateRequestAndTranslate(lang string, i interface{}) []model.ErrorResponse {
	var return_err []model.ErrorResponse
	errs := v.validateStruct(i)
	if errs != nil {
		for _, e := range errs {
			msg, err := v.trans.TranslateFieldError(lang, e)
			var errMsg model.ErrorResponse

			if err != nil {
				v.trans.TranslateErrorLogger(v.logger)(e.Field(), err)

				errMsg = model.ErrorResponse{Field: e.Field(), Error: "please check the documentation"}
			} else {
				errMsg = model.ErrorResponse{Field: e.Field(), Error: msg}
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
func (v *validatorWithTrans) ValidateAndTranslateAny(lang string, i interface{}, tag string) []model.ErrorResponse {
	var return_err []model.ErrorResponse
	errs := v.validateAny(i, tag)

	if errs != nil {
		for _, e := range errs {
			msg, err := v.trans.TranslateFieldError(lang, e)

			var errMsg model.ErrorResponse
			if err != nil {

				v.trans.TranslateErrorLogger(v.logger)(e.Field(), err)

				errMsg = model.ErrorResponse{Field: e.Field(), Error: "please check the documentation"}

			} else {
				errMsg = model.ErrorResponse{Field: e.Field(), Error: msg}
			}
			return_err = append(return_err, errMsg)

		}
		return return_err
	}
	return nil
}

func New(trans translator.Translator, logger *zerolog.Logger) Validator {
	return &validatorWithTrans{
		logger:    logger,
		validator: validator.New(),
		trans:     trans,
	}
}
