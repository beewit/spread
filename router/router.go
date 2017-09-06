package router

import (
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/handler"

	"fmt"

	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/beewit/spread/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Router() {
	e := echo.New()

	assetHandler := http.FileServer(rice.MustFindBox("app").HTTPBox())
	e.GET("/app", echo.WrapHandler(assetHandler))
	e.GET("/app/*", echo.WrapHandler(http.StripPrefix("/app/", assetHandler)))
	e.Use(middleware.Gzip())
	//e.Use(middleware.Logger())
	//e.Use(middleware.CSRF())
	//e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	//e.Static("/app", "app")
	e.File("/", "app/page/index.html")
	handlerConfig(e)
	go e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", global.Port)))

}

func handlerConfig(e *echo.Echo) {
	e.POST("/auth/identity", handler.Identity)

	e.POST("/push/push", handler.Push)

	e.GET("/api/template", api.GetTemplateByList)

	e.GET("/api/template/:id", api.GetTemplateById)

	e.GET("/ReceiveToken", handler.ReceiveToken)

	e.POST("/platform/list", handler.PlatformList)
	e.POST("/platform/bind", handler.PlatformUnionBind)
}
