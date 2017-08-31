package handler

import (
	"github.com/labstack/echo"
	"hive/hive-sso/utils"
	"github.com/beewit/spread/global"
)

func PlatformList(c echo.Context) error {
	return utils.Success(c, "", global.Platform)
}

func PlatformUnionBind(c echo.Context) error {
	platform := c.FormValue("platform")
	if platform == "" {
		return utils.Error(c, "请选择平台进行绑定", nil)
	}
	pv := global.Platform[platform]
	if pv < 0 {
		return utils.Error(c, "请正确选择平台进行绑定", nil)
	}
	//获取远程服务中的平台信息进行帐号绑定操作

	return utils.Success(c, "正在前往绑定中...", "")
}
