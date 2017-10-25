package handler

import (
	"encoding/json"
	"time"

	"fmt"
	"strings"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread/api"
	"github.com/beewit/spread/global"
	"github.com/beewit/wechat-ai/enum"
	"github.com/beewit/wechat-ai/send"
	"github.com/labstack/echo"
	"github.com/beewit/wechat-ai/ai"
)

var (
	UUid          string
	LoginMap      send.LoginMap
	SendStatusMsg string
	SendStatus    int
	ContactMap    map[string]send.User
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
	go addWechat(c, "1", area, types, convert.MustInt(groupCount), convert.MustInt(sleepTime))
	return utils.SuccessNull(c, "准备执行添加微信中..")
}

func addWechat(c echo.Context, pageIndex, area, types string, groupCount, sleepTime int) {
	r, err := api.GetWechatGroupListData(pageIndex, area, types)
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
	global.Navigate(global.LoadPage)
	gc := 0
	for i := 0; i < len(pageData.Data); i++ {
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
		addWechat(c, "1", area, types, groupCount, sleepTime)
	}
}

func SendWechatMsg(c echo.Context) error {
	msg := c.FormValue("msg")
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return utils.ErrorNull(c, "请输入待发送的内容")
	}
	/* 从微信服务器获取UUID */
	var err error
	UUid, err = send.GetUUIDFromWX()
	if err != nil {
		return utils.ErrorNull(c, "GetUUIDFromWX Error："+err.Error())
	}
	/* 根据UUID获取二维码 */
	base64Img, err := send.DownloadImage(enum.QRCODE_URL + UUid)
	if err != nil {
		return utils.ErrorNull(c, "DownloadImage Error："+err.Error())
	}
	SendStatusMsg = ""
	SendStatus = 0
	go sendWechat(msg)
	return utils.Success(c, "扫描登录微信网页服务", base64Img)
}

func GetSendWechatMsgStatus(c echo.Context) error {
	return utils.SuccessNullMsg(c, map[string]interface{}{"sendStatusMsg": SendStatusMsg, "sendStatus": SendStatus})
}

func sendWechat(msg string) {
	/* 轮询服务器判断二维码是否扫过暨是否登陆了 */
	var err error
	for {
		SendStatusMsg = "【" + UUid + "】正在验证登陆... ..."
		global.Log.Info(SendStatusMsg)
		status, msg := send.CheckLogin(UUid)
		if status == 200 {
			SendStatusMsg = "登陆成功,处理登陆信息..."
			global.Log.Info(SendStatusMsg)
			LoginMap, err = send.ProcessLoginInfo(msg)
			if err != nil {
				SendStatusMsg = "错误：登陆成功,处理登陆信息...，error：" + err.Error()
				global.Log.Info(SendStatusMsg)
				return
			}
			SendStatusMsg = "登陆信息处理完毕,正在初始化微信..."
			global.Log.Info(SendStatusMsg)
			err = send.InitWX(&LoginMap)
			if err != nil {
				if err != nil {
					SendStatusMsg = "错误：登陆信息处理完毕,正在初始化微信...，error：" + err.Error()
					global.Log.Info(SendStatusMsg)
					return
				}
			}
			SendStatusMsg = "初始化完毕,通知微信服务器登陆状态变更..."
			global.Log.Info(SendStatusMsg)
			err = send.NotifyStatus(&LoginMap)
			if err != nil {
				panic(err)
			}
			SendStatusMsg = "通知完毕,本次登陆信息获取成功"
			global.Log.Info(SendStatusMsg)
			//fmt.Println(enum.SKey + "\t\t" + loginMap.BaseRequest.SKey)
			//fmt.Println(enum.PassTicket + "\t\t" + loginMap.PassTicket)
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
	}
	SendStatusMsg = "开始获取联系人信息..."
	global.Log.Info(SendStatusMsg)
	ContactMap, err = send.GetAllContact(&LoginMap)
	if err != nil {
		SendStatusMsg = "错误：开始获取联系人信息...，error：" + err.Error()
		global.Log.Info(SendStatusMsg)
	}
	ss := convert.ToObjStr(ContactMap)

	global.Log.Info("联系人信息" + ss)
	SendStatusMsg = "【" + convert.ToString(len(ContactMap)) + "】准备群发消息..."
	global.Log.Info(SendStatusMsg)
	SendStatus = 1

	for k, v := range ContactMap {
		println(v.UserName)
		println(ContactMap[k].UserName)
		if len(v.UserName) > 40 {
			//为人，发送消息
			wxSendMsg := send.WxSendMsg{}
			wxSendMsg.Type = 1
			wxSendMsg.Content = msg
			wxSendMsg.FromUserName = LoginMap.SelfUserName
			wxSendMsg.ToUserName = v.UserName
			wxSendMsg.LocalID = fmt.Sprintf("%d", time.Now().Unix())
			wxSendMsg.ClientMsgId = wxSendMsg.LocalID
			bts, err := send.SendMsg(&LoginMap, wxSendMsg)
			if err != nil {
				SendStatusMsg = "错误：发送消息...，json:" + convert.ToObjStr(wxSendMsg) + "，error：" + err.Error()
				global.Log.Info(SendStatusMsg)
			} else {
				SendStatusMsg = v.NickName + "发送成功"
				global.Log.Info(SendStatusMsg + "，发送结果：" + string(bts))
			}
			time.Sleep(time.Second * 2)
		}
	}
	SendStatus = 2
	SendStatusMsg = "群发消息完成！"
	global.Log.Info(SendStatusMsg)
	global.PageSuccessMsg("群发消息完成！", global.Host+"?lastUrl=/app/page/admin/wechat/index.html")
}
