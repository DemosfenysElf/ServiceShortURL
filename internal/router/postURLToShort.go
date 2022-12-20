package router

import (
	"ServiceShortURL/internal/shorturlservice"
	"fmt"
	"github.com/labstack/echo"
	"io"
	"net/http"
)

func (s *Server) PostURLToShort(c echo.Context) error {
	defer c.Request().Body.Close()
	fmt.Println(">>>>>>>>>>A_1")
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println(">>>>>>>>>>A_2")
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}

	short := shorturlservice.SetURL(string(body), s.Cfg.Storage)

	write := []byte(s.Cfg.BaseURL + "/" + short)

	//if c.Request().Header.Get("Accept-Encoding") == "gzip" {
	//	fmt.Println(">>>>>>>>>>A_3")
	//	write, err = serviceCompress(write)
	//	if err != nil {
	//		fmt.Println("Compress fail")
	//	}
	//	c.Response().Header().Set("Content-Encoding", "gzip")
	//}
	fmt.Println(">>>>>>>>>>A_4")
	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Write(write)
	return nil
}
