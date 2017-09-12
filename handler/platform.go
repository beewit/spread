package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/api"
	"net/url"
	"strings"
	"github.com/beewit/spread/dao"
	"encoding/json"
	"github.com/beewit/spread/parser"
	"time"
)

func PlatformList(c echo.Context) error {
	m, err := api.GetPlatformList()
	if err != nil || m == nil {
		global.Log.Error(err.Error())
		return utils.Error(c, "获取平台信息失败", nil)
	}
	return utils.Success(c, "", m)
}

func UnionList(c echo.Context) error {
	m, err := dao.GetUnionList(global.Acc.Id)
	if err != nil || m == nil {
		global.Log.Error(err.Error())
		return utils.Error(c, "获取平台信息失败", nil)
	}
	return utils.Success(c, "", m)
}

func PlatformUnionBind(c echo.Context) error {
	t := c.FormValue("type")
	if t == "" {
		return utils.Error(c, "请选择平台进行绑定", nil)
	}
	pv := global.Platform[t]
	if pv < 0 {
		return utils.Error(c, "请正确选择平台进行绑定", nil)
	}
	//获取远程服务中的平台信息进行帐号绑定操作
	m, err := api.GetPlatformOne(t)
	if err != nil || m == nil {
		global.Log.Error(err.Error())
		return utils.Error(c, "绑定信息失败，获取平台信息失败", nil)
	}
	go UnionBind(m)
	return utils.Success(c, "正在前往绑定中...", "")
}

func UnionBind(m map[string]string) {
	if global.Acc == nil || global.Acc.Id <= 0 {
		global.PageAlertMsg("帐号未登陆，请登陆后操作", global.API_SSO_DOMAN+"?backUrl="+global.Host+"/ReceiveToken")
		return
	}
	//进入登陆页面
	lu := m["login_url"]
	domain := m["site_url"]
	platform := m["type"]
	identity := m["identity"]
	as := m["account_selector"]
	ps := m["password_selector"]
	global.Navigate(lu)
	//检测登陆状态
	flog, _ := checkLogin(domain, identity, platform, as, ps)
	if flog {
		infoUrl := m["info_url"]
		ns := m["info_nickname_selector"]
		ps := m["info_photo_selector"]
		if infoUrl != "" && (ns != "" || ps != "") {
			global.Navigate(infoUrl)
			time.Sleep(2 * time.Second)
			nickname := global.PageFindValue(ns)
			photo := global.PageFindValue(ps)
			if nickname != "" || photo != "" {
				uFlog, _ := dao.UpdateUnionPhoto(nickname, photo, platform, global.Acc.Id)
				global.Log.Warning("修改帐号昵称和帐号信息，状态：", uFlog)
			}
		}
		global.PageSuccessMsg("绑定帐号成功", global.Host)
	} else {
		global.PageErrorMsg("绑定帐号失败", global.Host)
	}
}

func checkLogin(domain, identity, platform, as, ps string) (bool, string) {
	flog := false
	result := ""
	i := 0
	for {
		global.Log.Info("检测登陆状态")
		//获取帐号密码
		acc := global.PageFindValue(as)
		pwd := global.PageFindValue(ps)
		if acc != "" {
			global.Log.Info("帐号：" + acc)
		}
		if pwd != "" {
			global.Log.Info("密码：" + pwd)
		}
		if acc != "" || pwd != "" {
			dbFlog, err := dao.SetUnion(platform, acc, pwd, global.Acc.Id)
			if err != nil {
				global.Log.Error("添加帐号绑定数据，异常：%v", err.Error())
			}
			global.Log.Warning("添加帐号绑定数据，状态：%v", dbFlog)
		}
		thisUrl, _ := global.Page.URL()
		u, _ := url.Parse(domain)
		if !strings.Contains(thisUrl, u.Host) {
			result = "已经不在本网站了，结束检测登陆状态"
			global.Log.Info(result)
			break
		}

		c, f := parser.CheckIdentity(identity)
		flog = f

		if flog {
			cookieJson, _ := json.Marshal(c)

			global.Log.Info(string(cookieJson[:]))
			dao.SetUnionCookies(domain, string(cookieJson), global.Acc.Id)
			result = "登陆成功"
			global.Log.Info(result)
			break
		}
		i++
		if i > 60*10*1000 {
			//10分钟退出循环
			break
		}
	}
	return flog, result
}
