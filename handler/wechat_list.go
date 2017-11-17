package handler

import (
	"errors"
	"fmt"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread/api"
	"github.com/beewit/spread/dao"
	"github.com/beewit/spread/global"
	"github.com/beewit/wechat-ai/enum"
	"github.com/beewit/wechat-ai/smartWechat"
	"github.com/labstack/echo"
	"strings"
	"time"
)

func GetWechatClientList(c echo.Context) error {
	return utils.SuccessNullMsg(c, global.WechatClientList)
}

func GetWechatAccountStatus(c echo.Context) error {
	uuid := c.FormValue("uuid")
	if uuid == "" {
		return utils.SuccessNullMsg(c, nil)
	}
	if global.WechatClientList[uuid] == nil {
		return utils.SuccessNullMsg(c, nil)
	}
	return utils.SuccessNullMsg(c, map[string]interface{}{
		"sendStatusMsg": global.WechatUUid[uuid].SendStatusMsg,
		"sendStatus":    global.WechatUUid[uuid].SendStatus})
}

func LoginWechatListAccount(c echo.Context) error {
	/* 从微信服务器获取UUID */
	UUid, err := smartWechat.GetUUIDFromWX()
	if err != nil {
		return utils.ErrorNull(c, "从微信服务器获取UUID失败")
	}
	/* 根据UUID获取二维码 */
	base64Img, err := smartWechat.DownloadImage(enum.QRCODE_URL + UUid)
	if err != nil {
		return utils.ErrorNull(c, "根据UUID获取微信登录二维码失败")
	}
	go LoginWechatList(UUid)
	return utils.Success(c, "扫描登录微信网页服务", map[string]string{"base64Img": base64Img, "UUid": UUid})
}

func LoginWechatList(UUid string) (err error) {
	WechatLoginCheck = true
	timeOut := 0
	for {
		//thisUrl, _ := global.Page.Page.URL()
		global.WechatUUid[UUid] = &smartWechat.WechatLoginStatus{}
		if WechatLoginCheck {
			global.WechatUUid[UUid].SendStatusMsg = "【" + UUid + "】正在验证登陆... ..."
			global.Log.Info(global.WechatUUid[UUid].SendStatusMsg)
			status, msg := smartWechat.CheckLogin(UUid)
			if status == 200 {
				global.WechatUUid[UUid].SendStatusMsg = "登陆成功,处理登陆信息..."
				global.Log.Info(global.WechatUUid[UUid].SendStatusMsg)
				var wc *smartWechat.WechatClient
				wc, err = smartWechat.ProcessLoginInfo(msg)
				if err != nil {
					global.WechatUUid[UUid].SendStatus = WECHAT_STATUS_FAIL
					global.WechatUUid[UUid].SendStatusMsg = "错误：登陆成功,处理登陆信息...，error：" + err.Error()
					global.Log.Info(global.WechatUUid[UUid].SendStatusMsg)
					return
				}
				selfUserName := wc.SelfUserName
				wc.UUID = UUid
				global.WechatClientList[selfUserName] = wc
				global.WechatUUid[UUid].SendStatusMsg = "登陆信息处理完毕,正在初始化微信..."
				global.Log.Info(global.WechatUUid[UUid].SendStatusMsg)
				global.Log.Info("global.WechatClientList[selfUserName]：%s", convert.ToObjStr(global.WechatClientList[selfUserName]))
				err = smartWechat.InitWX(global.WechatClientList[selfUserName])
				if err != nil {
					global.WechatUUid[UUid].SendStatusMsg = "【1】错误：登陆信息处理完毕,正在初始化微信...，error：" + err.Error()
					global.Log.Info(global.WechatUUid[UUid].SendStatusMsg)
					if global.WechatClientList[selfUserName].InitInfo.BaseResponse.Ret == 1010 {
						err = smartWechat.InitWX(global.WechatClientList[selfUserName])
					}
					if err != nil {
						global.WechatUUid[UUid].SendStatusMsg = "【2】错误：登陆信息处理完毕,正在初始化微信...，error：" + err.Error()
						global.WechatUUid[UUid].SendStatus = WECHAT_STATUS_FAIL
						global.Log.Info(global.WechatUUid[UUid].SendStatusMsg)
						return
					}
				}
				global.WechatUUid[UUid].SendStatusMsg = "初始化完毕,通知微信服务器登陆状态变更..."
				global.Log.Info(global.WechatUUid[UUid].SendStatusMsg)
				err = smartWechat.NotifyStatus(global.WechatClientList[selfUserName])
				if err != nil {
					global.WechatUUid[UUid].SendStatus = WECHAT_STATUS_FAIL
					global.WechatUUid[UUid].SendStatusMsg = "通知微信服务器状态变化失败：" + err.Error()
					return
				}
				global.WechatUUid[UUid].SendStatus = WECHAT_STATUS_SUCCESS
				global.WechatUUid[UUid].SendStatusMsg = "微信登陆成功"
				global.Log.Info(global.WechatUUid[UUid].SendStatusMsg)
				dao.InsertWechatLogin(convert.ToObjStr(global.WechatClientList[selfUserName]), global.Acc)
				//开启心跳检测
				SyncCheck(selfUserName)
				global.WechatClientList[selfUserName].LoginStatus = true
				dao.InsertHTTPCache(UUid, wc.SelfNickName, convert.ToObjStr(global.WechatClientList[selfUserName]), dao.WECHAT, global.Acc)
				break
			} else if status == 201 {
				global.WechatUUid[UUid].SendStatusMsg = "请在手机上确认登录"
				global.Log.Info(global.WechatUUid[UUid].SendStatusMsg)
			} else if status == 408 {
				global.WechatUUid[UUid].SendStatusMsg = "请扫描登录二维码"
				global.Log.Info(global.WechatUUid[UUid].SendStatusMsg)
			} else {
				global.WechatUUid[UUid].SendStatusMsg = fmt.Sprintf("未知情况，返回状态码：%d", status)
				global.Log.Info(global.WechatUUid[UUid].SendStatusMsg)
			}
			timeOut++
			time.Sleep(time.Second)
			if timeOut > 60*10 {
				global.WechatUUid[UUid].SendStatus = WECHAT_STATUS_FAIL
				global.WechatUUid[UUid].SendStatusMsg = "长时间未扫码，登录超时！"
				err = errors.New(global.WechatUUid[UUid].SendStatusMsg)
				return
			}
		} else {
			global.WechatUUid[UUid].SendStatus = WECHAT_STATUS_FAIL
			global.WechatUUid[UUid].SendStatusMsg = "前台已取消登录，结束登录验证"
			err = errors.New(global.WechatUUid[UUid].SendStatusMsg)
			return
		}
	}
	return
}

