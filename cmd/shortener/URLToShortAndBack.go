package main

import (
	"fmt"
	"github.com/labstack/echo"
	"io"
	"net/http"
)

func getShortToURL(c echo.Context) error {
	short := c.Request().URL.String()
	short = short[1:]
	if urlmap[short] == "" {
		c.Response().WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("shortURL is not exist")
	}
	c.Response().Header().Add("Location", urlmap[short])
	c.Response().WriteHeader(http.StatusTemporaryRedirect)
	c.Response().Header()
	return nil
}

func postURLToShort(c echo.Context) error {
	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}

	short := shortURL()
	if urlmap[short] != string(body) {
		for ; urlmap[short] != string(body); short = shortURL() {
			if _, ok := urlmap[short]; !ok {
				urlmap[short] = string(body)
				break
			}
		}
	}

	write := []byte("http://localhost:8080/" + short)
	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Write(write)
	return nil
}
