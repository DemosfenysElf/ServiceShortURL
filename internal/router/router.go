package router

import (
	"github.com/caarlos0/env"
	"github.com/labstack/echo"
	"log"
)

type ConfigURL struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}
type Server struct {
	Cfg  ConfigURL
	Serv *echo.Echo
}

func (s *Server) Router() error {

	errConfig := env.Parse(&s.Cfg)
	if errConfig != nil {
		log.Fatal(errConfig)
	}

	e := echo.New()

	e.GET("/:id", s.GetShortToURL)
	e.POST("/", s.PostURLToShort)
	e.POST("/api/shorten", s.APIShorten)

	errStart := e.Start(s.Cfg.ServerAddress)
	if errStart != nil {
		return errStart
	}
	return nil
}
