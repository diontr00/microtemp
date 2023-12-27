package controller

import (
	"net/http"
	"{{{mytemplate}}}/translator"
	"{{{mytemplate}}}/validator"

	"github.com/labstack/echo/v4"
)

type Maincontroller struct {
	Validator  validator.Validator
	Translator translator.Translator
}

type User struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required"`
}

func (m *Maincontroller) Hello(c echo.Context) error {
	var user User

	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	lang := c.Get("locale").(string)

	err := m.Validator.ValidateRequestAndTranslate(lang, &user)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"user": user.Name,
	})

}
