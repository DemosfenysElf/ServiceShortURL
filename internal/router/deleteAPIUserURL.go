package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo"
)

func (s *serverShortener) DeleteAPIUserURL(c echo.Context) error {
	fmt.Println("==>> DeleteAPIUserURL")
	s.wg.Add(1)
	defer s.wg.Done()
	var newlist []string

	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}

	err = json.Unmarshal(body, &newlist)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("unmarshal error")
	}

	go s.DB.DeleteURL(newlist)

	c.Response().WriteHeader(http.StatusAccepted)
	return nil
}
