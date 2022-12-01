package router

import (
	"github.com/labstack/echo"
)

func Router() error {

	e := echo.New()

	e.GET("/:id", getShortToURL)
	e.POST("/", postURLToShort)

	err := e.Start(":8080")
	if err != nil {
		return err
	}
	return nil
}
