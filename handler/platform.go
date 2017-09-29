package handler

import (
	"encoding/json"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread/api"
	"github.com/beewit/spread/dao"
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/parser"
	"github.com/labstack/echo"
	"time"
)

//func PlatformList(c echo.Context) error {
//	m, err := api.GetPlatformList(c)
//	if err != nil || m == nil {
//		global.Log.Error(err.Error())
//		return utils.Error(c, "获取平台信息失败", nil)
//	}
//	return utils.Success(c, "", m)
//}

func UnionList(c echo.Context) error {
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	page, err := dao.GetUnionListPage(global.Acc.Id, pageIndex, global.PAGE_SIZE)
	if err != nil {
		global.Log.Error(err.Error())
		return utils.Error(c, "获取绑定帐号异常", nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "", page)
}

func PlatformUnionBind(c echo.Context) error {
	t := c.FormValue("type")
	if t == "" {
		return utils.Error(c, "请选择平台进行绑定", nil)
	}
	//获取远程服务中的平台信息进行帐号绑定操作
	m, err := api.GetPlatformOne(t)
	if err != nil || m == nil {
		return utils.Error(c, "绑定信息失败，获取平台信息失败", nil)
	}
	go UnionBind(m)
	return utils.Success(c, "正在前往绑定中...", "")
}

func UnionBind(m map[string]interface{}) {
	if global.Acc == nil || global.Acc.Id <= 0 {
		global.PageAlertMsg("帐号未登陆，请登陆后操作", global.API_SSO_DOMAN+"?backUrl="+global.Host+"/ReceiveToken")
		return
	}
	//进入登陆页面
	platformId := convert.MustInt64(m["id"])
	lu := convert.ToString(m["login_url"])
	domain := convert.ToString(m["site_url"])
	platform := convert.ToString(m["type"])
	identity := convert.ToString(m["identity"])
	as := convert.ToString(m["account_selector"])
	ps := convert.ToString(m["password_selector"])
	iframe := convert.ToString(m["iframe"])
	global.Navigate(lu)
	parser.DeleteCookie()
	//检测登陆状态
	flog, platformAcc := checkLogin(domain, identity, platform, as, ps, iframe, platformId)
	if flog {
		infoUrl := convert.ToString(m["info_url"])
		ns := convert.ToString(m["info_nickname_selector"])
		ps := convert.ToString(m["info_photo_selector"])
		if infoUrl != "" && (ns != "" || ps != "") {
			global.Navigate(infoUrl)
			time.Sleep(2 * time.Second)
			nickname := global.PageFindValue(ns)
			photo := global.PageFindValue(ps)
			if nickname != "" || photo != "" {
				uFlog, _ := dao.UpdateUnionPhoto(nickname, photo, platformAcc, platformId, global.Acc.Id)
				global.Log.Warning("修改帐号昵称和帐号信息，状态：%v", uFlog)
			}
		}
		global.PageSuccessMsg("绑定帐号成功", global.Host+"?lastUrl=/app/page/admin/account/list.html")
	} else {
		global.PageErrorMsg("绑定帐号失败", global.Host+"?lastUrl=/app/page/admin/account/list.html")
	}
}

func checkLogin(domain, identity, platform, as, ps, iframeSeletor string, platformId int64) (bool, string) {
	if iframeSeletor != "" {
		time.Sleep(time.Second * 1)
		html, _ := global.Page.HTML()
		println(html)
		iframe, err := global.Page.Find(iframeSeletor).Elements()
		if err != nil {
			println(err.Error())
		}
		if len(iframe) <= 0 {
			return false, "切换iframe登陆失败"
		}
		err = global.Page.SwitchToRootFrameByName(iframe[0])
		if err != nil {
			println(err.Error())
		}
		defer global.Page.SwitchToParentFrame()
	}
	flog := false
	result := ""
	i := 0
	var platformAcc, platformPwd string
	for {
		global.Log.Info("检测登陆状态")
		//获取帐号密码
		a := global.PageFindValue(as)
		p := global.PageFindValue(ps)

		if a != "" {
			platformAcc = a
			global.Log.Info("帐号：" + platformAcc)
		}
		if p != "" {
			platformPwd = p
			global.Log.Info("密码：" + platformPwd)
		}

		if parser.CheckStopAtSite(domain) {
			result = "已经不在本网站了，结束检测登陆状态"
			global.Log.Info(result)
			break
		}

		c, f := parser.CheckIdentity(identity)
		flog = f

		if flog {
			if platformAcc != "" && platformPwd != "" {
				dbFlog, err := dao.SetUnion(platform, platformAcc, platformPwd, platformId, global.Acc.Id)
				if err != nil {
					global.Log.Error("添加帐号绑定数据，异常：%v", err.Error())
				}
				global.Log.Warning("添加帐号绑定数据，状态：%v", dbFlog)
			}
			cookieJson, _ := json.Marshal(c)

			global.Log.Info(string(cookieJson[:]))
			dao.SetUnionCookies(domain, string(cookieJson), platformId, global.Acc.Id, platformAcc)
			result = "登陆成功"
			break
		}
		i++
		if i > 60*10*1000 {
			//10分钟退出循环
			break
		}
	}
	if result != "" {
		global.Log.Info(result)
	}
	return flog, platformAcc
}
