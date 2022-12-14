package router

import (
	"ServiceShortURL/internal/shorturlservice"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"io"
	"net/http"
)

type urlJSON struct {
	URL string `json:"url"`
}

type shortURLJSON struct {
	ShortURL string `json:"result"`
}

func (s *Server) APIShorten(c echo.Context) error {
	urlJ := urlJSON{}
	shortURL := shortURLJSON{}
	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}

	json.Unmarshal(body, &urlJ)
	short := shorturlservice.SetURL(urlJ.URL)
	shortURL.ShortURL = s.Cfg.BaseURL + "/" + short

	shortU, err := json.Marshal(shortURL)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("marshal error")
	}

	c.Response().Header().Add("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Write(shortU)
	return nil
}
