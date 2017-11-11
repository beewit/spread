package handler

import (
	"encoding/json"
	"time"

	"fmt"
	"strings"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread/api"
	"github.com/beewit/spread/dao"
	"github.com/beewit/spread/global"
	"github.com/beewit/wechat-ai/ai"
	"github.com/beewit/wechat-ai/enum"
	"github.com/beewit/wechat-ai/smartWechat"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

var (
	SendStatusMsg    string
	SendStatus       int
	SyncCheckStatus  = false //时间开启了登录心跳检测
	WechatLoginCheck = true  //登录校验时间停止标记
)

const (
	WECHAT_STATUS_FAIL     = -1 //登录失败/其它地方登陆等
	WECHAT_STATUS_NOT      = 0  //未登录
	WECHAT_STATUS_SUCCESS  = 1  //登录成功
	WECHAT_STATUS_PROCESS  = 2  //获取联系人信息成功
	WECHAT_STATUS_COMPLETE = 3  //请求成功
)

func GetWechatFuncStatus(c echo.Context) error {
	flog := api.EffectiveFuncById(global.FUNC_WECHAT)
	return utils.SuccessNullMsg(c, flog)
}

func StartAddWechatGroup(c echo.Context) error {
	//权限校验
	if !api.EffectiveFuncById(global.FUNC_WECHAT) {
		return utils.ErrorNull(c, "微信营销功能还未开通，请开通此功能后使用")
	}
	area := c.FormValue("area")
	types := c.FormValue("type")
	groupCount := c.FormValue("groupCount")
	if groupCount == "" || !utils.IsValidNumber(groupCount) {
		groupCount = "3"
	}
	sleepTime := c.FormValue("sleepTime")
	if sleepTime == "" || !utils.IsValidNumber(sleepTime) {
		sleepTime = "30"
	}

	task := global.GetTask(global.TASK_WECHAT_ADD_GROUP)
	if task != nil && task.State {
		return utils.ErrorNull(c, "任务正在运行中，请勿重复执行")
	}
	go func() {
		defer func() {
			global.DelTask(global.TASK_WECHAT_ADD_GROUP)
		}()
		global.UpdateTask(global.TASK_WECHAT_ADD_GROUP, "准备开始添加微信群..")
		addWechat(area, types, convert.MustInt(groupCount), convert.MustInt(sleepTime))
	}()
	return utils.SuccessNull(c, "准备执行添加微信中..")
}

func addWechat(area, types string, groupCount, sleepTime int) {
	r, err := api.GetWechatGroupListData("1", area, types)
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
	var str string
	global.Navigate(global.LOAD_PAGE)
	gc := 0
	for i := 0; i < len(pageData.Data); i++ {
		//任务记录
		task := global.GetTask(global.TASK_WECHAT_ADD_GROUP)
		if task == nil || !task.State {
			str = fmt.Sprintf("【%s】任务已取消：%v ", global.TaskNameMap[global.TASK_WECHAT_ADD_GROUP], task == nil)
			global.Log.Info(str)
			global.PageSuccessMsg(str, global.Host+"?lastUrl=/app/page/admin/wechat/index.html")
			return
		}
		//更新任务记录
		global.UpdateTask(global.TASK_WECHAT_ADD_GROUP, fmt.Sprintf("正在前往获取微信群二维码：%s", convert.ToString(pageData.Data[i]["url"])))

		global.Page.Page.NextWindow()

		gc++
		if gc > groupCount {
			global.PageMsg(fmt.Sprintf("%v秒后再添加数据", sleepTime))
			time.Sleep(time.Duration(sleepTime) * time.Second)
			gc = 0
		}

		global.Navigate(convert.ToString(pageData.Data[i]["url"]))

		global.Page.RunScript(`$(".checkCode span:eq(1)").mouseover()`, nil, nil)
		time.Sleep(time.Second * 3)
		var of *enum.Offset
		err := global.Page.RunScript(`return  $(".shiftcode:eq(1) img").offset()`, nil, &of)
		if err != nil {
			global.Log.Error(" global.Page.RunScript ERROR:" + err.Error())
			continue
		}
		global.Log.Info("offset ： " + convert.ToObjStr(of))
		title, err := global.Page.Title()
		if err != nil {
			global.Log.Error("global.Page.Title ERROR:" + err.Error())
			continue
		}
		if of != nil {
			err = ai.Wechat(title, of)
			if err != nil {
				global.PageErrorMsg(err.Error()+"，已经停止添加微信群", global.Host+"?lastUrl=/app/page/admin/wechat/index.html")
				return
			}
			//记录添加记录
			flog, err := api.AddAccountWechatGroup(convert.MustInt64(pageData.Data[i]["id"]))
			if err != nil {
				global.Log.Error("AddAccountWechatGroup ERROR:" + err.Error())
				continue
			}
			global.Log.Info("AddAccountWechatGroup 添加：%v", flog)
		}
	}
	if pageData.Count == 1 {
		global.PageSuccessMsg("待添加的数据已完成", global.Host+"?lastUrl=/app/page/admin/wechat/index.html")
		return
	} else {
		addWechat(area, types, groupCount, sleepTime)
	}
}

func syncCheck() {
	global.Log.Info("准备开始微信登录状态心跳检测")
	defer func() {
		SyncCheckStatus = false
	}()
	for {
		if !loginCheck() {
			continue
		}
		SyncCheckStatus = true
		global.Log.Info("微信登录状态心跳检测")
		time.Sleep(time.Second)
	}
}

func loginCheck() bool {
	flog := false
	if global.WechatClient != nil {
		ret, selector, err := smartWechat.SyncCheck(global.WechatClient)
		if err == nil {
			flog = true
		} else {
			global.Log.Error("retcode：%v，selector：%v，ERROR：%v", ret, selector, err.Error())
		}
	}
	if !flog {
		global.WechatClient = nil
		dao.DeleteWechatLogin(global.Acc)
	}
	return flog
}

func LoginWechatCheck(c echo.Context) error {
	if global.WechatClient == nil {
		wl, err := dao.QueryWechatLogin(global.Acc.Id)
		if err != nil {
			global.Log.Error(err.Error())
		} else {
			if wl != "" {
				global.Log.Info("微信信息是数据库缓存")
				err = json.Unmarshal([]byte(wl), &global.WechatClient)
				if err != nil {
					global.Log.Info("微信信息是数据库缓存解析错误：")
					global.Log.Error(err.Error())
				}
			}
		}
	}
	if loginCheck() {
		SendStatus = WECHAT_STATUS_SUCCESS
		SendStatusMsg = "微信登陆成功"
		if !SyncCheckStatus {
			go syncCheck()
		}
		return utils.SuccessNullMsg(c, map[string]interface{}{"state": WECHAT_STATUS_SUCCESS, "wechatUser": global.WechatClient})
	}
	SendStatus = WECHAT_STATUS_FAIL
	SendStatusMsg = "微信未登陆"
	return utils.SuccessNullMsg(c, WECHAT_STATUS_FAIL)
}

func LoginWechat(c echo.Context) error {
	global.WechatClient = nil
	dao.DeleteWechatLogin(global.Acc)
	SendStatusMsg = "准备获取微信登录二维码"
	SendStatus = 0
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
	go loginWechat(UUid)
	return utils.Success(c, "扫描登录微信网页服务", base64Img)
}

func SendWechatMsg(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_WECHAT) {
		return utils.ErrorNull(c, "微信营销功能还未开通，请开通此功能后使用")
	}
	if global.WechatClient == nil {
		if global.WechatClient == nil {
			return utils.ErrorNull(c, "未登录，请重新扫码登录后群发消息")
		}
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

	go sendWechat(msg)
	return utils.SuccessNull(c, "后台发送中...")
}

func AddWechatUser(c echo.Context) error {
	if !api.EffectiveFuncById(global.FUNC_WECHAT) {
		return utils.ErrorNull(c, "微信营销功能还未开通，请开通此功能后使用")
	}
	if global.WechatClient == nil {
		return utils.ErrorNull(c, "未登录，请重新扫码登录后添加群成员")
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

		ContactMap, err := smartWechat.GetAllContact(global.WechatClient)
		if err != nil {
			global.Log.Info("获取联系人信息，ERROR：" + err.Error())
			global.PageMsg("添加微信群成员中断，获取联系人信息，ERROR：" + err.Error())
			return
		}
		if ContactMap == nil {
			global.PageMsg("添加微信群成员中断，没有获取到联系人")
			return
		}
		global.WechatClient.ContactMap = ContactMap
		sucNum := 0
		errNum := 0
		sleep := 10
		initInfo := global.WechatClient.InitInfo
		if initInfo != nil {
			global.UpdateTask(global.TASK_WECHAT_ADD_GROUP_USER, "准备添加微信群成员..")
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
					global.UpdateTask(global.TASK_WECHAT_ADD_GROUP_USER, fmt.Sprintf("正在添加微信群【%s】第【%v】群成员", v.NickName, uIndex))

					u := global.WechatClient.ContactMap[vv.UserName]
					if u.UserName == "" {
						vu := smartWechat.VerifyUser{}
						vu.Value = vv.UserName
						br, err := smartWechat.AddUser(global.WechatClient, content, []smartWechat.VerifyUser{vu})
						if err != nil {
							global.Log.Error("【%v】%v 发送请求错误：%s", v.NickName, vv.UserName, err.Error())
							errNum++
						} else {
							global.Log.Info("【%v】%v 发送请求【%v】状态：%v", v.NickName, vv.UserName, br.BaseResponse.Ret, br.BaseResponse.Ret == 0)
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
						global.Log.Info("延迟添加微信群成员，等待【%v】秒后继续添加", sleep)

						global.UpdateTask(global.TASK_WECHAT_ADD_GROUP_USER, fmt.Sprintf("延迟添加微信群成员，等待【%v】秒后继续添加", sleep))
						time.Sleep(time.Second * time.Duration(sleep))
					} else {
						global.Log.Error("我们是好友了啊【%v】%v", v.NickName, vv.UserName)
					}
				}
			}
		}
		global.PageMsg(fmt.Sprintf("添加微信群成员完成，请求成功：%s，失败：%s", sucNum, errNum))
	}()
	return utils.SuccessNull(c, "后台正在执行添加微信群成员，请务在短时间内重复发起！")
}

