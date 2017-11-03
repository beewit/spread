package handler

import (
	"fmt"
	"time"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread/api"
	"github.com/beewit/spread/global"
	"github.com/beewit/wechat-ai/smartQQ"
	"github.com/labstack/echo"
	"sync"
)

func GetQQFuncStatus(c echo.Context) error {
	flog := api.EffectiveFuncById(global.FUNC_QQ)
	return utils.SuccessNullMsg(c, flog)
}

func CancelLoginQQ(c echo.Context) error {
	global.QQClient.LoginCheck = false
	return utils.SuccessNull(c, "")
}

var LoginIng = false

func QQLogin(c echo.Context) error {
	if LoginIng {
		return utils.ErrorNull(c, "正在登录中，请勿重复点击登录")
	}
	_, err := global.QQClient.PtqrShow()
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "获取QQ登录二维码失败")
	}
	go func() {
		defer func() {
			LoginIng = false
		}()
		LoginIng = true
		global.QQClient, err = global.QQClient.CheckLogin(func(newQQ *smartQQ.QQClient, err error) {
			if newQQ.Login.Status {
				go LoadGroupInfo()
				go Pull()
			}
		})
		if err != nil {
			global.QQClient.Login.Desc = "登录失败，ERROR：" + err.Error()
			global.Log.Error(global.QQClient.Login.Desc)
		}
	}()
	return utils.Success(c, "扫描登录QQ网页服务", global.QQClient.LoginQrCode)
}

//加载群信息
func LoadGroupInfo() {
	if global.QQClient != nil && global.QQClient.GroupInfo != nil {
		for _, v := range global.QQClient.GroupInfo {
			reg, err := global.QQClient.GetGroupInfo(v.Code)
			if err != nil {
				global.Log.Error("【%s】加载群信息失败，ERROR：%s", v.Name, err.Error())
			} else {
				global.Log.Info("【%s】加载群信息结果：%s", v.Name, convert.ToObjStr(reg))
			}
			time.Sleep(time.Second * 3)
		}
	}
}

var smLoginQQCheck *sync.Mutex

func LoginQQCheck(c echo.Context) error {
	if global.QQClient == nil {
		return utils.SuccessNullMsg(c, nil)
	}
	if smLoginQQCheck == nil {
		smLoginQQCheck = new(sync.Mutex)
	}
	smLoginQQCheck.Lock()
	rep, err := global.QQClient.TestLogin()
	smLoginQQCheck.Unlock()
	if err == nil && rep.RetCode == 0 {
		return utils.SuccessNullMsg(c, map[string]interface{}{"QQUser": global.QQClient})
	}
	return utils.SuccessNullMsg(c, nil)
}

func GetQQStatus(c echo.Context) error {
	return utils.SuccessNullMsg(c, map[string]interface{}{"sendStatusMsg": global.QQClient.Login.Desc, "sendStatus": global.QQClient.Login.Status})
}

func Pull() {
	time.Sleep(time.Second * 5)
	pollResult, err := global.QQClient.Poll2(func(qq *smartQQ.QQClient, result smartQQ.QQResponsePoll) {
		if len(result.Result) > 0 && len(result.Result[0].Value.Content) > 0 {
			var message string
			if result.Result[0].PollType == "group_message" {
				group := qq.GroupInfo[result.Result[0].Value.GroupCode]
				if group.GId > 0 {
					message = " 【群消息 - " + group.Name + "】 "
				}
			}
			sendUser := qq.FriendsMap.Info[result.Result[0].Value.SendUin]
			if sendUser.Uin > 0 {
				message += "   -   发送人《" + qq.FriendsMap.Info[result.Result[0].Value.SendUin].Nick + "》"
			}
			for i := 0; i < len(result.Result[0].Value.Content); i++ {
				if i > 0 {
					message += convert.ToObjStr(result.Result[0].Value.Content[i])
				}
			}
			global.Log.Info("您有新消息了哦！ ==>> ", message)
		}
	})
	if err != nil {
		global.Log.Error("Poll2 , ERROR：", err.Error())
		return
	}
	global.Log.Info("QQClient -->Poll2 , Info：%s", convert.ToObjStr(pollResult))
}

