package router

import (
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/handler"

	"fmt"

	"github.com/labstack/echo"
)

func Router() {
	e := echo.New()
	e.Static("/app", "app")
	e.File("/", "app/page/index.html")
	handlerConfig(e)
	go e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", global.Port)))

}

func handlerConfig(e *echo.Echo) {
	e.POST("/auth/identity", handler.Identity)
}
