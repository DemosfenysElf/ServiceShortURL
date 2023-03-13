package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

// DeleteAPIUserURL e.DELETE("/api/user/urls")
// принимает список идентификаторов сокращённых URL для удаления в формате:
// [ "a", "b", "c", "d", ...]
// в случае успешного приёма запроса возвращает http.StatusAccepted.
func (s *serverShortener) DeleteAPIUserURL(c echo.Context) error {
	fmt.Println("==>> DeleteAPIUserURL")
	//s.WG.Wait()
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

	if len(newlist) == 0 {
		c.Response().WriteHeader(http.StatusAccepted)
		return nil
	}

	go s.Delete(c.Request().Context(), user, newlist)
	time.Sleep(1 * time.Second)
	c.Response().WriteHeader(http.StatusAccepted)
	return nil
}
