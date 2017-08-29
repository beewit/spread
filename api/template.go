package api

import (
	"encoding/json"

	"github.com/beewit/beekit/log"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/spread/global"
	"github.com/labstack/echo"
)

func GetTemplateByList(c echo.Context) error {
	body, err := uhttp.Cmd(uhttp.Request{
		Method: "GET",
		URL:    global.API_SERVICE_DOMAN + "/api/template",
	})
	if err != nil {
		return utils.Error(c, "请求失败,"+err.Error(), nil)
	}
	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err == nil {
		for k, v := range res {
			if k == "data" {
				data, err2 := json.Marshal(v)
				if err2 != nil {
					panic(err2)
				}
				var result []map[string]string
				println(string(data[:]))
				err2 = json.Unmarshal(data[:], &result)
				if err2 != nil {
					log.Logger.Error(err2.Error())
				}
				res[k] = &result
			}
		}
		return utils.Success(c, "获取数据成功", res)
	} else {
		return utils.Error(c, "获取数据失败,"+err.Error(), nil)
	}
}

func GetTemplateById(c echo.Context) error {
	println("--------------------")
	id := c.Param("id")
	body, err := uhttp.PostForm(global.API_SERVICE_DOMAN+"/api/template/"+id, nil)
	if err != nil {
		return utils.Error(c, "请求失败,"+err.Error(), nil)
	}
	var result map[string]string
	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err == nil {
		for k, v := range res {
			if k == "data" {
				data, err2 := json.Marshal(v)
				if err2 != nil {
					panic(err2)
				}
				println(string(data[:]))
				err2 = json.Unmarshal(data[:], &result)
				if err2 != nil {
					log.Logger.Error(err2.Error())
				}
				res[k] = &result
			}
		}
		return utils.Success(c, "获取数据成功", &result)
	} else {
		return utils.Error(c, "获取数据失败,"+err.Error(), nil)
	}
}
