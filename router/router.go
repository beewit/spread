package router

import (
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/handler"

	"fmt"

	"github.com/beewit/spread/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Router() {
	e := echo.New()

	e.Static("/app", "app")
	e.File("/", "app/page/index.html")
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	handlerConfig(e)
	go e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", global.Port)))

}

func handlerConfig(e *echo.Echo) {
	e.POST("/push/push", handler.Push, handler.Filter)

	e.POST("/api/template", api.GetTemplateByList, handler.Filter)
	e.POST("/api/template/:id", api.GetTemplateById, handler.Filter)
	e.POST("/api/template/update/refer/:id", api.UpdateTemplateReferById, handler.Filter)
	e.POST("/api/platform/list", api.GetPlatformList, handler.Filter)
	e.POST("/api/func/list", api.GetFuncByList, handler.Filter)
	e.POST("/api/account/func/list", api.GetAccountFuncList, handler.Filter)
	e.POST("/api/account/updatePwd", api.UpdatePwd, handler.Filter)

	e.POST("/api/rules/list", api.GetRules, handler.Filter)
	e.POST("/api/order/pay/list", api.GetOrderByList, handler.Filter)

	e.POST("/platform/bind", handler.PlatformUnionBind, handler.Filter)
	e.POST("/platform/union/list", handler.UnionList, handler.Filter)
	e.POST("/platform/union/login", handler.UnionLogin, handler.Filter)

	e.POST("/member/info", handler.GetMemberInfo, handler.Filter)
	e.GET("/ReceiveToken", handler.ReceiveToken)
	e.GET("/signOut", handler.SignOut)
}
