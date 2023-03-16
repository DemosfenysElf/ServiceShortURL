package router

import (
	"fmt"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo"
)

// GetPingDB проверяет соединение с базой данных.
// При успешной проверке хендлер должен вернуть HTTP-статус 200 OK,
// при неуспешной — 500 Internal Server Error
func (s *serverShortener) GetPingDB(c echo.Context) error {
	s.WG.Wait()
	fmt.Println("==>> PingDB")
	if s.DB == nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}
	if err := s.DB.Ping(c.Request().Context()); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
