package setup

import (
	"embed"
	"log"
	"{{{mytemplate}}}/translator"
)

//go:embed trans_file/*.toml
var trans_folder embed.FS

// set up new translator
func NewTranslator() translator.Translator {
	trans, err := translator.New(trans_folder, "trans_file")
	if err != nil {
		log.Fatalf("[Error] Reading Translation File %v \n", err)
	}
	return trans
}
