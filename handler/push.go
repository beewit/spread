package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/parser"
	"github.com/labstack/echo"
	"github.com/beewit/spread/dao"
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread/api"
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

		if m != nil && len(m) > 0 {
			for i := 0; i < len(m); i++ {
				platformName := convert.ToString(m[i]["platform_name"])
				list, err := dao.GetUnionList(convert.MustInt64(m[i]["platform_id"]), global.Acc.Id)
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
					rule := convert.ToString(m[i]["rule"])
					paramMap := map[string]string{
						"loginName": convert.ToString(platformAcc),
						"loginPwd":  convert.ToString(platformPwd),
						"title":     title,
						"content":   content}
					flog, rulest, err := parser.RunPush(rule, paramMap, platformAcc, platformId)
					if err != nil {
						global.Log.Error(platformName, " - > 发送失败", err.Error())
					} else {
						global.Log.Error(platformName, " - > 发送状态：", flog, "rulest:", rulest)
					}

				}
			}
		}
	}()

	return utils.Success(c, "正在发布中", "")
}

func PushComm(title string, content string, rule string) (bool, string, error) {
	//println("Title：", title, "，Content:", content, "，规则：", rule)

	paramMap := map[string]string{
		"loginName": "18223277005",
		"loginPwd":  "13696433488wb",
		"title":     title,
		"content":   content}
	return parser.RunPush(rule, paramMap, "", 1)
}
