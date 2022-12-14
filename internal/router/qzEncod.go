package router

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/labstack/echo"
	"io"
)

func (s Server) gzipHandle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println(">>>> Content: ", c.Request().Header.Get("Content-Encoding"))
		if c.Request().Header.Get("Content-Encoding") != "gzip" {
			return next(c)
		}

		qzBody, err := gzip.NewReader(c.Request().Body)
		if err != nil {
			return fmt.Errorf("qz is not exist")
		}

		body, err := io.ReadAll(qzBody)
		if err != nil {
			c.Error(echo.ErrInternalServerError)
			return fmt.Errorf("URL is not exist")
		}
		fmt.Println(string(body))
		stringReader := bytes.NewReader(body)

		c.Request().Body = io.NopCloser(stringReader)

		return next(c)
	}
}
