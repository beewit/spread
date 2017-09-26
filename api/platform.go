package api

import (
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils"
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils/convert"
	"fmt"
)

func GetPlatformList(c echo.Context) error {
	r, err := ApiPost("/api/platform", nil)
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "请求失败,"+err.Error())
	}
	return utils.ResultApi(c, r)
}

func GetPlatformOne(t string) (map[string]interface{}, error) {
	url := fmt.Sprintf("/api/platform/one?type=%s", t)
	r, err := ApiPost(url, nil)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}

	data, err2 := convert.Obj2Map(r.Data)
	if err2 != nil {
		global.Log.Error(err2.Error())
		return nil, err2
	}
	return data, nil
}
