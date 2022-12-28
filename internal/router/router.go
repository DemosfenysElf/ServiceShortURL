package router

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/labstack/echo"
	"io"
	"log"
)

type ConfigURL struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	Storage       string `env:"FILE_STORAGE_PATH"`
}

type Server struct {
	Cfg    ConfigURL
	Serv   *echo.Echo
	Writer io.Writer
}

func (s *Server) Router() error {

	errConfig := env.Parse(&s.Cfg)
	if errConfig != nil {
		log.Fatal(errConfig)
	}

	if s.Cfg.ServerAddress == "" {
		flag.StringVar(&s.Cfg.ServerAddress, "a", ":8080", "New SERVER_ADDRESS")
	}
	if s.Cfg.BaseURL == "" {
		flag.StringVar(&s.Cfg.BaseURL, "b", "http://localhost:8080", "New BASE_URL")
	}
	if s.Cfg.Storage == "" {
		flag.StringVar(&s.Cfg.Storage, "f", "shortsURl.log", "New FILE_STORAGE_PATH")
	}
	flag.Parse()

	e := echo.New()

	e.Use(s.gzipHandle)

	e.Use(s.serviceAuthentication)

	e.GET("/:id", s.GetShortToURL)
	e.POST("/", s.PostURLToShort)
	e.POST("/api/shorten", s.APIShorten)
	e.POST("/api/user/urls", s.APIUserURL)

	errStart := e.Start(s.Cfg.ServerAddress)

	if errStart != nil {
		return errStart
	}
	return nil

}
