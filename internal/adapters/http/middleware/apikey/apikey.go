package apikey

import (
	"net/http"
	"os"

	echo "github.com/labstack/echo/v4"
)

// OptionalCheckerFromEnv exige X-API-Key solo si existe la variable de entorno API_KEY.
func OptionalCheckerFromEnv() echo.MiddlewareFunc {
	want := os.Getenv("API_KEY")
	if want == "" {
		return func(next echo.HandlerFunc) echo.HandlerFunc { return next }
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("X-API-Key") != want {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "missing or invalid API key"})
			}
			return next(c)
		}
	}
}
