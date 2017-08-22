package handler

import (
	"net/http"
	"time"
	"github.com/beewit/spread/global"
	"encoding/json"
	"strings"
)

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
			result="设置Cookie成功"
			break
		}
	}
	return flog, result
}
