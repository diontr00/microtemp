package setup

import (
	"{{{mytemplate}}}/translator"
	"{{{mytemplate}}}/validator"
)
	"github.com/rs/zerolog"

// set up new validator
func NewValidator(trans translator.Translator , l *zerolog.Logger) validator.Validator {
	return validator.New(trans , l)
}
