package translator

import (
	"embed"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
	"io/fs"
	"path/filepath"
	"{{{mytemplate}}}/model"
)

// Translate parameter that will be supplied to dictionary file, e.g: required "This field is required, plese supplied"
type TranslateParam map[string]interface{}

type Translator interface {
	// Translate Arbitrary message
	TranslateMessage(lang string, key string, param TranslateParam, plurals interface{}) (string, error)
	// Use with Validator to translate valifdation error to other locale or more meaningful error
	TranslateFieldError(lang string, field model.FieldError) (string, error)
}

type uTtrans struct {
	fs.FileInfo
	bundle *i18n.Bundle
}

// return walkdir function that load messsage file to locale bundle
func newloadMessageFile(f embed.FS, bundle *i18n.Bundle) func(path string, di fs.DirEntry, err error) error {
	return func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if di.IsDir() {
			return nil
		}
		_, err = bundle.LoadMessageFileFS(f, path)

		if err != nil {
			return err
		}
		return nil
	}
}

// return new translator when setup this server, when return error is not nil
// you should panic
func New(fs embed.FS, transFolderName string) (Translator, error) {
	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	loadMessage := newloadMessageFile(fs, bundle)
	err := filepath.WalkDir(transFolderName, loadMessage)
	if err != nil {
		return nil, err
	}

	return &uTtrans{bundle: bundle}, nil
}

func (u *uTtrans) TranslateFieldError(
	lang string,
	fe model.FieldError,
) (string, error) {

	// message localizer with message bundle have been setup
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
		return "", err
	}

	return message, nil
}

// Translate message with associated key , plurals define the plurals variable that determine more than one form
func (u *uTtrans) TranslateMessage(
	lang string,
	key string,
	para TranslateParam,
	plurals interface{},
) (string, error) {
	localizer := i18n.NewLocalizer(u.bundle, lang)
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: para,
		PluralCount:  plurals,
	})
	if err != nil {
		return "", err
	}
	return message, nil
}