func GetSendWechatMsgStatus(c echo.Context) error {
	return utils.SuccessNullMsg(c, map[string]interface{}{"sendStatusMsg": SendStatusMsg, "sendStatus": SendStatus})
}

func CancelLoginWechat(c echo.Context) error {
	WechatLoginCheck = false
	return utils.SuccessNull(c, "")
}

func loginWechat(UUid string) (err error) {
	WechatLoginCheck = true
	timeOut := 0
	for {
		thisUrl, _ := global.Page.Page.URL()
		if WechatLoginCheck && strings.Contains(thisUrl, "/app/page/admin/wechat/index.html") {
			SendStatusMsg = "【" + UUid + "】正在验证登陆... ..."
			global.Log.Info(SendStatusMsg)
			status, msg := smartWechat.CheckLogin(UUid)
			if status == 200 {
				SendStatusMsg = "登陆成功,处理登陆信息..."
				global.Log.Info(SendStatusMsg)
				global.WechatClient, err = smartWechat.ProcessLoginInfo(msg)
				if err != nil {
					SendStatus = WECHAT_STATUS_FAIL
					SendStatusMsg = "错误：登陆成功,处理登陆信息...，error：" + err.Error()
					global.Log.Info(SendStatusMsg)
					return
				}
				SendStatusMsg = "登陆信息处理完毕,正在初始化微信..."
				global.Log.Info(SendStatusMsg)
				global.Log.Info("global.WechatClient：%s", convert.ToObjStr(global.WechatClient))
				err = smartWechat.InitWX(global.WechatClient)
				if err != nil {
					SendStatusMsg = "【1】错误：登陆信息处理完毕,正在初始化微信...，error：" + err.Error()
					global.Log.Info(SendStatusMsg)
					if global.WechatClient.InitInfo.BaseResponse.Ret == 1010 {
						err = smartWechat.InitWX(global.WechatClient)
					}
					if err != nil {
						SendStatusMsg = "【2】错误：登陆信息处理完毕,正在初始化微信...，error：" + err.Error()
						SendStatus = WECHAT_STATUS_FAIL
						global.Log.Info(SendStatusMsg)
						return
					}
				}
				SendStatusMsg = "初始化完毕,通知微信服务器登陆状态变更..."
				global.Log.Info(SendStatusMsg)
				err = smartWechat.NotifyStatus(global.WechatClient)
				if err != nil {
					SendStatus = WECHAT_STATUS_FAIL
					SendStatusMsg = "通知微信服务器状态变化失败：" + err.Error()
					return
				}
				SendStatus = WECHAT_STATUS_SUCCESS
				SendStatusMsg = "微信登陆成功"
				global.Log.Info(SendStatusMsg)
				dao.InsertWechatLogin(convert.ToObjStr(global.WechatClient), global.Acc)
				//开启心跳检测
				syncCheck()
				break
			} else if status == 201 {
				SendStatusMsg = "请在手机上确认登录"
				global.Log.Info(SendStatusMsg)
			} else if status == 408 {
				SendStatusMsg = "请扫描登录二维码"
				global.Log.Info(SendStatusMsg)
			} else {
				SendStatusMsg = fmt.Sprintf("未知情况，返回状态码：%d", status)
				global.Log.Info(SendStatusMsg)
			}
			timeOut++
			time.Sleep(time.Second)
			if timeOut > 60*10 {
				SendStatus = WECHAT_STATUS_FAIL
				SendStatusMsg = "长时间未扫码，登录超时！"
				err = errors.New(SendStatusMsg)
				return
			}
		} else {
			err = errors.New("前台已取消登录，结束登录验证")
			return
		}
	}
	return
}

