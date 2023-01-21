package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func (s *serverShortener) GetShortToURL(c echo.Context) error {
	s.WG.Wait()

	fmt.Println("==>> GetShortToURL")
	short := c.Request().URL.String()
	short = short[1:]

	url, err := s.GetURL(short)
	fmt.Println(url)

	if err != nil {
		sErr := err.Error()
		if strings.Contains(sErr, "deleted") {
			c.Response().WriteHeader(http.StatusGone)
			return nil
		}
		fmt.Println(">>>>>Get Status: ", http.StatusBadRequest)
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}
	fmt.Println(">>>>>Get Status: ", http.StatusTemporaryRedirect)
	c.Response().Header().Add("Location", url)
	c.Response().WriteHeader(http.StatusTemporaryRedirect)
	c.Response().Header()
	return nil
}
