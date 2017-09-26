package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/sclevine/agouti"

	"github.com/beewit/beekit/redis"
	"github.com/sclevine/agouti/api"
)

func TestJianshu(t *testing.T) {
	driver := agouti.ChromeDriver(agouti.ChromeOptions("args", []string{
		"--start-maximized",
		"--disable-infobars",
		"--app=http://www.jianshu.com/",
		"--webkit-text-size-adjust"}))
	driver.Start()
	page, err := driver.NewPage()
	if err != nil {
		println(err.Error())
	} else {
		li, _ := page.Find(".note-list").All(".have-img").Elements()

		var s string
		for i := range li {
			s, _ = li[i].GetText()

			println(s)
		}
	}
}
func TestQQ(t *testing.T) {
	driver := agouti.ChromeDriver(agouti.ChromeOptions("args", []string{
		"--start-maximized",
		"--disable-infobars",
		"--app=https://i.qq.com/?rd=1",
		"--webkit-text-size-adjust"}))
	driver.Start()
	page, err := driver.NewPage()
	if err != nil {
		println(err.Error())
	} else {
		//html, _ := page.HTML()
		//println(html)
		flog, ce := setCookieLogin(page, "qqzone294477044")
		if ce != nil {
			println(ce.Error())
			return
		}
		if flog {
			page.Navigate("https://user.qzone.qq.com/294477044")
			time.Sleep(time.Second * 3)
			src, ei := page.Find(".head-avatar img").Attribute("src")
			if ei != nil {
				println(ei.Error())
			}
			if src != "" {
				println("头像", src)
			} else {
				flog = false
			}
		}
		if !flog {
			page.Navigate("https://i.qq.com/?rd=1")
			time.Sleep(time.Second * 3)
			ifarme, ee := page.Find("#login_frame").Elements()
			if ee != nil {
				println(ee.Error())
			}
			e2 := page.SwitchToRootFrameByName(ifarme[0])
			if e2 != nil {
				println(e2.Error())
				return
			}
			text, e3 := page.Find("#switcher_plogin").Text()
			if e3 != nil {
				println(e3.Error())
				return
			}
			println("登陆按钮", text)
			e4 := page.Find("#switcher_plogin").Click()
			if e4 != nil {
				println("登陆失败", e4.Error())
				return
			}
			page.FindByID("u").Fill("294477044")
			time.Sleep(time.Second * 1)
			page.FindByID("p").Fill("zxblovelc520~")
			time.Sleep(time.Second * 2)
			page.FindByID("login_button").Click()
			time.Sleep(time.Second * 2)
			page.SwitchToParentFrame()
			time.Sleep(time.Second * 3)
			c, err5 := page.GetCookies()
			if err5 != nil {
				println("登陆失败", e4.Error())
				return
			}
			cookieJson, _ := json.Marshal(c)
			//println("cookie", string(cookieJson[:]))
			redis.Cache.SetString("qqzone294477044", string(cookieJson[:]))

		}
		//数据
		count := 0
		getInfo(page, count)
	}
}

func getInfo(page *agouti.Page, count int) {
	list, e6 := page.Find("#feed_friend_list").All(".f-single").Elements()
	if e6 != nil {
		println("获取好友数据失败", e6.Error())
		return
	}
	println("总数量", len(list))
	var s string
	var e7 error
	var ele *api.Element
	for i := range list {

		println("---------------------------------------------------------\r\n")

		s, e7 = list[i].GetAttribute("id")
		if e7 != nil {
			println("错误：", e7.Error())
		}
		println("id：", s)
		ele, e7 = list[i].GetElement(api.Selector{"css selector", ".user-pto img"})
		if e7 != nil {
			println("错误：", e7.Error())
		}
		s, e7 = ele.GetAttribute("src")
		println("头像：", s)

		ele, _ = list[i].GetElement(api.Selector{"css selector", ".user-pto a"})
		if e7 != nil {
			println("错误：", e7.Error())
		}
		s, e7 = ele.GetAttribute("href")
		println("空间链接：", s)

		ele, _ = list[i].GetElement(api.Selector{"css selector", ".f-single-content"})
		if e7 != nil {
			println("错误：", e7.Error())
		}
		s, e7 = ele.GetText()
		println("发表内容：", s)

		ele, _ = list[i].GetElement(api.Selector{"css selector", ".qz_feed_plugin"})
		if e7 != nil {
			println("错误：", e7.Error())
		}
		s, e7 = ele.GetText()
		println("浏览量：", s)

		ele, _ = list[i].GetElement(api.Selector{"css selector", ".comments-list"})
		if e7 != nil {
			println("错误：", e7.Error())
		}
		s, e7 = ele.GetText()
		println("评论：", s)

		ele, _ = list[i].GetElement(api.Selector{"css selector", ".user-list"})
		if e7 != nil {
			println("错误：", e7.Error())
		}
		s, e7 = ele.GetText()
		println("点赞：", s)

		println("---------------------------------------------------------\r\n")
	}
	page.RunScript("document.documentElement.scrollTop=document.body.clientHeight;", nil, nil)
	time.Sleep(time.Second * 3)
	count++
	if count < 3 {
		getInfo(page, count)
	}
}

func setCookieLogin(page *agouti.Page, key string) (bool, error) {
	cookieRd, err := redis.Cache.GetString(key)
	if err != nil {
		return false, err
	}
	if cookieRd == "" {
		return false, nil
	}
	var cks = []*http.Cookie{}
	err = json.Unmarshal([]byte(cookieRd), &cks)
	if err != nil {
		return false, err
	}
	for i := range cks {
		cc := cks[i]
		page.SetCookie(cc)
	}
	return true, nil
}
