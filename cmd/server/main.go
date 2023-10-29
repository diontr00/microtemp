package main

import (
	"context"
	"sync"
	"time"

	"{{{mytemplate}}}/api/route"
	"{{{mytemplate}}}/config"
)

var (
	app_    *config.Applications
	appOnce sync.Once
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	appOnce.Do(func() {
		app_ = config.NewApp(ctx)
	})

	route.Setup(&route.RouteConfig{
		Env:        &app_.Env.App,
		Timeout:    app_.Env.App.Timeout,
		Validator:  app_.Validator,
		Translator: app_.Translator,
		Rest:       app_.Rest,
	})
	app_.Start()
}
