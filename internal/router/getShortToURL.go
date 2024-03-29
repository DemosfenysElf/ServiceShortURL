package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// GetShortToURL GET("/:id")
// принимает в качестве URL-параметра идентификатор сокращённого URL
// и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location
func (s *ServerShortener) GetShortToURL(c echo.Context) error {
	s.WG.Wait()

	fmt.Println("==>> GetShortToURL")
	short := c.Request().URL.String()
	short = short[1:]

	url, err := s.GetLongURL(c.Request().Context(), short)

	if err != nil {
		sErr := err.Error()
		if strings.Contains(sErr, "deleted") {
			c.Response().WriteHeader(http.StatusGone)
			return nil
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}
	fmt.Println("==>> GetShortToURL", url)
	c.Response().Header().Add("Location", url)
	c.Response().WriteHeader(http.StatusTemporaryRedirect)
	c.Response().Header()
	return nil
}
