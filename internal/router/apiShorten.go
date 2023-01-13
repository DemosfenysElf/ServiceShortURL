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

type urlJSON struct {
	URL string `json:"url"`
}

type shortURLJSON struct {
	ShortURL string `json:"result"`
}

func (s *Server) APIShorten(c echo.Context) error {
	fmt.Println("==>> APIShorten")
	urlJ := urlJSON{}
	shortURL := shortURLJSON{}
	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}

	json.Unmarshal(body, &urlJ)
	short, setErr := s.SetURL(urlJ.URL)
	shortURL.ShortURL = s.Cfg.BaseURL + "/" + short

	shortU, err := json.Marshal(shortURL)
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
		sErr := setErr.Error()
		if strings.Contains(sErr, pgerrcode.UniqueViolation) {
			c.Response().WriteHeader(http.StatusConflict)
		}
	} else {
		c.Response().WriteHeader(http.StatusCreated)
	}
	c.Response().Write(shortU)
	return nil
}
