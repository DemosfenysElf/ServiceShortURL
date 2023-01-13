package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/labstack/echo"
)

type urlAPIShortenBath struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type shortURLApiShortenBath struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (s *Server) APIShortenBatch(c echo.Context) error {
	fmt.Println("==>> APIShortenBatch")
	urlBath := []urlAPIShortenBath{}
	shortURLBath := []shortURLApiShortenBath{}
	shortURLOne := shortURLApiShortenBath{}

	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}

	json.Unmarshal(body, &urlBath)
	var setErr error
	var short string
	for i := range urlBath {
		short, setErr = s.SetURL(urlBath[i].OriginalURL)
		shortURLOne.ShortURL = s.Cfg.BaseURL + "/" + short
		shortURLOne.CorrelationID = urlBath[i].CorrelationID
		shortURLBath = append(shortURLBath, shortURLOne)
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
	if setErr != nil {
		setErr := setErr.Error()
		if strings.Contains(setErr, pgerrcode.UniqueViolation) {
			c.Response().WriteHeader(http.StatusConflict)
		}
	} else {
		c.Response().WriteHeader(http.StatusCreated)
	}
	c.Response().Write(shortU)
	return nil
}
