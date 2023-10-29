package setup

import (
	"{{{mytemplate}}}/config/env"
	"{{{mytemplate}}}/json"
	"{{{mytemplate}}}/rest"
	"{{{mytemplate}}}/translator"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
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

		type JSON map[string]string

		var e error
		switch code {
		case http.StatusNotFound:
			e = c.JSON(code, JSON{"error": trans.TranslateMessage(locale, "notfound-Error", translator.TranslateParam{
				"Method": c.Request().Method,
				"Route":  c.Request().URL.Path,
			}, nil)})

			logger.Info().Err(err).Send()

		case http.StatusUnauthorized:
			e = c.JSON(code, JSON{"error": trans.TranslateMessage(locale, "unauthorized-Error", nil, nil)})
			logger.Info().Err(err).Send()

		case http.StatusInternalServerError:
			e = c.JSON(code, JSON{"error": trans.TranslateMessage(locale, "internal-Error", nil, nil)})

		case http.StatusBadRequest:
			//
			e = c.JSON(code, JSON{"error": trans.TranslateMessage(locale, "badrequest-Error", translator.TranslateParam{"Reason": detail}, nil)})
			logger.Info().Err(err).Send()

		case http.StatusMethodNotAllowed:

			e = c.JSON(code, JSON{"error": trans.TranslateMessage(locale, "methodnotallow-error", translator.TranslateParam{
				"Method": c.Request().Method,
				"Route":  c.Request().URL.Path,
			}, nil)})

		default:
			e = c.JSON(code, JSON{"error": trans.TranslateMessage(locale, "internal-Error", nil, nil)})

			logger.Error().Err(err).Msg("Un-handle error")

		}
		if e != nil {
			c.Logger().Error(e)
		}

	}

}
