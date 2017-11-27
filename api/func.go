package api

import (
	"fmt"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread/global"
	"github.com/labstack/echo"
)

func GetFuncByList(c echo.Context) error {
	r, err := getFuncByList(convert.ToString(global.Acc.Id), c.FormValue("pageIndex"))
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "请求失败,"+err.Error())
	}
	if r.Data != nil {

		group, _ := GetFuncGroup()

		nd, _ := convert.Obj2Map(r.Data)
		nd["group"] = group
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

func GetFuncByPlatformIdsAndAccId(platformIds, accId string) (utils.ResultParam, error) {
	url := fmt.Sprintf("/api/func/account/list?platformIds=%s&accId=%s", platformIds, accId)
	return ApiPost(url, nil)
}

func GetEffectiveFuncById(funcId string) (utils.ResultParam, error) {
	url := fmt.Sprintf("/api/func/account/funcId?funcId=%s", funcId)
	return ApiPost(url, nil)
}

func EffectiveFuncById(funcId int64) bool {
	r, err := GetEffectiveFuncById(convert.ToString(funcId))
	if err != nil {
		global.Log.Error("EffectiveFuncById ERROR %s", err.Error())
		return false
	}
	return r.Ret == utils.SUCCESS_CODE
}

func GetFuncGroup() (utils.ResultParam, error) {
	url := fmt.Sprintf("/api/func/account/group")
	return ApiPost(url, nil)
}
