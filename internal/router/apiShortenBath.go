package router

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"io"
	"net/http"
)

type urlApiShortenBath struct {
	CorrelationId string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type shortURLApiShortenBath struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (s *Server) APIShortenBatch(c echo.Context) error {
	fmt.Println("==>> APIShortenBatch")
	urlBath := []urlApiShortenBath{}
	shortURLBath := []shortURLApiShortenBath{}
	shortURLOne := shortURLApiShortenBath{}

	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}

	json.Unmarshal(body, &urlBath)

	for i := range urlBath {
		short := s.SetURL(urlBath[i].OriginalURL)
		shortURLOne.ShortURL = s.Cfg.BaseURL + "/" + short
		shortURLOne.CorrelationId = urlBath[i].CorrelationId
		shortURLBath = append(shortURLBath, shortURLOne)
		//shortURLBath[i].ShortURL = s.Cfg.BaseURL + "/" + short
		//shortURLBath[i].CorrelationId = urlBath[i].CorrelationId
	}

	shortU, err := json.Marshal(shortURLBath)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("marshal error")
	}

	if c.Request().Header.Get("Accept-Encoding") == "gzip" {
		shortU, err = serviceCompress(shortU)
		if err != nil {
			fmt.Println("Compress fail")
		}
		c.Response().Header().Set("Content-Encoding", "gzip")
	}

	c.Response().Header().Add("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Write(shortU)
	return nil
}
