package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo"
)

func (s *serverShortener) PingDB(c echo.Context) error {
	fmt.Println("==>> PingDB")
	if s.DB == nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.DB.Ping(ctx); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
