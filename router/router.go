package router

import (
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/handler"

	"fmt"

	"os"

	"github.com/beewit/spread/api"
	//"github.com/beewit/spread/static"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var e *echo.Echo

func Router() {
	defer func() {
		if err := recover(); err != nil {
			global.Log.Error("《程序出现严重错误，终止运行！》，ERROR：%v", err)
		}
	}()
	e = echo.New()
	file, _ := os.OpenFile("web.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	e.Logger.SetOutput(file)

	e.Static("/app", "app")
	e.File("/", "app/page/index.html")

	//e.GET("/*", echo.WrapHandler(static.Handler))
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	handlerConfig()
	//go e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", global.Port)))

	e.Start(fmt.Sprintf(":%v", global.Port))
}

func handlerConfig() {
	e.GET("/", handler.Index)

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
	e.POST("/wechat/send/message", handler.SendWechatMsg, handler.Filter)
	e.POST("/wechat/group/get/sendStatus", handler.GetSendWechatMsgStatus, handler.Filter)
	e.POST("/wechat/funcStatus", handler.GetWechatFuncStatus, handler.Filter)

	e.POST("/wechat/login", handler.LoginWechat, handler.Filter)
	e.POST("/wechat/login/check", handler.LoginWechatCheck, handler.Filter)
	e.POST("/wechat/add/user", handler.AddWechatUser, handler.Filter)
	e.POST("/wechat/cancel/login", handler.CancelLoginWechat, handler.Filter)

	e.POST("/wechat/client/list", handler.GetWechatClientList, handler.Filter)
	e.POST("/wechat/list/accountStatus", handler.GetWechatAccountStatus, handler.Filter)
	e.POST("/wechat/list/loginAccount", handler.LoginWechatListAccount, handler.Filter)
	e.POST("/wechat/list/sendMsg", handler.SendWechatListMsg, handler.Filter)
	e.POST("/wechat/list/addUser", handler.AddWechatListUser, handler.Filter)

	e.POST("/qq/login", handler.QQLogin, handler.Filter)
	e.POST("/qq/funcStatus", handler.GetQQFuncStatus, handler.Filter)
	e.POST("/qq/cancel/login", handler.CancelLoginQQ, handler.Filter)
	e.POST("/qq/login/check", handler.LoginQQCheck, handler.Filter)
	e.POST("/qq/status", handler.GetQQStatus, handler.Filter)
	e.POST("/qq/send/message", handler.SendQQMessage, handler.Filter)
	e.POST("/qq/search/group", handler.SearchQQGroup, handler.Filter)
	e.POST("/qq/add/group", handler.AddQQGroup, handler.Filter)
	e.POST("/qq/add/qq", handler.AddQQ, handler.Filter)
	e.POST("/qq/group/members", handler.GetQQGroupMembers, handler.Filter)
	e.POST("/qq/group/one/members", handler.GetQQGroupMembersByQQ, handler.Filter)
	e.POST("/qq/account/save", handler.SaveQQAccount, handler.Filter)
	e.POST("/qq/account/get", handler.GetQQAccount, handler.Filter)
	e.POST("/qq/client/list", handler.GetQQClientList, handler.Filter)

	e.GET("/task.js", handler.GetTask, handler.Filter)
	e.GET("/task/stop.js", handler.StopTask, handler.Filter)

	e.GET("/ReceiveToken", handler.ReceiveToken)
	e.GET("/signOut", handler.SignOut)
}

func Stop() {
	if e != nil {
		e.Close()
	}
}
