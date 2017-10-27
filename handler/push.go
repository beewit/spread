package handler

import (
	"fmt"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread/api"
	"github.com/beewit/spread/dao"
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/parser"
	"github.com/labstack/echo"
	"strings"
	"time"
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

	task := global.GetTask(global.TASK_PLATFORM_PUSH)
	if task != nil && task.State {
		return utils.ErrorNull(c, "任务正在运行中，请勿重复执行")
	}

	go func() {

		defer func() {
			global.DelTask(global.TASK_PLATFORM_PUSH)
		}()

		var names []string
		done := 0
		if m != nil && len(m) > 0 {
			for j := 0; j < len(m); j++ {
				platformName := convert.ToString(m[j]["platform_name"])
				names = append(names, platformName)
			}
			str := fmt.Sprintf("准备开始发送,发送平台：%s，发送标题：%s", strings.Join(names, ","), title)
			global.UpdateTask(global.TASK_PLATFORM_PUSH, str)

			//多帐号同时操作的时候进行切换帐号判断
			oldPlatformAcc := map[int64]string{}
			for j := 0; j < len(m); j++ {
				platformName := convert.ToString(m[j]["platform_name"])
				platformId := convert.MustInt64(m[j]["platform_id"])

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

					//任务记录
					task = global.GetTask(global.TASK_PLATFORM_PUSH)
					if task == nil || !task.State {
						str = fmt.Sprintf("【%s】任务已取消", global.TaskNameMap[global.TASK_PLATFORM_PUSH])
						global.Log.Info(str)
						global.PageSuccessMsg(str, global.Host+"?lastUrl=/app/page/admin/index.html")
						return
					}
					//更新任务记录
					global.UpdateTask(global.TASK_PLATFORM_PUSH, fmt.Sprintf("账号【%s】正在发送：%s", platformAcc, m[j]["name"]))

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

					flog, completFlog, result, resultMap, err := parser.RunPush(rule, paramMap, platformAcc, platformId, switchAccount)
					var logs string
					if err != nil {
						logs = fmt.Sprintf("《%s》 - > 异常：%v，result:%s", platformName, err.Error(), result)
						global.Log.Error(logs)
					} else {
						logs = fmt.Sprintf("《%s》 - > 状态：%v，result:%s", platformName, flog, result)
						global.Log.Error(logs)
					}
					if completFlog {
						done++
						global.PageSuccessMsg(platformName+" - > 发送结果成功", "")

						if completFlog {
							//执行成功数据
							if resultMap != nil {
								global.Log.Info("成功后的日志记录")
								resultMap["status"] = 1
							}
						} else {
							//执行失败数据
							global.Log.Info("失败后的日志记录")
						}
						if resultMap == nil {
							resultMap = map[string]interface{}{}
							resultMap["status"] = 0
						}
						iw, _ := utils.NewIdWorker(1)
						id, _ := iw.NextId()
						resultMap["id"] = id
						resultMap["type"] = 0
						resultMap["func_id"] = m[j]["id"]
						resultMap["func_name"] = m[j]["name"]
						resultMap["platform_id"] = platformId
						resultMap["platform_name"] = platformName
						resultMap["ct_time"] = utils.CurrentTime()
						resultMap["title"] = title
						resultMap["content"] = content
						resultMap["logs"] = logs
						resultMap["account_union_id"] = list[i]["id"]
						resultMap["account_id"] = global.Acc.Id
						flog, err := dao.SetFuncLogs(resultMap)
						if err != nil {
							logs = fmt.Sprintf("《%s》 - > 保存日志异常：%v ", platformName, err.Error())
							global.Log.Error(logs)
						} else {
							if flog {
								logs = fmt.Sprintf("《%s》 - > 保存日志成功", platformName)
							} else {
								logs = fmt.Sprintf("《%s》 - > 保存日志失败", platformName)
							}
							global.Log.Info(logs)
						}
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

func PushComm(title string, content string, rule string) (bool, bool, string, map[string]interface{}, error) {
	//println("Title：", title, "，Content:", content, "，规则：", rule)

	paramMap := map[string]string{
		"loginName": "登陆帐号",
		"loginPwd":  "登陆密码",
		"title":     title,
		"content":   content}
	return parser.RunPush(rule, paramMap, "", 1, false)
}
