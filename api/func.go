package api

import (
	"github.com/labstack/echo"
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"fmt"
)

func GetFuncByList(c echo.Context) error {
	r, err := getFuncByList(convert.ToString(global.Acc.Id), c.FormValue("pageIndex"))
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "请求失败,"+err.Error())
	}
	if r.Data != nil {
		nd, _ := convert.Obj2Map(r.Data)
		if nd != nil {
			delete(nd, "rule")
			r.Data = nd
		}
	}
	return utils.ResultApi(c, r)
}

func GetAccountFuncList(c echo.Context) error {
	r, err := getAccountFuncList(c.FormValue("pageIndex"))
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "请求失败,"+err.Error())
	}
	if r.Data != nil {
		nd, _ := convert.Obj2Map(r.Data)
		if nd != nil {
			delete(nd, "rule")
			r.Data = nd
		}
	}
	return utils.ResultApi(c, r)
}

func getFuncByList(accId, pageIndex string) (utils.ResultParam, error) {
	url := fmt.Sprintf("/api/func/list?accId=%s&pageIndex=%s", accId, pageIndex)
	return ApiPost(url, nil)
}

func getAccountFuncList(pageIndex string) (utils.ResultParam, error) {
	url := fmt.Sprintf("/api/account/func/list?accId=%s&pageIndex=%s&pageSize=20",
		convert.ToString(global.Acc.Id), pageIndex)
	return ApiPost(url, nil)
}

func GetFuncAllByIdsAndAccId(funcIds, accId string) (utils.ResultParam, error) {
	url := fmt.Sprintf("/api/func/account/list?funcIds=%s&accId=%s", funcIds, accId)
	return ApiPost(url, nil)
}
