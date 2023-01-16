package router

import (
	"flag"
	"fmt"
	"io"
	"log"

	"ServiceShortURL/internal/shorturlservice"

	"github.com/caarlos0/env"
	"github.com/labstack/echo"
)

type ConfigURL struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	Storage       string `env:"FILE_STORAGE_PATH"`
	ConnectDB     string `env:"DATABASE_DSN"`
}

type URLServer struct {
	Cfg    ConfigURL
	Serv   *echo.Echo
	Writer io.Writer
	DB     shorturlservice.DatabaseInterface
	shorturlservice.StorageInterface
}

func (s *URLServer) Router() error {

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
	if s.Cfg.ConnectDB == "" {
		flag.StringVar(&s.Cfg.ConnectDB, "d", "postgres://postgres:0000@localhost:5432/postgres", "New DATABASE_DSN")
	}
	flag.Parse()

	//s.Cfg.Storage = ""
	//s.Cfg.ConnectDB = ""

	if err := s.startBD(); err == nil {
		fmt.Println(">>>>use BD<<<<", s.Cfg.ConnectDB)
	} else if s.Cfg.Storage != "" {
		fmt.Println(">>>>use storage<<<<")
		s.StorageInterface = &shorturlservice.FileStorage{
			FilePath: s.Cfg.Storage,
		}
	} else {
		fmt.Println(">>>>use memory<<<<")
		s.StorageInterface = shorturlservice.InitMem()
	}

	//
	e := echo.New()

	e.Use(s.gzipHandle)
	e.Use(s.serviceAuthentication)

	e.GET("/:id", s.GetShortToURL)
	e.GET("/api/user/urls", s.APIUserURL)
	e.GET("/ping", s.PingDB)

	e.POST("/", s.PostURLToShort)
	e.POST("/api/shorten/batch", s.APIShortenBatch)
	e.POST("/api/shorten", s.APIShorten)

	errStart := e.Start(s.Cfg.ServerAddress)

	if errStart != nil {
		return errStart
	}
	return nil
}

func (s *URLServer) startBD() error {
	if s.Cfg.ConnectDB == "" {
		return fmt.Errorf("error s.Cfg.ConnectDB == nil")
	}

	// DB connection
	DB, errInit := shorturlservice.InitDB()
	if errInit != nil {
		return errInit
	}

	if errConnect := DB.Connect(s.Cfg.ConnectDB); errConnect != nil {
		return errConnect
	}
	//defer DB.Close()

	s.StorageInterface = DB
	s.DB = DB
	return nil
}
