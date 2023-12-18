package config

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	oslog "log"
	"os"
	"os/signal"
	"{{{mytemplate}}}/config/env"
	"{{{mytemplate}}}/config/setup"
	"{{{mytemplate}}}/rest"
	"{{{mytemplate}}}/translator"
	"{{{mytemplate}}}/validator"
)

type Applications struct {
	Env        *env.Env
	Rest       rest.RestServer[echo.Context]
	Logfile    io.WriteCloser
	Translator translator.Translator
	Validator  validator.Validator
}

// Start rest server and register clean up function
func (a *Applications) Start() {
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)

	logger := a.Rest.GetGlobalLogger()

	go func() {
		<-terminate
		fmt.Printf("Gratefully Shutdown %s , Doing Cleanup Task...😷 \n", a.Env.App.AppName)
		ctx, cancel := context.WithTimeout(context.Background(), a.Env.App.Timeout)
		defer func() {
			cancel()
			if a.Logfile != nil {
				err := a.Logfile.Close()
				if err != nil {
					oslog.Fatalf("Error Closing Log File : %v", err)
				}
			}
		}()
		if err := a.Rest.Shutdown(ctx); err != nil {

			logger.Fatal().Err(err).Msg("❌ Shut down error")
		}
	}()

	err := a.Rest.Listen(a.Env.App.ListenPort)
	if err != nil {
		logger.Fatal().Err(err).Msg("❌ Shut listen error")
	}

}

func NewApp(ctx context.Context) *Applications {
	var app = new(Applications)
	app.Translator = setup.NewTranslator()
	env := env.NewEnv(ctx)
	app.Env = env
	app.Validator = setup.NewValidator(app.Translator)
	logger, logfile := setup.NewLogger(app.Env)
	app.Logfile = logfile
	app.Rest = setup.NewRest(env, app.Translator, &logger)
	app.Rest.SetGlobalLogger(&logger)

	return app
}
