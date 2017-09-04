package parser

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"net/url"

	"github.com/beewit/spread/global"
	"github.com/sclevine/agouti"
	"github.com/beewit/spread/dao"
)

type PushJson struct {
	Title     string         `json:"title"`
	Domain    string         `json:"domain"`
	LoginUrl  string         `json:"loginUrl"`
	Identity  string         `json:"identity"`
	WriterUrl string         `json:"writerUrl"`
	Sleep     int64          `json:"sleep"`
	Fill      []PushFillJson `json:"fill"`
	Login     []PushFillJson `json:"login"`
}

type PushFillJson struct {
	Selector     string                 `json:"selector"`
	SelectorName string                 `json:"selectorName"`
	SelectorVal  string                 `json:"selectorVal"`
	Handle       string                 `json:"handle"`
	Sleep        int64                  `json:"sleep"`
	Js           string                 `json:"js"`
	JsParam      map[string]interface{} `json:"jsParam"`
	Param        string                 `json:"param"`
	Result       string                 `json:"result"`
}

const (
	Click       = "Click"
	DoubleClick = "DoubleClick"
	Check       = "Check"
	Uncheck     = "Uncheck"
	Select      = "Select"
	Submit      = "Submit"
	RunScript   = "Js"
	Fill        = "Fill"
	Text        = "Text"
)

const (
	Find         = "Selector"
	FindByID     = "ID"
	FindByXPath  = "XPath"
	FindByLink   = "Link"
	FindByLabel  = "Label"
	FindByButton = "Button"
	FindByName   = "Name"
	FindByClass  = "Class"
)

func getPushJson(rule string) (*PushJson, error) {
	var pj PushJson
	err := json.Unmarshal([]byte(rule), &pj)
	if err != nil {
		panic(err)
		return &PushJson{}, err
	} else {
		return &pj, nil
	}
}

func RunPush(rule string, paramMap map[string]string) (bool, string, error) {
	pj, err := getPushJson(rule)
	if err != nil {
		return false, pj.Title + "解析配置规则失败", err
	}
	if len(pj.Fill) <= 0 {
		return false, pj.Title + "无发布规则配置", nil
	}
	var flog bool
	if pj.Login != nil && len(pj.Login) > 0 {
		global.Navigate(pj.LoginUrl)
		//Login Set UserName And Password
		global.Log.Info("Login Set UserName And Password")
		for i := 0; i < len(pj.Login); i++ {
			handleSelection(&pj.Login[i], paramMap)
		}
		flog, _ = checkLogin(pj.Domain, pj.Identity)
	} else {
		global.Navigate(pj.Domain)
		//Login Identity
		_, flog = CheckIdentity(pj.Identity)
		if !flog {
			flog, _ = setCookieLogin(pj.Domain)
			if !flog {
				global.Navigate(pj.LoginUrl)
				//Sending messages to users requires landing
			} else {
				//设置Cookie
				global.Log.Info("设置Cookie成功")
			}
			flog, _ = checkLogin(pj.Domain, pj.Identity)
		} else {
			global.Log.Info("已经是登陆状态")
		}
	}
	if !flog {
		return false, pj.Title + "登陆失败", nil
	}
	global.Navigate(pj.WriterUrl)
	if pj.Sleep > 0 {
		time.Sleep(time.Duration(pj.Sleep) * time.Second)
	}
	for i := 0; i < len(pj.Fill); i++ {
		handleSelection(&pj.Fill[i], paramMap)
	}
	return true, "全部执行完成", nil
}

func handleSelection(p *PushFillJson, paramMap map[string]string) (bool, string, error) {
	var jsResult string
	switch p.Handle {
	case Click:
		global.Log.Info("Click：", findSelection(p.Selector, p.SelectorName).String())
		findSelection(p.Selector, p.SelectorName).Click()
		break
	case DoubleClick:
		findSelection(p.Selector, p.SelectorName).DoubleClick()
		break
	case Check:
		findSelection(p.Selector, p.SelectorName).Check()
	case Uncheck:
		findSelection(p.Selector, p.SelectorName).Uncheck()
	case Select:
		findSelection(p.Selector, p.SelectorName).Select(p.SelectorVal)
	case Submit:
		findSelection(p.Selector, p.SelectorName).Submit()
	case Fill:
		var text string
		if p.SelectorVal == "" {
			text = paramMap[p.Param]
		} else {
			text = p.SelectorVal
		}
		global.Log.Info("Fill：", findSelection(p.Selector, p.SelectorName).String(), p.Selector, p.SelectorName, text)
		findSelection(p.Selector, p.SelectorName).Fill(text)
		break
	case Text:
		result, _ := findSelection(p.Selector, p.SelectorName).Text()
		if p.Result != "" {
			if strings.Contains(result, p.Result) {
				global.Log.Info("执行结果：", result)
				return true, result, nil
			}
		}
	case RunScript:
		if p.JsParam != nil {
			for key, value := range p.JsParam {
				v := string(value.(string))
				if strings.Contains(v, "/v") {
					v = strings.Replace(v, "/v", "", 1)
					p.JsParam[key] = paramMap[v]
				}
				global.Log.Info("JsParam", key, p.JsParam[key])
			}
		}
		global.Log.Info("执行JS：", p.Js)
		global.Page.RunScript(p.Js, p.JsParam, &jsResult)
		break
	}
	if p.Sleep > 0 {
		time.Sleep(time.Duration(p.Sleep) * time.Second)
	}
	return true, p.Handle + "执行完成", nil
}

