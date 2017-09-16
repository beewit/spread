package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/api"
	"github.com/beewit/spread/dao"
	"github.com/beewit/spread/global"
	"github.com/labstack/echo"
)

func CheckClientLogin() *global.Account {
	token, err := dao.QueryLoginToken()
	if err != nil {
		global.Log.Error(err.Error())
		panic(err)
	}
	if token == "" {
		return nil
	}
	return api.CheckClientLogin(token)
}

func Filter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		acc := CheckClientLogin()
		if acc == nil {
			go global.Navigate(global.API_SSO_DOMAN + "?backUrl=" + global.Host + "/ReceiveToken")
			return utils.AuthFail(c, "登陆信息已失效，请重新登陆")
		}
		return next(c)
	}
}
