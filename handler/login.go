package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/api"
	"github.com/beewit/spread/dao"
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/static"
	"github.com/labstack/echo"
)

func Index(c echo.Context) error {
	return utils.ResultHtml(c, string(static.FileAppPageIndexHTML))
}

func SetClientToken(token string) *global.Account {
	acc := api.CheckClientLogin(token)
	if acc != nil {
		//insert sqllite
		flog, err := dao.InsertToken(token, acc)
		if err != nil {
			global.Log.Error(err.Error())
			return nil
		}
		if flog {
			return acc
		}
		global.Log.Error("Token 写入本地数据库失败")
		return nil
	}
	return nil
}

func ReceiveToken(c echo.Context) error {
	token := c.FormValue("token")
	if token != "" {
		acc := SetClientToken(token)
		if acc != nil {
			global.Acc = acc
			global.Acc.Token = token
			return utils.Redirect(c, global.Host)
		}
	}
	return utils.RedirectAndAlert(c, "登陆失败", global.API_SSO_DOMAIN+"?backUrl="+global.Host+"/ReceiveToken")
}

func SignOut(c echo.Context) error {
	api.DeleteToken(global.Acc.Token)
	dao.DeleteToken(global.Acc)
	global.Acc = nil
	return utils.RedirectAndAlert(c, "", global.API_SSO_DOMAIN+"?backUrl="+global.Host+"/ReceiveToken")
}
