package api

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
	"fmt"
)

func GetWechatGroupList(c echo.Context) error {
	r, err := GetWechatGroupListData(c.FormValue("pageIndex"), c.FormValue("area"), c.FormValue("type"))
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "请求失败,"+err.Error())
	}
	return utils.ResultApi(c, r)
}

func GetWechatGroupClass(c echo.Context) error {
	url := "/api/wechat/group/class"
	r, err := ApiPost(url, nil)
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "请求失败,"+err.Error())
	}
	return utils.ResultApi(c, r)
}

func GetWechatGroupListData(pageIndex, area, types string) (*utils.ResultParam, error) {
	url := fmt.Sprintf("/api/wechat/group/list?pageIndex=%s&area=%s&type=%s", pageIndex, area, types)
	r, err := ApiPost(url, nil)
	if err != nil {
		global.Log.Error(err.Error())
		return nil, err
	}
	return &r, nil
}
