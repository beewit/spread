package handler

import (
	"net/http"
	"time"
	"github.com/beewit/spread/global"
	"encoding/json"
	"strings"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/beekit/utils"
	"github.com/labstack/echo"
)

func CheckClientLogin(token string) bool {
	b, err := uhttp.PostForm(global.API_SSO_DOMAN+"/pass/checkToken?token="+token, nil)
	if err != nil {
		global.Log.Error(err.Error())
		return false
	}
	rp := utils.ToResultParam(b)
	if rp.Ret != 200 {
		return false
	}
	return true
}

func SetClientToken(token string) bool {
	if CheckClientLogin(token) {
		//insert sqllite
		flog, err := global.InsertToken(token)
		if err != nil {
			global.Log.Error(err.Error())
			return false
		}
		return flog
	}
	return false
}

func ReceiveToken(c echo.Context) error {
	token := c.FormValue("token")
	if token != "" {
		if SetClientToken(token) {
			return utils.Redirect(c, global.Host)
		}
	}
	return utils.RedirectAndAlert(c, "登陆失败", global.API_SSO_DOMAN)
}

func SetLoginCookie() (bool, string) {
	jsAccount, _ := global.RD.GetString("jianshu")
	if jsAccount == "" {
		return false, "Redis无效"
	}
	var cks = []*http.Cookie{}
	err := json.Unmarshal([]byte(jsAccount), &cks)
	if err != nil {
		return false, "解码失败"
	}
	global.Page.Navigate("http://www.jianshu.com")
	for i := range cks {
		cc := cks[i]
		global.Page.SetCookie(cc)
	}
	time.Sleep(500 * time.Millisecond)
	flog, result := CheckLogin("jianshu", "www.jianshu.com", "remember_user_token")
	return flog, result
}

func CheckLogin(t string, domain string, identity string) (bool, string) {
	flog := false
	result := ""
	for {
		println("检测登陆状态")
		url, _ := global.Page.URL()
		if !strings.Contains(url, domain) {
			result = "已经不在本网站了，结束检测登陆状态"
			println(result)
			break
		}
		c, _ := global.Page.GetCookies()
		for _, apiCookie := range c {
			if apiCookie.Name == identity && apiCookie.Value != "" {
				flog = true
				break
			}
		}

		time.Sleep(1 * time.Second)

		if flog {
			cookieJson, _ := json.Marshal(c)

			println(cookieJson)
			global.RD.SetString(t, cookieJson)

			//global.Page.Navigate(global.Url)
			result = "设置Cookie成功"
			break
		}
	}
	return flog, result
}
