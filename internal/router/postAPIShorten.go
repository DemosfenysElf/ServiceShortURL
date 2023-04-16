package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

type urlJSON struct {
	URL string `json:"url"`
}

type shortURLJSON struct {
	ShortURL string `json:"result"`
}

// PostAPIShorten POST("/api/shorten")
// принимает в теле запроса JSON-объект {"url":"<some_url>"}
// возвращает в ответ объект {"result":"<shorten_url>"}.
func (s *serverShortener) PostAPIShorten(c echo.Context) error {
	s.WG.Wait()
	fmt.Println("==>> APIShorten")
	urlJ := urlJSON{}
	shortURL := shortURLJSON{}
	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}

	err = json.Unmarshal(body, &urlJ)
	if err != nil {
		c.Response().WriteHeader(http.StatusNoContent)
		return fmt.Errorf("unmarshal error")
	}

	if len(urlJ.URL) == 0 {
		c.Response().WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("URL is nil")
	}
	short, setErr := s.SetURL(c.Request().Context(), urlJ.URL)
	shortURL.ShortURL = s.Cfg.BaseURL + "/" + short

	shortU, err := json.Marshal(shortURL)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	if c.Request().Header.Get("Accept-Encoding") == "gzip" {
		shortU, err = shorturlservice.ServiceCompress(shortU)
		if err != nil {
			fmt.Println("Compress fail")
		}
		c.Response().Header().Set("Content-Encoding", "gzip")
	}

	c.Response().Header().Add("Content-Type", "application/json")
	if setErr != nil {
		sErr := setErr.Error()
		if strings.Contains(sErr, pgerrcode.UniqueViolation) {
			c.Response().WriteHeader(http.StatusConflict)
			c.Response().Write(shortU)
			return nil
		}
	}
	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Write(shortU)
	return nil
}
