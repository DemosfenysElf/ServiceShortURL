package router

import (
	"github.com/labstack/echo"
)

func Router() error {

	e := echo.New()

	e.GET("/:id", GetShortToURL)
	e.POST("/", PostURLToShort)
	e.POST("/api/shorten", APIShorten)

	err := e.Start(":8080")
	if err != nil {
		return err
	}
	return nil
}
