package json

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Custom JSON serializer
type JSONSerializer[T any] interface {
	Serialize(context T, i interface{}, indent string) error
	Deserialize(context T, i interface{}) error
}

// TODO : Implement fast json

type jsonBindingError struct {
	Field  any `json:"field"`
	Got    any `json:"got"`
	Expect any `json:"expect"`
}

type goJSON struct {
}

func (s *goJSON) Serialize(c echo.Context, i interface{}, indent string) error {

	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}

	return enc.Encode(i)
}

func (s *goJSON) Deserialize(c echo.Context, i interface{}) error {

	err := json.NewDecoder(c.Request().Body).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {

		return c.JSON(http.StatusBadRequest, map[string]jsonBindingError{"error": {Field: ute.Field, Expect: ute.Type.Name(), Got: ute.Value}})

	} else if se, ok := err.(*json.SyntaxError); ok {

		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: error=%v\n", se.Error())).SetInternal(err)
	}

	return err
}

func NewCustomJson() JSONSerializer[echo.Context] {
	return new(goJSON)
}
