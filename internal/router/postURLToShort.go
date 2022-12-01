package router

import (
	"fmt"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"shortURLService"
)

func postURLToShort(c echo.Context) error {

	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}

	short := shortURLService.SetURL(string(body))

	write := []byte("http://localhost:8080/" + short)
	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Write(write)
	return nil
}
