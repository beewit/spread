package router

import (
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/handler"

	"fmt"

	"github.com/beewit/spread/api"
	"github.com/beewit/spread/static"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io"
	"os"
)

var e *echo.Echo

type LoggerConfig struct {
	// 可选。默认值是 DefaultLoggerConfig.Format.
	Format string `json:"format"`
	// Output 是记录日志的位置。
	// 可选。默认值是 os.Stdout.
	Output io.Writer
}

func Router() {
	e = echo.New()
	file, _ := os.OpenFile("web.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	//	Output: file,
	//}))
	//e.Logger.SetLevel(log.OFF)
	e.Logger.SetOutput(file)

	//e.Static("/app", "app")
	e.GET("/*", echo.WrapHandler(static.Handler))
	e.File("/", "app/page/index.html")
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	handlerConfig()
	//go e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", global.Port)))
	go e.Start(fmt.Sprintf(":%v", global.Port))
}

func handlerConfig() {
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
	e.POST("/api/wechat/group/list", api.GetWechatGroupList, handler.Filter)
	e.POST("/api/wechat/group/class", api.GetWechatGroupClass, handler.Filter)

	e.POST("/platform/bind", handler.PlatformUnionBind, handler.Filter)
	e.POST("/platform/union/list", handler.UnionList, handler.Filter)
	e.POST("/platform/union/login", handler.UnionLogin, handler.Filter)
	e.POST("/platform/union/delete", handler.UnionDelete, handler.Filter)

	e.POST("/member/info", handler.GetMemberInfo, handler.Filter)
	e.POST("/member/bindWechat", handler.CreateWechatQrCode, handler.Filter)

	e.POST("/wechat/group/start/add", handler.StartAddWechatGroup, handler.Filter)
	e.POST("/wechat/group/start/send", handler.SendWechatMsg, handler.Filter)
	e.POST("/wechat/group/get/sendStatus", handler.GetSendWechatMsgStatus, handler.Filter)
	e.POST("/wechat/funcStatus", handler.GetWechatFuncStatus, handler.Filter)

	e.GET("/ReceiveToken", handler.ReceiveToken)
	e.GET("/signOut", handler.SignOut)
}

func Stop() {
	e.Close()
}