func SyncCheck(SelfUserName string) {
	global.Log.Info("准备开始微信登录状态心跳检测")
	defer func() {
		SyncCheckStatus = false
	}()
	for {
		if !LoginCheck(SelfUserName) {
			continue
		}
		SyncCheckStatus = true
		global.Log.Info("微信登录状态心跳检测")
		time.Sleep(time.Second)
	}
}

func LoginCheck(SelfUserName string) bool {
	flog := false
	if global.WechatClientList[SelfUserName] != nil {
		ret, selector, err := smartWechat.SyncCheck(global.WechatClientList[SelfUserName])
		if err == nil {
			flog = true
		} else {
			global.Log.Error("retcode：%v，selector：%v，ERROR：%v", ret, selector, err.Error())
		}
	}
	if !flog {
		global.WechatClientList[SelfUserName].LoginStatus = false
	}
	return flog
}

func AddWechatListUser(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_WECHAT) {
		return utils.ErrorNull(c, "微信营销功能还未开通，请开通此功能后使用")
	}

	SelfUserNames := c.FormValue("SelfUserNames")

	if SelfUserNames == "" {
		return utils.ErrorNull(c, "请选择可发送的微信或添加登录微信账号")
	}

	task := global.GetTask(global.TASK_WECHAT_ADD_GROUP_USER)
	if task != nil && task.State {
		return utils.ErrorNull(c, "任务正在运行中，请勿重复执行")
	}

	content := c.FormValue("content")
	go func() {

		defer func() {
			global.DelTask(global.TASK_WECHAT_ADD_GROUP_USER)
		}()

		SelfUserName := strings.Split(SelfUserNames, ",")
		sucNum := 0
		errNum := 0
		sleep := 10
		for _, j := range SelfUserName {
			if global.WechatClientList[j] == nil {
				global.PageMsg(fmt.Sprintf("【%s】微信号未登录或登录已失效", j))
				continue
			}
			nick := global.WechatClientList[j].SelfNickName
			ContactMap, err := smartWechat.GetAllContact(global.WechatClientList[j])
			if err != nil {
				global.Log.Info(fmt.Sprintf("【%s】微信账号，获取联系人信息，ERROR：%s", nick, err.Error()))
				global.PageMsg(fmt.Sprintf("【%s】微信账号，添加微信群成员中断，获取联系人信息，ERROR：%s", nick, err.Error()))
				return
			}
			if ContactMap == nil {
				global.PageMsg(fmt.Sprintf("【%s】微信账号，添加微信群成员中断，没有获取到联系人", nick))
				return
			}
			global.WechatClientList[j].ContactMap = ContactMap
			initInfo := global.WechatClientList[j].InitInfo
			if initInfo != nil {
				global.UpdateTask(global.TASK_WECHAT_ADD_GROUP_USER, fmt.Sprintf("【%s】微信账号，准备添加微信群成员..", nick))
				var str string
				for _, v := range initInfo.AllContactList {
					uIndex := 0
					for _, vv := range v.MemberList {
						uIndex++
						//任务记录
						task = global.GetTask(global.TASK_WECHAT_ADD_GROUP_USER)
						if task == nil || !task.State {
							str = fmt.Sprintf("【%s】任务已取消", global.TaskNameMap[global.TASK_WECHAT_ADD_GROUP_USER])
							global.Log.Info(str)
							global.PageMsg(str)
							return
						}
						//更新任务记录
						global.UpdateTask(global.TASK_WECHAT_ADD_GROUP_USER, fmt.Sprintf("【%s】微信账号，正在添加微信群【%s】第【%v】群成员", nick, v.NickName, uIndex))

						u := global.WechatClientList[j].ContactMap[vv.UserName]
						if u.UserName == "" {
							vu := smartWechat.VerifyUser{}
							vu.Value = vv.UserName
							br, err := smartWechat.AddUser(global.WechatClientList[j], content, []smartWechat.VerifyUser{vu})
							if err != nil {
								global.Log.Error("【%s】微信账号，【%v】%v 发送请求错误：%s", nick, v.NickName, vv.UserName, err.Error())
								errNum++
							} else {
								global.Log.Info("【%s】微信账号，【%v】%v 发送请求【%v】状态：%v", nick, v.NickName, vv.UserName, br.BaseResponse.Ret, br.BaseResponse.Ret == 0)
								if br.BaseResponse.Ret == 0 {
									sucNum++
									sleep = utils.NewRandom().Number(2)
								} else {
									if br.BaseResponse.Ret == enum.WECHAT_RESPONE_FREQUENTLY {
										//15分钟后继续
										sleep = 60 * 15
									} else {
										sleep = 60
									}
									errNum++
								}
							}
							str = fmt.Sprintf("【%s】微信账号，延迟添加微信群成员，等待【%v】秒后继续添加", nick, sleep)
							global.Log.Info(str)

							global.UpdateTask(global.TASK_WECHAT_ADD_GROUP_USER, str)
							time.Sleep(time.Second * time.Duration(sleep))
						} else {
							global.Log.Error("【%s】微信账号，我们是好友了啊【%v】%v", nick, v.NickName, vv.UserName)
						}
					}
				}
			}
		}
		global.PageMsg(fmt.Sprintf("添加微信群成员完成，请求成功：%s，失败：%s", sucNum, errNum))
	}()
	return utils.SuccessNull(c, "后台正在执行添加微信群成员，请务在短时间内重复发起！")
}

