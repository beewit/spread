package api

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
	"github.com/labstack/echo"
)

func GetTemplateByList(c echo.Context) error {
	r, err := ApiPost("/api/template", nil)
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "请求失败,"+err.Error())
	}
	return utils.ResultApi(c, r)
}

func UpdateTemplateReferById(c echo.Context) error {
	id := c.Param("id")
	r, err := ApiPost("/api/template/update/refer/"+id, nil)
	if err != nil {
		return utils.ErrorNull(c, "请求失败,"+err.Error())
	}
	return utils.ResultApi(c, r)
}

func GetTemplateById(c echo.Context) error {
	id := c.Param("id")
	r, err := ApiPost("/api/template/"+id, nil)
	if err != nil {
		return utils.ErrorNull(c, "请求失败,"+err.Error())
	}
	return utils.ResultApi(c, r)
}
