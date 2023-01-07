package router

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"

	"github.com/labstack/echo"
	"net/http"
)

type DatabaseInterface interface {
	Connect(connStr string) error
	Close() error
	Ping() error
}

type Database struct {
	connection *sql.DB
	ctx        context.Context
}

func InitDB() (*Database, error) {
	return &Database{ctx: context.Background()}, nil
}

func (db *Database) Connect(connStr string) (err error) {
	db.connection, err = sql.Open("pgx", connStr)
	if err != nil {
		return err
	}
	return nil
}
func (db *Database) Close() error {
	return db.Close()
}

func (db *Database) Ping() error {
	ctx, cancel := context.WithTimeout(db.ctx, 1*time.Second)
	defer cancel()
	if err := db.connection.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Server) PingBD(c echo.Context) error {
	if err := s.DB.Ping(); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
	}
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
