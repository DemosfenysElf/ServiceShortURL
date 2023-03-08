package router

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"sync"

	"ServiceShortURL/internal/shorturlservice"

	"github.com/caarlos0/env"
	"github.com/labstack/echo"
)

// SERVER_ADDRESS адрес запуска HTTP-сервера.
// BASE_URL базовый адрес результирующего сокращённого URL.
// FILE_STORAGE_PATH путь до файла должен.
// DATABASE_DSN адрес подключения к БД.
type ConfigURL struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	Storage       string `env:"FILE_STORAGE_PATH"`
	ConnectDB     string `env:"DATABASE_DSN"`
}

type serverShortener struct {
	Cfg    ConfigURL
	Serv   *echo.Echo
	Writer io.Writer
	WG     *sync.WaitGroup
	DB     shorturlservice.DatabaseService
	shorturlservice.StorageInterface
}

func InitServer() *serverShortener {
	return &serverShortener{WG: new(sync.WaitGroup)}
}

// Router - роутер
func (s *serverShortener) Router() error {
	s.InitRouter()

	e := echo.New()

	e.Use(s.mwGzipHandle)
	e.Use(s.mwAuthentication)

	e.GET("/:id", s.GetShortToURL)
	e.GET("/api/user/urls", s.GetAPIUserURL)
	e.GET("/ping", s.GetPingDB)

	e.POST("/", s.PostURLToShort)
	e.POST("/api/shorten/batch", s.PostAPIShortenBatch)
	e.POST("/api/shorten", s.PostAPIShorten)

	e.DELETE("/api/user/urls", s.DeleteAPIUserURL)

	RegisterPprof(e, "/debug/pprof")

	errStart := e.Start(s.Cfg.ServerAddress)

	if errStart != nil {
		return errStart
	}
	return nil
}

// Пакет хендлеров pprof.
func RegisterPprof(e *echo.Echo, prefixOptions string) {
	prefixRouter := e.Group(prefixOptions)
	{
		prefixRouter.GET("/", handler(pprof.Index))
		prefixRouter.GET("/allocs", handler(pprof.Handler("allocs").ServeHTTP))
		prefixRouter.GET("/block", handler(pprof.Handler("block").ServeHTTP))
		prefixRouter.GET("/cmdline", handler(pprof.Cmdline))
		prefixRouter.GET("/goroutine", handler(pprof.Handler("goroutine").ServeHTTP))
		prefixRouter.GET("/heap", handler(pprof.Handler("heap").ServeHTTP))
		prefixRouter.GET("/mutex", handler(pprof.Handler("mutex").ServeHTTP))
		prefixRouter.GET("/profile", handler(pprof.Profile))
		prefixRouter.POST("/symbol", handler(pprof.Symbol))
		prefixRouter.GET("/symbol", handler(pprof.Symbol))
		prefixRouter.GET("/threadcreate", handler(pprof.Handler("threadcreate").ServeHTTP))
		prefixRouter.GET("/trace", handler(pprof.Trace))
	}
}

func handler(h http.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}
}

func (s *serverShortener) startBD() error {
	if s.Cfg.ConnectDB == "" {
		return fmt.Errorf("error s.Cfg.ConnectDB == nil")
	}

	DB := &shorturlservice.Database{}

	if errConnect := DB.Connect(s.Cfg.ConnectDB); errConnect != nil {
		return errConnect
	}

	s.StorageInterface = DB
	s.DB = DB
	return nil
}

// InitRouter парсим флаги
// флаг -a, отвечающий за адрес запуска HTTP-сервера (переменная SERVER_ADDRESS);
// флаг -b, отвечающий за базовый адрес результирующего сокращённого URL (переменная BASE_URL);
// флаг -f, отвечающий за путь до файла с сокращёнными URL (переменная FILE_STORAGE_PATH);
// флаг -d, отвечающий за путь до DB (переменная DATABASE_DSN).
//
// Подключаемся к БД, если не получается проверяем к файлу с данными,
// если не получается, то храним данные в памяти.
func (s *serverShortener) InitRouter() {
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

	//для быстрого локального тестирования деградации
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

}
