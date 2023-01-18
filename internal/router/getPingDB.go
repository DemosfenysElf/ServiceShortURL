package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo"
)

func (s *serverShortener) GetPingDB(c echo.Context) error {
	s.WG.Wait()
	fmt.Println("==>> PingDB")
	if s.DB == nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.DB.Ping(ctx); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
