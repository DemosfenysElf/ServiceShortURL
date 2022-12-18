package router

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/labstack/echo"
	"log"
)

type ConfigURL struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	//envDefault:":8080"`
	BaseURL string `env:"BASE_URL"`
	//envDefault:"http://localhost:8080"`
	Storage string `env:"FILE_STORAGE_PATH"`
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

	e.GET("/:id", s.GetShortToURL)
	e.POST("/", s.PostURLToShort)
	e.POST("/api/shorten", s.APIShorten)

	errStart := e.Start(s.Cfg.ServerAddress)

	if errStart != nil {
		return errStart
	}
	return nil
}
