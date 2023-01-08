package router

import (
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
	if len(body) == 0 {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}
	short := s.SetURL(string(body))

	write := []byte(s.Cfg.BaseURL + "/" + short)

	if c.Request().Header.Get("Accept-Encoding") == "gzip" {
		write, err = serviceCompress(write)

		if err != nil {
			fmt.Println("Compress fail")
		}

		c.Response().Header().Set("Content-Encoding", "gzip")
	}

	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Write(write)
	return nil
}
