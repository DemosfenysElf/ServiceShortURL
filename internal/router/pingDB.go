package router

import (
	"fmt"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo"
)

func (s *URLServer) PingDB(c echo.Context) error {
	fmt.Println("==>> PingDB")
	if s.DB == nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
	}
	if err := s.DB.Ping(); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
	}
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
