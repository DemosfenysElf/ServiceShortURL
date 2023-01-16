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

type urlAPIShortenBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type shortURLApiShortenBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (s *URLServer) APIShortenBatch(c echo.Context) error {
	fmt.Println("==>> APIShortenBatch")
	urlBatch := []urlAPIShortenBatch{}
	shortURLBatch := []shortURLApiShortenBatch{}
	shortURLOne := shortURLApiShortenBatch{}

	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}

	err = json.Unmarshal(body, &urlBatch)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("unmarshal error")
	}
	var setErr error
	var short string
	for i := range urlBatch {
		short, setErr = s.SetURL(urlBatch[i].OriginalURL)
		shortURLOne.ShortURL = s.Cfg.BaseURL + "/" + short
		shortURLOne.CorrelationID = urlBatch[i].CorrelationID
		shortURLBatch = append(shortURLBatch, shortURLOne)
	}

	shortU, err := json.Marshal(shortURLBatch)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
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
	}
	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Write(shortU)
	return nil
}
