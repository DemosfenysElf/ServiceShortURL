package router

import (
	"github.com/caarlos0/env"
	"github.com/labstack/echo"
	"log"
)

type ConfigUrl struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}
type Server struct {
	Cfg ConfigUrl
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
