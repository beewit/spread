package api

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
	"github.com/labstack/echo"
)

func GetRules(c echo.Context) error {
	r, err := ApiPost("/api/rules/list", nil)
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "请求失败,"+err.Error())
	}
	return utils.ResultApi(c, r)
}
