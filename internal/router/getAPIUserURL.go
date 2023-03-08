package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

type userURLstruct struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// GetAPIUserURL e.GET("/api/user/urls")
// возвращает пользователю все когда-либо сокращённые им URL в формате
// [{"short_url":"http://...","original_url":"http://..."},...]
// при отсутствии сокращённых пользователем URL: http.StatusNoContent
func (s *serverShortener) GetAPIUserURL(c echo.Context) error {
	s.WG.Wait()
	fmt.Println("==>> GetAPIUserURL")

	userCookies := shorturlservice.GetCookieValue(c.Request().Cookies())

	allURL := make([]userURLstruct, 0)

	consumerURL, err := shorturlservice.NewConsumer(s.Cfg.Storage)
	if err != nil {
		log.Fatal(err)
	}
	defer consumerURL.Close()

	for readURL, err := consumerURL.ReadURLInfo(); err == nil; readURL, err = consumerURL.ReadURLInfo() {
		if readURL.CookiesAuthentication.ValueUser == userCookies {
			element := userURLstruct{
				ShortURL:    s.Cfg.BaseURL + "/" + readURL.ShortURL,
				OriginalURL: readURL.URL,
			}
			allURL = append(allURL, element)
		}
	}

	if len(allURL) == 0 {
		c.Response().WriteHeader(http.StatusNoContent)
		return nil
	}

	allURLJSON, err := json.Marshal(allURL)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("marshal error")
	}

	c.Response().Header().Add("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	c.Response().Write(allURLJSON)
	return nil
}
