package router

import (
	"ServiceShortURL/internal/shorturlservice"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

type userURLstruct struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s *Server) APIUserURL(c echo.Context) error {

	userCookies := shorturlservice.GetStructCookies()
	fmt.Println("<====API==1=======>", userCookies)

	allURL := make([]userURLstruct, 0)

	consumerURL, err := shorturlservice.NewConsumer(s.Cfg.Storage)
	if err != nil {
		log.Fatal(err)
	}
	defer consumerURL.Close()

	for {
		readURL, err := consumerURL.ReadURLInfo()
		if err != nil {
			break
		}
		fmt.Println(">>>>>>Api>0>>>", readURL)
		if readURL.CookiesAuthentication.ValueUser == userCookies.ValueUser {
			fmt.Println(">>>>>>Api>1>>>", readURL.CookiesAuthentication.ValueUser)
			fmt.Println(">>>>>>Api>2>>>", userCookies.ValueUser)
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
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("marshal error")
	}

	c.Response().Header().Add("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	c.Response().Write(allURLJSON)
	return nil
}
