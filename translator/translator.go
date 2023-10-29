package translator

import (
	"embed"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
	"io/fs"
	"{{{mytemplate}}}/model"
	"path/filepath"
)

// Translate parameter e.g {count : 3}
type TranslateParam map[string]interface{}

type Translator interface {
	// Translate Arbitrary message
	TranslateMessage(lang string, key string, param TranslateParam, plurals interface{}) string
	// Validate struct fileds validation error helper
	TranslateRequest(lang string, field model.FieldError) string
}

type uTtrans struct {
	fs.FileInfo
	bundle *i18n.Bundle
}

func New(fs embed.FS, transFolderName string) (Translator, error) {
	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	dirs, err := fs.ReadDir(transFolderName)
	if err != nil {

		return nil, err
	}

	for _, file := range dirs {
		file_path := filepath.Join(transFolderName, file.Name())

		_, err = bundle.LoadMessageFileFS(fs, file_path)
		if err != nil {
			return nil, err

		}
	}

	return &uTtrans{bundle: bundle}, err
}

func (u *uTtrans) TranslateRequest(
	lang string,
	fe model.FieldError,
) string {

	localizer := i18n.NewLocalizer(u.bundle, lang)

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: fe.Tag(),
		TemplateData: map[string]interface{}{
			"Field": fe.Field(),
			"Param": fe.Param(),
		},
		PluralCount: nil,
	})
	if err != nil {
		return "Error Translate Message , Refer to the documentaion"
	}

	return message
}

// Translate message with associated key , plurals define the plurals variable that determine more than one form
func (u *uTtrans) TranslateMessage(
	lang string,
	key string,
	para TranslateParam,
	plurals interface{},
) string {
	localizer := i18n.NewLocalizer(u.bundle, lang)
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: para,
		PluralCount:  plurals,
	})
	if err != nil {
		return "Error TranslateMessage"
	}
	return message
}
