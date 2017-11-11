package handler

import (
	"fmt"
	"time"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/spread/api"
	"github.com/beewit/spread/dao"
	"github.com/beewit/spread/global"
	"github.com/beewit/wechat-ai/ai"
	"github.com/beewit/wechat-ai/smartQQ"
	"github.com/labstack/echo"
	"strings"
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
	if global.QQClient != nil && global.QQClient.GroupInfoMap != nil {
		for _, v := range global.QQClient.GroupInfoMap {
			_, err := global.QQClient.GetGroupInfo(v.Code)
			if err != nil {
				global.Log.Error("【%s】加载群信息失败，ERROR：%s", v.Name, err.Error())
			} else {
				global.Log.Info("【%s】加载群信息完成", v.Name)
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
			//var message string
			//if result.Result[0].PollType == "group_message" {
			//	group := qq.GroupInfoMap[result.Result[0].Value.GroupCode]
			//	if group.GId > 0 {
			//		message = " 【群消息 - " + group.Name + "】 "
			//	}
			//}
			//sendUser := qq.FriendsMap.Info[result.Result[0].Value.SendUin]
			//if sendUser.Uin > 0 {
			//	message += "   -   发送人《" + qq.FriendsMap.Info[result.Result[0].Value.SendUin].Nick + "》"
			//}
			//for i := 0; i < len(result.Result[0].Value.Content); i++ {
			//	if i > 0 {
			//		message += convert.ToObjStr(result.Result[0].Value.Content[i])
			//	}
			//}
			//global.Log.Info("您有新消息了哦！ ==>> ", message)
		}
	})
	if err != nil {
		global.Log.Error("Poll2 , ERROR：", err.Error())
		return
	}
	global.Log.Info("QQClient -->Poll2 , Info：%s", convert.ToObjStr(pollResult))
}

func GetQQGroupMembers(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_QQ) {
		return utils.ErrorNull(c, "QQ营销功能还未开通，请开通此功能后使用")
	}
	if global.QQClient == nil || !global.QQClient.Login.Status {
		return utils.ErrorNull(c, "未登录，请重新扫码登录后发送QQ消息")
	}
	qq := c.FormValue("qq")
	if qq == "" || !utils.IsValidNumber(qq) {
		return utils.ErrorNull(c, "群QQ错误")
	}
	if global.QQClient.Group2Map == nil {
		global.QQClient.GetMyGroupList()
	}
	if global.QQClient.Group2Map == nil {
		return utils.ErrorNull(c, "加载QQ群失败")
	}
	v := global.QQClient.Group2Map[convert.MustInt64(qq)]
	if v.QQ > 0 {
		global.QQClient.GetGroupMembers(v)
		return utils.SuccessNull(c, "加载群信息成功")
	} else {
		return utils.ErrorNull(c, fmt.Sprintf("查询群QQ%s失败", qq))
	}
}

func GetQQGroupMembersByQQ(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_QQ) {
		return utils.ErrorNull(c, "QQ营销功能还未开通，请开通此功能后使用")
	}
	if global.QQClient == nil || !global.QQClient.Login.Status {
		return utils.ErrorNull(c, "未登录，请重新扫码登录后发送QQ消息")
	}
	qq := c.FormValue("qq")
	if qq == "" || !utils.IsValidNumber(qq) {
		return utils.ErrorNull(c, "群QQ错误")
	}
	if global.QQClient.GroupMembersMap == nil {
		return utils.ErrorNull(c, "请先加载QQ群成员")
	}
	v := global.QQClient.GroupMembersMap[convert.MustInt64(qq)]
	if len(v.Mems) > 0 {
		return utils.Success(c, "加载群成员信息成功", v)
	} else {
		return utils.ErrorNull(c, "加载群成员失败")
	}
}

