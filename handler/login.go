package handler

import (
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils"
	"github.com/labstack/echo"
	"github.com/beewit/spread/api"
	"github.com/beewit/beekit/log"
)

func SetClientToken(token string) *global.Account {
	acc := api.CheckClientLogin(token)
	if acc != nil {
		//insert sqllite
		flog, err := global.InsertToken(token)
		if err != nil {
			global.Log.Error(err.Error())
			return nil
		}
		if flog {
			return acc
		}
		log.Logger.Error("Token 写入本地数据库失败")
		return nil
	}
	return nil
}

func ReceiveToken(c echo.Context) error {
	token := c.FormValue("token")
	if token != "" {
		acc := SetClientToken(token)
		if acc != nil {
			return utils.Redirect(c, global.Host)
		}
	}
	return utils.RedirectAndAlert(c, "登陆失败", global.API_SSO_DOMAN)
}