func findSelection(selector string, selectorName string) *agouti.Selection {
	switch selector {
	case FindByID:
		return global.Page.FindByID(selectorName)
	case FindByClass:
		return global.Page.FindByClass(selectorName)
	case FindByName:
		return global.Page.FindByName(selectorName)
	case FindByButton:
		return global.Page.FindByButton(selectorName)
	case FindByLabel:
		return global.Page.FindByLabel(selectorName)
	case FindByLink:
		return global.Page.FindByLink(selectorName)
	case FindByXPath:
		return global.Page.FindByXPath(selectorName)
	default:
		return global.Page.Find(selectorName)
	}
}

func checkLogin(domain string, identity string) (bool, string) {
	flog := false
	result := ""
	i := 0
	for {
		global.Log.Info("检测登陆状态")
		thisUrl, _ := global.Page.URL()
		u, _ := url.Parse(domain)
		if !strings.Contains(thisUrl, u.Host) {
			result = "已经不在本网站了，结束检测登陆状态"
			global.Log.Info(result)
			break
		}

		c, f := CheckIdentity(identity)
		flog = f

		if flog {
			cookieJson, _ := json.Marshal(c)

			global.Log.Info(string(cookieJson[:]))
			//global.RD.SetString(domain, cookieJson)
			dao.SetUnionCookies(domain, string(cookieJson), global.Acc.Id)
			result = "设置Cookie成功"
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

func setCookieLogin(doMan string) (bool, error) {
	cookieRd, err := dao.GetUnionCookies(doMan, global.Acc.Id) //global.RD.GetString(doMan)
	if err != nil {
		return false, err
	}
	if cookieRd == "" {
		return false, nil
	}
	global.Log.Info("获取Redis：" + cookieRd)
	var cks = []*http.Cookie{}
	err = json.Unmarshal([]byte(cookieRd), &cks)
	if err != nil {
		return false, err
	}
	global.Navigate(doMan)
	for i := range cks {
		cc := cks[i]
		global.Page.SetCookie(cc)
	}
	return true, nil
}

func CheckIdentity(identity string) ([]*http.Cookie, bool) {
	c, err := global.Page.GetCookies()
	if err != nil {
		return nil, false
	}
	for _, apiCookie := range c {
		if apiCookie.Name == identity && apiCookie.Value != "" {
			return c, true
		}
	}
	return nil, false
}

func main() {
	pf := &PushFillJson{
		"ID",
		"share-modal",
		"",
		"Click",
		1000,
		"alert(1)",
		map[string]interface{}{"hiveHtml": "hiveHtml", "host": "host"},
		"title",
		"已发布",
	}

	var pfs []PushFillJson

	pfs = append(pfs, *pf)
	pf = &PushFillJson{
		"ID2",
		"share-modal2",
		"123456",
		"Text",
		1000,
		"alert(1)",
		map[string]interface{}{"hiveHtml": "hiveHtml", "host": "host"},
		"title",
		"已发布",
	}
	pfs = append(pfs, *pf)
	st := PushJson{
		"简书",
		"http://www.jianshu.com",
		"https://www.jianshu.com/sign_in",
		"remember_user_token",
		"http://www.jianshu.com/writer#/",
		1000,
		pfs,
		nil,
	}

	b, err := json.Marshal(st)

	if err != nil {
		global.Log.Info("encoding faild")
	} else {
		j := string(b)
		var pj PushJson
		global.Log.Info(j)
		err := json.Unmarshal(b, &pj)
		if err != nil {
			global.Log.Info(err.Error())
		} else {
			global.Log.Info("结果：" + pj.Title)
		}
	}

}
