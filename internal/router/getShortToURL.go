package router

import (
	"ServiceShortURL/internal/shorturlservice"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

func (s *Server) GetShortToURL(c echo.Context) error {
	short := c.Request().URL.String()
	short = short[1:]

	url, err := shorturlservice.GetURL(short, s.Cfg.Storage)
	if err != nil {
		log.Fatal(err)
	}
	if url == "" {
		c.Response().WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("shortURL is not exist")
	}

	c.Response().Header().Add("Location", url)
	c.Response().WriteHeader(http.StatusTemporaryRedirect)
	c.Response().Header()
	return nil
}