func SendQQMessage(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_QQ) {
		return utils.ErrorNull(c, "QQ营销功能还未开通，请开通此功能后使用")
	}
	if global.QQClient == nil || !global.QQClient.Login.Status {
		return utils.ErrorNull(c, "未登录，请重新扫码登录后发送QQ消息")
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
		if global.QQClient.FriendsMap.Info != nil && global.QQClient.GroupInfoMap != nil {
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
			if global.QQClient.GroupInfoMap != nil {
				global.Log.Info("准备开始发送群消息")
				for _, v := range global.QQClient.GroupInfoMap {
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

func SearchQQGroup(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_QQ) {
		return utils.ErrorNull(c, "QQ营销功能还未开通，请开通此功能后使用")
	}
	if global.QQClient == nil || !global.QQClient.Login.Status {
		return utils.ErrorNull(c, "未登录，请重新扫码登录后操作")
	}
	task := global.GetTask(global.TASK_QQ_SEND_MESSAGE)
	if task != nil && task.State {
		return utils.ErrorNull(c, "正在发送中，请勿重复执行")
	}
	keyword, pageStr, cityStr := c.FormValue("keyword"), c.FormValue("page"), c.FormValue("city")
	if strings.Trim(keyword, "") == "" {
		return utils.ErrorNull(c, "请输入搜索群关键词")
	}
	page, city := 0, 0
	if pageStr != "" && utils.IsValidNumber(pageStr) {
		page = convert.MustInt(pageStr)
	}
	if cityStr != "" && utils.IsValidNumber(cityStr) {
		city = convert.MustInt(cityStr)
	}
	if page == 0 {
		global.QQClient.SearchGroup.Total = 0
		global.QQClient.SearchGroup.SearchGroupList = nil
		global.QQClient.SearchKeyWord = keyword
	}
	groupSearch, err := global.QQClient.GetGroupSearch(keyword, city, page)
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "获取QQ群失败，ERROR:"+err.Error())
	}
	if groupSearch.RetCode.RetCode != 0 {
		return utils.ErrorNull(c, "获取QQ群失败，请重新登录QQ操作")
	}
	global.QQClient.SearchGroup.Total = groupSearch.Total
	if len(groupSearch.SearchGroupList) > 0 {
		for _, v := range groupSearch.SearchGroupList {
			global.QQClient.SearchGroup.SearchGroupList = append(global.QQClient.SearchGroup.SearchGroupList, v)
		}
	}
	return utils.SuccessNullMsg(c, groupSearch)
}

func AddQQGroup(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_QQ) {
		return utils.ErrorNull(c, "QQ营销功能还未开通，请开通此功能后使用")
	}
	if global.QQClient == nil || !global.QQClient.Login.Status {
		return utils.ErrorNull(c, "未登录，请重新扫码登录后操作")
	}
	if len(global.QQClient.SearchGroup.SearchGroupList) <= 0 {
		return utils.ErrorNull(c, "请先搜索群关键词后再启动QQ加群")
	}

	task := global.GetTask(global.TASK_QQ_ADD_FRIEND)
	if task != nil && task.State {
		return utils.ErrorNull(c, "正在添加好友中，不能同时进行")
	}
	task = global.GetTask(global.TASK_QQ_ADD_GROUP)
	if task != nil && task.State {
		return utils.ErrorNull(c, "正在加群中，请勿重复执行")
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
	qq := c.FormValue("qq")
	pwd := c.FormValue("pwd")
	remark := c.FormValue("remark")
	if qq == "" || pwd == "" {
		return utils.ErrorNull(c, "请设置QQ账号密码")
	}
	go func() {
		defer func() {
			global.DelTask(global.TASK_QQ_ADD_GROUP)
		}()
		global.UpdateTask(global.TASK_QQ_ADD_GROUP, "准备开发添加QQ好友..")
		utils.Close("QQ")
		time.Sleep(time.Second * 3)
		println(qq, pwd)
		err := ai.QQLogin(convert.MustInt64(qq), pwd)
		if err != nil {
			global.Log.Error(err.Error())
			global.PageMsg(err.Error())
			return
		}
		//获取最新群列表，进行排除添加用
		groupList, err := global.QQClient.GetMyGroupList()
		//连续5次错误放弃加群
		errCount, count := 0, 0
		var str string
		for _, v := range global.QQClient.SearchGroup.SearchGroupList {
			if IsExistGroup(groupList, v.GId) {
				str = fmt.Sprintf(" QQ群【%v】%s已添加过了", v.GId, v.Name)
				global.Log.Info(str)
				continue
			}
			count++
			if count > groupCount {
				global.UpdateTask(global.TASK_QQ_ADD_GROUP, fmt.Sprintf("延迟【%v】秒后添加QQ群", sleepTime))
				time.Sleep(time.Duration(sleepTime) * time.Second)
				count = 0
			}
			task := global.GetTask(global.TASK_QQ_ADD_GROUP)
			if task == nil || !task.State {
				str = fmt.Sprintf("【%s】已取消了", global.TaskNameMap[global.TASK_QQ_ADD_GROUP])
				global.Log.Info(str)
				global.PageMsg(str)
				return
			}
			//更新任务记录
			global.UpdateTask(global.TASK_QQ_ADD_GROUP, fmt.Sprintf("正在添加QQ群【%v】%s", v.GId, v.Name))
			err = ai.AddQQGroup(v.GId, remark)
			if err != nil {
				errCount++
				str = fmt.Sprintf("添加QQ群【%v】%s失败，原因：%s", v.GId, v.Name, err.Error())
				global.Log.Error(str)
				global.PageErrorMsg(str, "")
			} else {
				errCount = 0
				str = fmt.Sprintf("添加QQ群【%v】%s成功", v.GId, v.Name)
				global.PageSuccessMsg(str, "")
			}

			global.UpdateTask(global.TASK_QQ_ADD_GROUP, str)
			time.Sleep(time.Second * time.Duration(utils.NewRandom().Number(1)))
			//连续错误5次停止发送
			if errCount > 5 {
				global.PageMsg("连续5次以上添加QQ群失败，终止添加，请稍后重试！")
			}
		}
	}()
	return utils.SuccessNull(c, "'正在启动添加群成员")
}

func IsExistGroup(groupList map[int64]smartQQ.Group2, qq int64) bool {
	for _, v2 := range groupList {
		if v2.QQ == qq {
			return true
		}
	}
	return false
}

func AddQQ(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_QQ) {
		return utils.ErrorNull(c, "QQ营销功能还未开通，请开通此功能后使用")
	}
	if global.QQClient == nil || !global.QQClient.Login.Status {
		return utils.ErrorNull(c, "未登录，请重新扫码登录后操作")
	}
	if len(global.QQClient.MemberMap) <= 0 {
		return utils.ErrorNull(c, "请先获取群成员后再启动QQ加好友")
	}
	task := global.GetTask(global.TASK_QQ_ADD_GROUP)
	if task != nil && task.State {
		return utils.ErrorNull(c, "正在加群中，不能同时进行")
	}
	task = global.GetTask(global.TASK_QQ_ADD_FRIEND)
	if task != nil && task.State {
		return utils.ErrorNull(c, "正在加好友中，请勿重复执行")
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
	gqq := c.FormValue("gqq")
	qq := c.FormValue("qq")
	pwd := c.FormValue("pwd")
	remark := c.FormValue("remark")
	if qq == "" || pwd == "" {
		return utils.ErrorNull(c, "请设置QQ账号密码")
	}
	go func() {
		defer func() {
			global.DelTask(global.TASK_QQ_ADD_FRIEND)
		}()
		global.UpdateTask(global.TASK_QQ_ADD_FRIEND, "准备开发添加QQ好友..")

		utils.Close("QQ")
		time.Sleep(time.Second * 3)
		err := ai.QQLogin(convert.MustInt64(qq), pwd)
		if err != nil {
			global.PageMsg(err.Error())
			global.Log.Error(err.Error())
			return
		}
		//连续5次错误放弃加群
		errCount, count := 0, 0
		var str string
		for _, v := range global.QQClient.MemberMap {
			if gqq != "" {
				//不是当前待加的QQ群
				if v.GroupQQ != convert.MustInt64(gqq) {
					continue
				}
			}
			count++
			if count > groupCount {
				global.UpdateTask(global.TASK_QQ_ADD_FRIEND, fmt.Sprintf("延迟【%v】秒后添加QQ好友", sleepTime))
				time.Sleep(time.Duration(sleepTime) * time.Second)
				count = 0
			}
			task := global.GetTask(global.TASK_QQ_ADD_FRIEND)
			if task == nil || !task.State {
				str = fmt.Sprintf("【%s】已取消了", global.TaskNameMap[global.TASK_QQ_ADD_FRIEND])
				global.Log.Info(str)
				global.PageMsg(str)
				return
			}
			//更新任务记录
			global.UpdateTask(global.TASK_QQ_ADD_FRIEND, fmt.Sprintf("正在添加QQ好友【%v】%s", v.QQ, v.Nick))
			err = ai.AddQQFriend(v.QQ, remark)
			if err != nil {
				errCount++
				str = fmt.Sprintf("添加QQ好友【%v】%s失败，原因：%s", v.QQ, v.Nick, err.Error())
				global.Log.Error(str)
				global.PageErrorMsg(str, "")
			} else {
				errCount = 0
				str = fmt.Sprintf("添加QQ好友【%v】%s成功", v.QQ, v.Nick)
				global.PageSuccessMsg(str, "")
			}

			global.UpdateTask(global.TASK_QQ_ADD_FRIEND, str)
			time.Sleep(time.Second * time.Duration(utils.NewRandom().Number(1)))
			//连续错误5次停止发送
			if errCount > 5 {
				global.PageMsg("连续5次以上添加QQ好友失败，终止添加，请稍后重试！")
			}
		}
	}()
	return utils.SuccessNull(c, "'正在启动添加QQ好友")
}

func SaveQQAccount(c echo.Context) error {
	qq := c.FormValue("qq")
	pwd := c.FormValue("pwd")
	remark := c.FormValue("remark")
	if qq == "" || pwd == "" {
		return utils.ErrorNull(c, "请设置QQ账号密码")
	}
	flog, err := dao.SetUnion(enum.QQ, qq, pwd, remark, enum.QQ_ID, global.Acc.Id)
	if err != nil {
		str := fmt.Sprintf("保存QQ账号失败，原因：%s", err.Error())
		global.Log.Error(str)
		return utils.ErrorNull(c, str)
	} else {
		if flog {
			return utils.SuccessNull(c, "保存QQ账号成功！")
		} else {
			return utils.ErrorNull(c, "保存QQ账号失败！")
		}
	}
}

func GetQQAccount(c echo.Context) error {
	m, err := dao.GetUnionList(enum.QQ_ID, global.Acc.Id)
	if err != nil {
		str := fmt.Sprintf("查询QQ账号异常，原因：%s", err.Error())
		global.Log.Error(str)
		return utils.SuccessNull(c, "")
	} else {
		return utils.SuccessNullMsg(c, m)
	}
}

func UpdateFriend(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_QQ) {
		return utils.ErrorNull(c, "QQ营销功能还未开通，请开通此功能后使用")
	}
	if global.QQClient == nil || !global.QQClient.Login.Status {
		return utils.ErrorNull(c, "未登录，请重新扫码登录后操作")
	}
	_, err := global.QQClient.GetFriendList()
	if err != nil {
		str := fmt.Sprintf("更新好友错误，原因：%s", err.Error())
		global.Log.Error(str)
		return utils.ErrorNull(c, "更新失败")
	} else {
		return utils.SuccessNullMsg(c, "更新成功")
	}
}

func UpdateGroup(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_QQ) {
		return utils.ErrorNull(c, "QQ营销功能还未开通，请开通此功能后使用")
	}
	if global.QQClient == nil || !global.QQClient.Login.Status {
		return utils.ErrorNull(c, "未登录，请重新扫码登录后操作")
	}
	_, err := global.QQClient.GetMyGroupList()
	if err != nil {
		str := fmt.Sprintf("更新群错误，原因：%s", err.Error())
		global.Log.Error(str)
		return utils.ErrorNull(c, "更新失败")
	} else {
		return utils.SuccessNullMsg(c, "更新成功")
	}
}
