package router

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo"
	"net/http"
)

func (s *Server) PingDB(c echo.Context) error {
	if s.DB == nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
	}
	if err := s.DB.Ping(); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
	}
	c.Response().WriteHeader(http.StatusOK)
	return nil
}