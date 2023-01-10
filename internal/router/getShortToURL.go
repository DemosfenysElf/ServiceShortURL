package router

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func (s *Server) GetShortToURL(c echo.Context) error {
	fmt.Println("==>> GetShortToURL")
	short := c.Request().URL.String()
	short = short[1:]

	url, err := s.GetURL(short)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("shortURL is not exist")
	}
	c.Response().Header().Add("Location", url)
	c.Response().WriteHeader(http.StatusTemporaryRedirect)
	c.Response().Header()
	return nil
}