func SendQQMessage(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_QQ) {
		return utils.ErrorNull(c, "QQ营销功能还未开通，请开通此功能后使用")
	}
	if global.QQClient == nil || !global.QQClient.Login.Status {
		return utils.ErrorNull(c, "未登录，请重新扫码登录后添加QQ好友")
	}
	task := global.GetTask(global.TASK_QQ_SEND_MESSAGE)
	if task != nil && task.State {
		return utils.ErrorNull(c, "正在发送中，请勿重复执行")
	}
	content := c.FormValue("msg")
	if content == "" {
		return utils.ErrorNull(c, "发送内容不能为空")
	}
	groupCountStr := c.FormValue("groupCount")
	if groupCountStr == "" || !utils.IsValidNumber(groupCountStr) {
		groupCountStr = "3"
	}
	sleepTimeStr := c.FormValue("sleepTime")
	if sleepTimeStr == "" || !utils.IsValidNumber(sleepTimeStr) {
		sleepTimeStr = "30"
	}

	groupCount := convert.MustInt(groupCountStr)
	sleepTime := convert.MustInt64(sleepTimeStr)

	global.Log.Info("QQ消息发送内容：%s", content)
	go func() {
		defer func() {
			global.DelTask(global.TASK_QQ_SEND_MESSAGE)
		}()
		if global.QQClient.FriendsMap.Info != nil && global.QQClient.GroupInfo != nil {
			global.UpdateTask(global.TASK_QQ_SEND_MESSAGE, "准备开发发送QQ消息..")
			global.QQClient.StatusMessage = "准备开发发送QQ消息！"
			global.Log.Info(global.QQClient.StatusMessage)
			time.Sleep(time.Second * 20)
			var sleep int
			var str string
			errCount := 0
			count := 0

			if global.QQClient.FriendsMap.Info != nil {
				global.Log.Info("准备开始发送好友消息")
				for _, v := range global.QQClient.FriendsMap.Info {
					count++
					if count > groupCount {
						global.UpdateTask(global.TASK_QQ_SEND_MESSAGE, fmt.Sprintf("延迟【%v】秒后发送QQ消息", sleepTime))
						time.Sleep(time.Duration(sleepTime) * time.Second)
						count = 0
					}
					task := global.GetTask(global.TASK_QQ_SEND_MESSAGE)
					if task == nil || !task.State {
						str = fmt.Sprintf("【%s】已取消了", global.TaskNameMap[global.TASK_QQ_SEND_MESSAGE])
						global.Log.Info(str)
						global.PageMsg(str)
						return
					}
					//更新任务记录
					global.UpdateTask(global.TASK_QQ_SEND_MESSAGE, fmt.Sprintf("正在发送QQ给用户【%s】", v.Nick))
					res, err := global.QQClient.SendMsg(v.Uin, content)
					if err != nil {
						errCount++
						global.QQClient.StatusMessage = fmt.Sprintf("发送消息发生错误，ERROR：%s", err.Error())
						sleep = utils.NewRandom().Number(2)
					} else {
						if res.RetCode == 0 {
							errCount = 0
							global.QQClient.StatusMessage = "发送给【" + v.Nick + "】成功"
							sleep = utils.NewRandom().Number(1)
						} else {
							global.QQClient.StatusMessage = "发送给【" + v.Nick + "】失败，"
							sleep = utils.NewRandom().Number(2)
						}
					}
					global.UpdateTask(global.TASK_QQ_SEND_MESSAGE, global.QQClient.StatusMessage)
					time.Sleep(time.Second)
					global.UpdateTask(global.TASK_QQ_SEND_MESSAGE, fmt.Sprintf("延迟【%v】秒后发送QQ消息", sleep))
					time.Sleep(time.Second * time.Duration(sleep))
					//连续错误5次停止发送
					if errCount > 5 {
						global.PageMsg("连续5次以上发送失败，终止发送，请稍后重试！")
					}
				}
			}
			count = 0
			if global.QQClient.GroupInfo != nil {
				global.Log.Info("准备开始发送群消息")
				for _, v := range global.QQClient.GroupInfo {
					count++
					if count > groupCount {
						global.UpdateTask(global.TASK_QQ_SEND_MESSAGE, fmt.Sprintf("延迟【%v】秒后发送QQ群消息", sleepTime))
						time.Sleep(time.Duration(sleepTime) * time.Second)
						count = 0
					}
					task := global.GetTask(global.TASK_QQ_SEND_MESSAGE)
					if task == nil || !task.State {
						str = fmt.Sprintf("【%s】已取消了", global.TaskNameMap[global.TASK_QQ_SEND_MESSAGE])
						global.Log.Info(str)
						global.PageMsg(str)
						return
					}
					//更新任务记录
					global.UpdateTask(global.TASK_QQ_SEND_MESSAGE, fmt.Sprintf("正在发送QQ群【%s】", v.Name))
					res, err := global.QQClient.SendQunMsg(v.GId, content)
					if err != nil {
						errCount++
						global.QQClient.StatusMessage = fmt.Sprintf("发送QQ群消息发生错误，ERROR：%s", err.Error())
						sleep = utils.NewRandom().Number(2)
					} else {
						if res.RetCode == 0 {
							errCount = 0
							global.QQClient.StatusMessage = "发送给QQ群【" + v.Name + "】成功"
							sleep = utils.NewRandom().Number(1)
						} else {
							global.QQClient.StatusMessage = "发送给QQ群【" + v.Name + "】失败，"
							sleep = utils.NewRandom().Number(2)
						}
					}
					global.UpdateTask(global.TASK_QQ_SEND_MESSAGE, global.QQClient.StatusMessage)
					time.Sleep(time.Second)
					global.UpdateTask(global.TASK_QQ_SEND_MESSAGE, fmt.Sprintf("延迟【%v】秒后发送QQ群消息", sleep))
					time.Sleep(time.Second * time.Duration(sleep))
					//连续错误5次停止发送
					if errCount > 5 {
						global.PageMsg("连续5次以上发送失败，终止发送，请稍后重试！")
					}
				}
			}

			global.QQClient.StatusMessage = "QQ发消息任务完成！"
			global.Log.Info(global.QQClient.StatusMessage)
			global.PageSuccessMsg(global.QQClient.StatusMessage, global.Host+"?lastUrl=/app/page/admin/qq/index.html")
		} else {
			global.PageMsg("好友列表未获取到！")
		}
	}()
	return utils.SuccessNull(c, "后台发送中...")
}
