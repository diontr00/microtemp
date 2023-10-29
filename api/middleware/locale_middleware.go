package middleware

import (
	"{{{mytemplate}}}/rest"

	"github.com/labstack/echo/v4"
)

func NewLocaleMiddleWare() rest.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			locale := c.QueryParam("locale")
			switch locale {
			case "en":
				locale = "en-US"
			case "vi":
				locale = "vi-VN"
			default:
				locale = "en-US"
			}
			c.Set("locale", locale)
			return next(c)
		}
	}

}
