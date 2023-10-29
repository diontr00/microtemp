package setup

import (
	"{{{mytemplate}}}/translator"
	"{{{mytemplate}}}/validator"
)

// set up new validator
func NewValidator(trans translator.Translator) validator.Validator {
	return validator.New(trans)
}
