package middleware

import (
	"github.com/labstack/echo/v4/middleware"
	"{{{mytemplate}}}/rest"
)

func NewRecoverMiddleWare() rest.MiddlewareFunc {
	return rest.MiddlewareFunc(middleware.Recover())

}
