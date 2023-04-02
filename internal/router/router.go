package router

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"sync"

	"ServiceShortURL/internal/shorturlservice"

	"github.com/caarlos0/env"
	"github.com/labstack/echo"
)

// ConfigURL
// SERVER_ADDRESS адрес запуска HTTP-сервера.
// BASE_URL базовый адрес результирующего сокращённого URL.
// FILE_STORAGE_PATH путь до файла должен.
// DATABASE_DSN адрес подключения к БД.
type ConfigURL struct {
	ServerAddress string `env:"SERVER_ADDRESS" json:"server_address,omitempty"`
	BaseURL       string `env:"BASE_URL" json:"base_url,omitempty"`
	Storage       string `env:"FILE_STORAGE_PATH" json:"file_storage_path,omitempty"`
	ConnectDB     string `env:"DATABASE_DSN" json:"database_dsn,omitempty"`
	EnableHTTPS   bool   `env:"ENABLE_HTTPS" json:"enable_https,omitempty"`
	Config        string `env:"CONFIG"`
}

type serverShortener struct {
	Cfg    ConfigURL
	Serv   *echo.Echo
	Writer io.Writer
	WG     *sync.WaitGroup
	DB     shorturlservice.DatabaseService
	shorturlservice.StorageInterface
	GeneratorUsers shorturlservice.GeneratorUser
}

// InitServer инициализация сервера
func InitServer() *serverShortener {
	return &serverShortener{WG: new(sync.WaitGroup), GeneratorUsers: shorturlservice.RandomGeneratorUser{}}
}

// InitTestServer инициализация сервера для тестов (пока только при работе с БД)
func InitTestServer() *serverShortener {
	return &serverShortener{WG: new(sync.WaitGroup), GeneratorUsers: shorturlservice.TestGeneratorUser{}}
}

// Router - роутер
func (s *serverShortener) Router() error {
	s.InitRouter()

	e := echo.New()

	e.Use(s.mwGzipHandle)
	e.Use(s.MWAuthentication)

	e.GET("/:id", s.GetShortToURL)
	e.GET("/api/user/urls", s.GetAPIUserURL)
	e.GET("/ping", s.GetPingDB)

	e.POST("/", s.PostURLToShort)
	e.POST("/api/shorten/batch", s.PostAPIShortenBatch)
	e.POST("/api/shorten", s.PostAPIShorten)

	e.DELETE("/api/user/urls", s.DeleteAPIUserURLs)

	RegisterPprof(e, "/debug/pprof")

	var errStart error
	if s.Cfg.EnableHTTPS {
		errStart = e.StartAutoTLS(s.Cfg.ServerAddress)
	} else {
		errStart = e.Start(s.Cfg.ServerAddress)
	}

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

// startBD подключение к БД
func (s *serverShortener) startBD() error {
	if s.Cfg.ConnectDB == "" {
		return fmt.Errorf("error s.Cfg.ConnectDB == nil")
	}

	DB := &shorturlservice.Database{RandomShort: &shorturlservice.RandomGenerator{}}

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

	if s.Cfg.Config == "" {
		flag.StringVar(&s.Cfg.Config, "-config", "", "New CONFIG")
		flag.StringVar(&s.Cfg.Config, "-c", "", "New CONFIG")
	}
	if s.Cfg.ServerAddress == "" {
		flag.StringVar(&s.Cfg.ServerAddress, "a", "", "New SERVER_ADDRESS")
	}
	if s.Cfg.BaseURL == "" {
		flag.StringVar(&s.Cfg.BaseURL, "b", "", "New BASE_URL")
	}
	if s.Cfg.Storage == "" {
		flag.StringVar(&s.Cfg.Storage, "f", "", "New FILE_STORAGE_PATH")
	}
	if s.Cfg.ConnectDB == "" {
		flag.StringVar(&s.Cfg.ConnectDB, "d", "", "New DATABASE_DSN")
	}
	if !s.Cfg.EnableHTTPS {
		flag.BoolVar(&s.Cfg.EnableHTTPS, "s", false, "New ENABLE_HTTPS")
	}
	flag.Parse()

	// чтение файла настроек
	var cfgFile ConfigURL
	if s.Cfg.Config != "" {
		file, err := os.ReadFile(s.Cfg.Config)
		if err != nil {
			log.Fatal(err)
		}
		if err = json.Unmarshal(file, &cfgFile); err != nil {
			log.Fatal(err)
		}
	}

	//проверяем данные и заполняем не заполненые
	if (s.Cfg.ServerAddress == "") && (cfgFile.ServerAddress != "") {
		s.Cfg.ServerAddress = cfgFile.ServerAddress
	} else {
		if s.Cfg.ServerAddress == "" {
			s.Cfg.ServerAddress = ":8080"
		}
	}
	if (s.Cfg.BaseURL == "") && (cfgFile.BaseURL != "") {
		s.Cfg.BaseURL = cfgFile.BaseURL
	} else {
		if s.Cfg.BaseURL == "" {
			s.Cfg.BaseURL = "http://localhost:8080"
		}
	}
	if (s.Cfg.Storage == "") && (cfgFile.Storage != "") {
		s.Cfg.Storage = cfgFile.Storage
	} else {
		if s.Cfg.Storage == "" {
			s.Cfg.Storage = "shortsURl.log"
		}
	}
	if (s.Cfg.ConnectDB == "") && (cfgFile.ConnectDB != "") {
		s.Cfg.ConnectDB = cfgFile.ConnectDB
	} else {
		if s.Cfg.ConnectDB == "" {
			s.Cfg.ConnectDB = "postgres://postgres:0000@localhost:5432/postgres"
		}
	}
	if !s.Cfg.EnableHTTPS && cfgFile.EnableHTTPS {
		s.Cfg.EnableHTTPS = cfgFile.EnableHTTPS
	}

	//для быстрого локального тестирования деградации
	//s.Cfg.Storage = ""
	//s.Cfg.ConnectDB = ""

	//выбор и инициализация хранилища
	if err := s.startBD(); err == nil {
		fmt.Println(">>>>use BD<<<<", s.Cfg.ConnectDB)
	} else if s.Cfg.Storage != "" {
		fmt.Println(">>>>use storage<<<<")
		s.StorageInterface = &shorturlservice.FileStorage{
			FilePath:    s.Cfg.Storage,
			RandomShort: &shorturlservice.RandomGenerator{},
		}
	} else {
		fmt.Println(">>>>use memory<<<<")
		s.StorageInterface = shorturlservice.InitMem()
	}
}
