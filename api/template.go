package api

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/spread/global"
	"github.com/labstack/echo"
)

func GetTemplateByList(c echo.Context) error {
	body, err := uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    global.API_SERVICE_DOMAN + "/api/template",
	})
	if err != nil {
		return utils.Error(c, "请求失败,"+err.Error(), nil)
	}
	return utils.ResultApi(c, utils.ToResultParam(body))
}

func UpdateTemplateRefterById(c echo.Context) error {
	id := c.Param("id")
	body, err := uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    global.API_SERVICE_DOMAN + "/api/template/update/refer/" + id,
	})
	if err != nil {
		return utils.Error(c, "请求失败,"+err.Error(), nil)
	}
	return utils.ResultApi(c, utils.ToResultParam(body))
}

func GetTemplateById(c echo.Context) error {
	id := c.Param("id")
	body, err := uhttp.PostForm(global.API_SERVICE_DOMAN+"/api/template/"+id, nil)
	if err != nil {
		return utils.Error(c, "请求失败,"+err.Error(), nil)
	}
	return utils.ResultApi(c, utils.ToResultParam(body))
}
