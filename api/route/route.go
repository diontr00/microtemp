package route

import (
	"time"

	"github.com/labstack/echo/v4"
	"{{{mytemplate}}}/api/controller"
	"{{{mytemplate}}}/api/middleware"
	"{{{mytemplate}}}/config/env"
	"{{{mytemplate}}}/rest"
	"{{{mytemplate}}}/translator"
	"{{{mytemplate}}}/validator"
)

type RouteConfig struct {
	Env        *env.RestEnv
	Timeout    time.Duration
	Rest       rest.RestServer[echo.Context]
	Validator  validator.Validator
	Translator translator.Translator
}

func Setup(config *RouteConfig) {

	main_controller := &controller.Maincontroller{
		Validator:  config.Validator,
		Translator: config.Translator,
	}

	router := config.Rest
	router.Use(middleware.NewRecoverMiddleWare())
	router.Use(middleware.NewLocaleMiddleWare())
	config.Rest.Post("/", main_controller.Hello)
}
