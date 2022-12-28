package router

import (
	"ServiceShortURL/internal/shorturlservice"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

type massiveURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s *Server) APIUserURL(c echo.Context) error {
	userCoockies := shorturlservice.GetStructCoockies()
	allURL := make([]massiveURL, 0)
	element := massiveURL{}
	consumerURL, err := shorturlservice.NewConsumer(s.Cfg.Storage)
	if err != nil {
		log.Fatal(err)
	}
	defer consumerURL.Close()

	for {
		readURL, err := consumerURL.ReadURL()
		if err != nil {
			break
		}
		if readURL.CookiesAuthentication.ValueUser == userCoockies.ValueUser {
			element.ShortURL = s.Cfg.BaseURL + readURL.ShortURL
			element.OriginalURL = readURL.URL
			allURL = append(allURL, element)
		}
	}
	if len(allURL) == 0 {
		c.Response().WriteHeader(http.StatusNoContent)

		return nil
	}
	allURLJSON, err := json.Marshal(allURL)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("marshal error")
	}

	c.Response().Header().Add("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	c.Response().Write(allURLJSON)
	return nil
}
