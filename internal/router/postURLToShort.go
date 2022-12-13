package router

import (
	"ServiceShortURL/internal/shorturlservice"
	"fmt"
	"github.com/labstack/echo"
	"io"
	"net/http"
)

func (s *Server) PostURLToShort(c echo.Context) error {

	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}

	short := shorturlservice.SetURL(string(body))

	//write := []byte(s.Cfg.BaseURL + short)
	write := []byte("http://localhost:8080/" + short)
	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Write(write)
	return nil
}
