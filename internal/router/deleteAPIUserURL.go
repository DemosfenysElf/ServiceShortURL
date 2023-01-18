package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

func (s *serverShortener) DeleteAPIUserURL(c echo.Context) error {
	fmt.Println("==>> DeleteAPIUserURL")
	s.WG.Add(1)
	defer s.WG.Done()
	var newlist []string

	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}
	user := shorturlservice.GetCookieValue(c.Request().Cookies())

	err = json.Unmarshal(body, &newlist)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("unmarshal error")
	}

	go s.DB.DeleteURL(user, newlist)

	c.Response().WriteHeader(http.StatusAccepted)
	return nil
}
