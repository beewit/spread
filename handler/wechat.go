package handler

import (
	"encoding/json"
	"math"
	"time"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	sapi "github.com/beewit/spread/api"
	"github.com/beewit/spread/global"
	"github.com/beewit/wechat-ai/api"
	"github.com/beewit/wechat-ai/enum"
	"github.com/labstack/echo"
	"github.com/sclevine/agouti"
)

var page *agouti.Page

func StartAddWechatGroup(c echo.Context) error {
	go addWechat(c, c.FormValue("pageIndex"), c.FormValue("area"), c.FormValue("type"))
	return utils.SuccessNull(c, "准备执行添加微信中..")
}

func addWechat(c echo.Context, pageIndex, area, types string) {
	r, err := sapi.GetWechatGroupListData(pageIndex, area, types)
	if err != nil {
		global.PageMsg("获取微信群信息失败")
		return
	}
	if r.Ret != utils.SUCCESS_CODE {
		global.PageMsg(r.Msg)
		return
	}
	bt, err := json.Marshal(r.Data)
	if err != nil {
		global.PageMsg("获取微信群信息失败")
		return
	}
	var pageData *utils.PageData
	json.Unmarshal(bt, &pageData)
	if err != nil {
		global.PageMsg("Failed to open page." + err.Error())
		return
	}
	if pageData == nil || pageData.Count <= 0 {
		global.PageMsg("暂无微信群信息")
		return
	}
	global.Navigate("http://www.baidu.com")

	for i := 0; i < len(pageData.Data); i++ {
		println(1)
		global.Navigate(convert.ToString(pageData.Data[i]["url"]))
		global.Page.RunScript(`$(".checkCode span:eq(1)").mouseover()`, nil, nil)
		time.Sleep(time.Second * 3)
		var of *enum.Offset
		println(2)
		global.Page.RunScript(`return  $(".shiftcode:eq(1) img").offset()`, nil, &of)
		title, err := global.Page.Title()
		if err != nil {

			println(3)
			global.Log.Error("error:" + err.Error())
			continue
		}

		println(4)
		if of != nil {

			println(5)
			err = api.Wechat(title, of)
			if err != nil {

				println(6)
				global.PageMsg(err.Error() + "，已经停止添加微信群")
				return
			}
		}
	}
	if pageData.PageIndex == int(math.Ceil(float64(pageData.Count)/float64(pageData.PageSize))) {
		//global.Page.CloseWindow()
		//global.Page.Destroy()
		global.PageMsg("添加数据完成")
		return
	} else {
		addWechat(c, convert.ToString(pageData.PageIndex+1), area, types)
	}
}
