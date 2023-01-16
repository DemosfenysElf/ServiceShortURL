package router

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/labstack/echo"
)

func (s *URLServer) PostURLToShort(c echo.Context) error {
	fmt.Println("==>> PostURLToShort")

	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}
	if len(body) == 0 {
		c.Response().WriteHeader(http.StatusNoContent)
		return nil
	}
	short, setErr := s.SetURL(string(body))

	write := []byte(s.Cfg.BaseURL + "/" + short)

	if c.Request().Header.Get("Accept-Encoding") == "gzip" {
		write, err = serviceCompress(write)

		if err != nil {
			fmt.Println("Compress fail")
		}

		c.Response().Header().Set("Content-Encoding", "gzip")
	}
	if setErr != nil {
		sErr := setErr.Error()
		if strings.Contains(sErr, pgerrcode.UniqueViolation) {
			c.Response().WriteHeader(http.StatusConflict)
		}
	} else {
		c.Response().WriteHeader(http.StatusCreated)
	}
	c.Response().Write(write)
	return nil
}
