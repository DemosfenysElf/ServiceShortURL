package router

import (
	"ServiceShortURL/internal/shortURLService"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func GetShortToURL(c echo.Context) error {
	short := c.Request().URL.String()
	short = short[1:]

	url := shortURLService.GetURL(short)

	if url == "" {
		c.Response().WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("shortURL is not exist")
	}
	c.Response().Header().Add("Location", url)
	c.Response().WriteHeader(http.StatusTemporaryRedirect)
	c.Response().Header()
	return nil
}
