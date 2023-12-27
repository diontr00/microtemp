package setup

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"net/http"
	"{{{mytemplate}}}/config/env"
	"{{{mytemplate}}}/json"
	"{{{mytemplate}}}/model"
	"{{{mytemplate}}}/rest"
	"{{{mytemplate}}}/translator"
)

func NewRest(config *env.Env, trans translator.Translator, logger *zerolog.Logger) rest.RestServer[echo.Context] {

	r := rest.NewRest(logger)

	r.SetReadTimeout(config.App.ReadTimeout)
	r.SetWriteTimeout(config.App.WriteTimeout)

	r.SetErrorHandler(newErrrHandler(trans, logger))

	customJson := json.NewCustomJson()

	r.SetBodyParser(customJson)

	if !config.App.Http2 {
		r.DisableHTTP2()
	}

	if config.App.Production {
		r.HideBanner()
	}
	return r
}

func newErrrHandler(trans translator.Translator, logger *zerolog.Logger) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		code := http.StatusInternalServerError
		locale := c.Get("locale").(string)
		var detail interface{}
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			detail = he.Message
		}

		var msg string
		logTransErr := trans.TranslateErrorLogger(logger)

		switch code {
		case http.StatusNotFound:

			msg, err = trans.TranslateMessage(locale, "Notfound-Error", translator.TranslateParam{
				"Method": c.Request().Method,
				"Route":  c.Request().URL.Path,
			}, nil)

			if err != nil {
				logger.Warn().AnErr("Notfound-Error", err).Msg("Translate Message")
				logTransErr("notfound-Error", err)
				msg = "Not Found"

			}

		case http.StatusUnauthorized:
			msg, err = trans.TranslateMessage(locale, "Unauthorized-Error", nil, nil)
			if err != nil {
				logTransErr("Unauthorized-Error", err)
				msg = "Unauthorized"
			}

		case http.StatusBadRequest:

			msg, err = trans.TranslateMessage(locale, "badrequest-Error", translator.TranslateParam{"Reason": detail}, nil)
			if err != nil {
				logTransErr("Badrequest-Error", err)
				msg = "Bad Request"

			}

		case http.StatusMethodNotAllowed:

			msg, err = trans.TranslateMessage(locale, "Notallowed-Error", translator.TranslateParam{
				"Method": c.Request().Method,
				"Route":  c.Request().URL.Path,
			}, nil)
			if err != nil {
				logTransErr("Notallowed-Error", err)
				msg = "Not Allowed"
			}

		default:
			msg, err = trans.TranslateMessage(locale, "Internal-Error", nil, nil)
			if err != nil {
				logTransErr("internal-Error", err)
				msg = "Internal"
			}

		}

		e := c.JSON(code, model.ErrorResponse{Error: msg})
		if e != nil {
			logger.Err(e).Send()
		}

	}

}
