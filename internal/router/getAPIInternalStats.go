package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"

	"ServiceShortURL/internal/shorturlservice"
)

type statsServer struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}

// GetAPIInternalStats e.GET("/api/internal/stats")
// возвращает пользователю количество пользователей и сокращенных ссылок
// {"urls": <int>, "users": <int>}
func (s *ServerShortener) GetAPIInternalStats(c echo.Context) error {
	s.WG.Wait()
	stats := statsServer{}

	if s.Cfg.Storage == testStorageURL {
		storageUsers = testStorageUsers
	}

	consumerURL, err := shorturlservice.NewConsumer(storageUsers)
	if err != nil {
		log.Fatal(err)
	}
	defer consumerURL.Close()

	//считаем пользователей из файла
	for readUsers, errRead := consumerURL.ReadUser(); errRead == nil; readUsers, errRead = consumerURL.ReadUser() {
		mapUsers := make(map[string]bool)
		user := readUsers.ValueUser
		//если такого пользователя нет в мапе, то добавляем в мапу пользователя и увеличиваем счётчик пользователей
		if !mapUsers[user] {
			mapUsers[user] = true
			stats.Users++
		}
	}

	//считаем ссылки из бд
	stats.Urls, err = s.DB.GetCount(c.Request().Context())
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	statsJSON, err := json.Marshal(stats)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("marshal error")
	}

	c.Response().Header().Add("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	c.Response().Write(statsJSON)
	return nil
}