func sendWechat(msg string) (err error) {
	defer func() {
		global.DelTask(global.TASK_WECHAT_SEND_MESSAGE)
	}()
	var str string
	global.UpdateTask(global.TASK_WECHAT_SEND_MESSAGE, "准备开发群发微信消息..")

	/* 轮询服务器判断二维码是否扫过暨是否登陆了 */
	SendStatusMsg = "开始获取联系人信息..."
	global.Log.Info(SendStatusMsg)
	ContactMap, err := smartWechat.GetAllContact(global.WechatClient)
	if err != nil {
		SendStatusMsg = "获取联系人信息，ERROR：" + err.Error()
		global.Log.Info(SendStatusMsg)
		return
	}
	cm := convert.ToObjStr(ContactMap)
	if cm == "" {
		SendStatusMsg = "没有获取到待发送的联系人信息"
		global.Log.Info(SendStatusMsg)
		return
	}
	global.Log.Info("联系人信息" + cm)
	SendStatusMsg = "【" + convert.ToString(len(ContactMap)) + "】准备群发消息..."
	global.Log.Info(SendStatusMsg)
	SendStatus = WECHAT_STATUS_PROCESS
	global.WechatClient.ContactMap = ContactMap
	for k, v := range ContactMap {
		//任务记录
		task := global.GetTask(global.TASK_WECHAT_SEND_MESSAGE)
		if task == nil || !task.State {
			str = fmt.Sprintf("【%s】已取消了", global.TaskNameMap[global.TASK_WECHAT_SEND_MESSAGE])
			global.Log.Info(str)
			global.PageMsg(str)
			return
		}
		//更新任务记录
		global.UpdateTask(global.TASK_WECHAT_SEND_MESSAGE, fmt.Sprintf("正在发送微信给用户【%s】", v.NickName))

		global.Log.Info(v.UserName)
		global.Log.Info(ContactMap[k].UserName)
		if len(v.UserName) > 40 {
			//给所有人都发送消息
			wxSendMsg := smartWechat.WxSendMsg{}
			wxSendMsg.Type = 1
			wxSendMsg.Content = msg
			wxSendMsg.FromUserName = global.WechatClient.SelfUserName
			wxSendMsg.ToUserName = v.UserName
			wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
			wxSendMsg.ClientMsgId = wxSendMsg.LocalID
			bts, err := smartWechat.SendMsg(global.WechatClient, wxSendMsg)
			if err != nil {
				SendStatusMsg = "错误：发送消息...，json:" + convert.ToObjStr(wxSendMsg) + "，error：" + err.Error()
				global.Log.Info(SendStatusMsg)
			} else {
				SendStatusMsg = v.NickName + "发送成功"
				global.Log.Info(SendStatusMsg + "，发送结果：" + string(bts))
			}
			sleep := utils.NewRandom().Number(1)
			global.UpdateTask(global.TASK_WECHAT_SEND_MESSAGE, fmt.Sprintf("延迟【%v】秒后发送微信消息", sleep))
			time.Sleep(time.Second * time.Duration(sleep))
		}
	}
	SendStatus = WECHAT_STATUS_COMPLETE
	SendStatusMsg = "微信群发消息任务完成！"
	global.Log.Info(SendStatusMsg)
	global.PageSuccessMsg(SendStatusMsg, global.Host+"?lastUrl=/app/page/admin/wechat/index.html")
	err = nil
	return
}
