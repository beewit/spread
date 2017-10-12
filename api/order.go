package api

import (
	"github.com/labstack/echo"
	"fmt"
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
)

func GetOrderByList(c echo.Context) error {
	r, err := getOrderByList(convert.ToString(global.Acc.Id), c.FormValue("pageIndex"))
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "请求失败,"+err.Error())
	}
	return utils.ResultApi(c, r)
}

func getOrderByList(accId, pageIndex string) (utils.ResultParam, error) {
	url := fmt.Sprintf("/api/order/pay/list?accId=%s&pageIndex=%s", accId, pageIndex)
	return ApiPost(url, nil)
}
