package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread/api"
	"github.com/beewit/spread/dao"
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/parser"
	"github.com/labstack/echo"
	"strings"
	"time"
	"fmt"
)

func Push(c echo.Context) error {
	title := c.FormValue("title")
	content := c.FormValue("content")
	funIds := c.FormValue("funIds")
	rp, err := api.GetFuncAllByIdsAndAccId(funIds, convert.ToString(global.Acc.Id))
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "获取待发送的网站模块失败，err："+err.Error())
	}
	m, err := convert.Obj2ListMap(rp.Data)
	go func() {
		//println("开始执行简书分发")
		////   utils.JsonPath("parser", "./jianshu.json")
		//rule := utils.Read("./parser/jianshu.json")
		//flog, result, err2 := PushComm(title, content, rule)
		//println("简书分发", flog, result, err2)
		//
		//println("开始执行新浪分发")
		////utils.JsonPath("parser", "./sina.json")
		//rule = utils.Read("./parser/sina.json")
		//flog, result, err2 = PushComm(title, content, rule)
		//println("微博分发", flog, result, err2)
		//
		//println("开始执行知乎分发")
		////utils.JsonPath("parser", "./zhihu.json")
		//rule = utils.Read("./parser/zhihu.json")
		//flog, result, err2 = PushComm(title, content, rule)
		//println("知乎分发", flog, result, err2)
		var names []string
		done := 0
		if m != nil && len(m) > 0 {
			//多帐号同时操作的时候进行切换帐号判断
			oldPlatformAcc := map[int64]string{}
			for j := 0; j < len(m); j++ {
				platformName := convert.ToString(m[j]["platform_name"])
				platformId := convert.MustInt64(m[j]["platform_id"])
				names = append(names, platformName)

				list, err := dao.GetUnionList(platformId, global.Acc.Id)
				if err != nil {
					global.PageMsg("[" + platformName + "]查找平台绑定帐号失败")
					continue
				}
				if list == nil || len(list) <= 0 {
					global.PageMsg("[" + platformName + "]未绑定平台帐号，请进入《帐号》->《平台帐号绑定》->《新增平台帐号》 ->点击要绑定的平台帐号")
					continue
				}
				for i := 0; i < len(list); i++ {
					platformAcc := convert.ToString(list[i]["platform_account"])
					platformId := convert.MustInt64(list[i]["platform_id"])
					platformPwd := convert.ToString(list[i]["platform_password"])
					rule := convert.ToString(m[j]["rule"])
					paramMap := map[string]string{
						"loginName": convert.ToString(platformAcc),
						"loginPwd":  convert.ToString(platformPwd),
						"title":     title,
						"content":   content}
					switchAccount := false
					if oldPlatformAcc[platformId] != "" && oldPlatformAcc[platformId] != platformAcc {
						switchAccount = true
					}
					oldPlatformAcc[platformId] = platformAcc

					flog, completFlog, result, err := parser.RunPush(rule, paramMap, platformAcc, platformId, switchAccount)
					if err != nil {
						global.Log.Error(platformName, " - > 发送失败")
						global.Log.Error("《%s》 - > 异常：%v，result:%s", platformName, err.Error(), result)
					} else {
						global.Log.Error("《%s》 - > 状态：%v，result:%s", platformName, flog, result)
					}
					if completFlog {
						done++
						global.PageSuccessMsg(platformName+" - > 发送结果成功", "")
					} else {
						global.PageErrorMsg(platformName+" - > 发送结果失败", "")
					}
					time.Sleep(time.Second * 2)
				}
			}
		}
		tip := fmt.Sprintf("《%s》全部发布完成，发布成功：%v，发布失败：%v", strings.Join(names, ","), done, len(names)-done)
		global.Log.Info(tip)
		time.Sleep(time.Second * 2)
		if strings.Contains(global.PageUrl(), global.Host) {
			global.PageSuccessMsg(tip, "")
		} else {
			global.PageSuccessMsg(tip, global.Host+"?lastUrl=/app/page/admin/index.html")
		}
	}()

	return utils.Success(c, "正在发布中", "")
}

func PushComm(title string, content string, rule string) (bool, bool, string, error) {
	//println("Title：", title, "，Content:", content, "，规则：", rule)

	paramMap := map[string]string{
		"loginName": "18223277005",
		"loginPwd":  "13696433488wb",
		"title":     title,
		"content":   content}
	return parser.RunPush(rule, paramMap, "", 1, false)
}