func SendWechatListMsg(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_WECHAT) {
		return utils.ErrorNull(c, "微信营销功能还未开通，请开通此功能后使用")
	}
	selfUserNames := c.FormValue("SelfUserNames")
	if selfUserNames == "" {
		return utils.ErrorNull(c, "请选择可发送的微信或添加登录微信账号")
	}
	msg := c.FormValue("msg")
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return utils.ErrorNull(c, "请输入待发送的内容")
	}

	task := global.GetTask(global.TASK_WECHAT_SEND_MESSAGE)
	if task != nil && task.State {
		return utils.ErrorNull(c, "正在发送中，请勿重复执行")
	}

	go SendWechatList(msg, selfUserNames)
	return utils.SuccessNull(c, "后台发送中...")
}

func SendWechatList(msg, selfUserNames string) (err error) {
	defer func() {
		global.DelTask(global.TASK_WECHAT_SEND_MESSAGE)
	}()
	var str string
	global.UpdateTask(global.TASK_WECHAT_SEND_MESSAGE, "准备开发群发微信消息..")

	selfUserName := strings.Split(selfUserNames, ",")
	for _, j := range selfUserName {
		if global.WechatClientList[j] == nil {
			global.Log.Info(fmt.Sprintf("【%s】登录账号已失效，请重新登录", j))
			continue
		}
		nick := global.WechatClientList[j].SelfNickName
		wu := global.WechatUUid[global.WechatClientList[j].UUID]
		/* 轮询服务器判断二维码是否扫过暨是否登陆了 */
		wu.SendStatusMsg = fmt.Sprintf("微信账号【%s】，开始获取联系人信息...", nick)
		global.Log.Info(wu.SendStatusMsg)
		ContactMap, err := smartWechat.GetAllContact(global.WechatClientList[j])
		if err != nil {
			wu.SendStatusMsg = fmt.Sprintf("微信账号【%s】，获取联系人信息，ERROR："+err.Error(), nick)
			global.Log.Info(wu.SendStatusMsg)
			continue
		}
		cm := convert.ToObjStr(ContactMap)
		if cm == "" {
			wu.SendStatusMsg = fmt.Sprintf("微信账号【%s】，没有获取到待发送的联系人信息", nick)
			global.Log.Info(wu.SendStatusMsg)
			continue
		}
		global.Log.Info(fmt.Sprintf("微信账号【%s】，联系人信息："+cm, nick))
		wu.SendStatusMsg = "【" + convert.ToString(len(ContactMap)) + "】准备群发消息..."
		global.Log.Info(wu.SendStatusMsg)
		wu.SendStatus = WECHAT_STATUS_PROCESS
		global.WechatClientList[j].ContactMap = ContactMap
		for k, v := range ContactMap {
			//任务记录
			task := global.GetTask(global.TASK_WECHAT_SEND_MESSAGE)
			if task == nil || !task.State {
				str = fmt.Sprintf("微信账号【%s】，【%s】已取消了", nick, global.TaskNameMap[global.TASK_WECHAT_SEND_MESSAGE])
				global.Log.Info(str)
				global.PageMsg(str)
				break
			}
			//更新任务记录
			global.UpdateTask(global.TASK_WECHAT_SEND_MESSAGE, fmt.Sprintf("微信账号【%s】，正在发送微信给用户【%s】", nick, v.NickName))

			global.Log.Info(v.UserName)
			global.Log.Info(ContactMap[k].UserName)
			if len(v.UserName) > 40 {
				//给所有人都发送消息
				wxSendMsg := smartWechat.WxSendMsg{}
				wxSendMsg.Type = 1
				wxSendMsg.Content = msg
				wxSendMsg.FromUserName = global.WechatClientList[j].SelfUserName
				wxSendMsg.ToUserName = v.UserName
				wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
				wxSendMsg.ClientMsgId = wxSendMsg.LocalID
				bts, err := smartWechat.SendMsg(global.WechatClientList[j], wxSendMsg)
				if err != nil {
					wu.SendStatusMsg = fmt.Sprintf("微信账号【%s】，错误：发送消息...，json:%s，error：%s", nick, convert.ToObjStr(wxSendMsg), err.Error())
					global.Log.Info(wu.SendStatusMsg)
				} else {
					wu.SendStatusMsg = fmt.Sprintf("微信账号【%s】，《%s》发送成功", nick, v.NickName)
					global.Log.Info(wu.SendStatusMsg + "，发送结果：" + string(bts))
				}
				sleep := utils.NewRandom().Number(1)
				global.UpdateTask(global.TASK_WECHAT_SEND_MESSAGE, fmt.Sprintf("微信账号【%s】，延迟【%v】秒后发送微信消息", nick, sleep))
				time.Sleep(time.Second * time.Duration(sleep))
			}
		}
		wu.SendStatus = WECHAT_STATUS_COMPLETE
		wu.SendStatusMsg = "微信群发消息任务完成！"
	}
	err = nil
	return
}
