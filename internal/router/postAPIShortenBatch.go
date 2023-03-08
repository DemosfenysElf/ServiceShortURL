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
	ID          string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type shortURLApiShortenBatch struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

// PostAPIShortenBatch e.POST("/api/shorten/batch)
// принимающий в теле запроса множество URL для сокращения в формате
// [{"correlation_id": "<строковый идентификатор>","original_url": "<URL для сокращения>"},...]
// возвращает данные в формате:
// [{"correlation_id": "<строковый идентификатор из объекта запроса>","short_url": "<результирующий сокращённый URL>"},...]
func (s *serverShortener) PostAPIShortenBatch(c echo.Context) error {
	s.WG.Wait()
	fmt.Println("==>> APIShortenBatch")
	urlBatch := []urlAPIShortenBatch{}
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
	shortURLBatch := make([]shortURLApiShortenBatch, 0, len(urlBatch))
	var setErr error
	var short string
	for i := range urlBatch {
		short, setErr = s.SetURL(urlBatch[i].OriginalURL)
		shortURLOne.ShortURL = s.Cfg.BaseURL + "/" + short
		shortURLOne.ID = urlBatch[i].ID
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
