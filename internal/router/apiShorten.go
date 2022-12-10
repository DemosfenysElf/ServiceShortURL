package router

import (
	"ServiceShortURL/internal/shorturlservice"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"io"
	"net/http"
)

type urlJson struct {
	Url string `json:"url"`
}

type shortUrlJson struct {
	ShortUrl string `json:"result"`
}

func ApiShorten(c echo.Context) error {
	urlJ := urlJson{}
	shortUrl := shortUrlJson{}
	defer c.Request().Body.Close()
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("URL is not exist")
	}

	json.Unmarshal(body, &urlJ)
	short := shorturlservice.SetURL(urlJ.Url)
	shortUrl.ShortUrl = "http://localhost:8080/" + short

	shortU, err := json.Marshal(shortUrl)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("Marshal error")
	}

	c.Response().Header().Add("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Write(shortU)
	return nil
}
