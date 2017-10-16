package parser

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"net/url"

	"errors"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/dao"
	"github.com/beewit/spread/global"
	"github.com/sclevine/agouti"
	"math/rand"
)

type PushJson struct {
	Title        string            `json:"title"`
	Domain       string            `json:"domain"`
	LoginUrl     string            `json:"loginUrl"`
	Identity     string            `json:"identity"`
	WriterUrl    []string          `json:"writerUrl"`
	Sleep        int64             `json:"sleep"`
	Fill         []PushFillJson    `json:"fill"`
	Login        []PushFillJson    `json:"login"`
	SwitchIframe map[string]string `json:"switchIframe"`
	TimeOut      int               `json:"timeOut"`
}

type PushFillJson struct {
	Selector     string                 `json:"selector"`
	SelectorName string                 `json:"selectorName"`
	SelectorVal  string                 `json:"selectorVal"`
	Attr         string                 `json:"attr"`
	Handle       string                 `json:"handle"`
	Sleep        int64                  `json:"sleep"`
	Js           string                 `json:"js"`
	JsParam      map[string]interface{} `json:"jsParam"`
	Param        string                 `json:"param"`
	Result       string                 `json:"result"`
	Field        string                 `json:"field"`
	Status       string                 `json:"status"`
	SwitchIframe string                 `json:"switchIframe"`
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
	PageURL     = "PageURL"
	Attr        = "Attr"
	Value       = "Value"
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

const (
	Executing_Status = "Executing"
	Complete_Status  = "Complete"
)

func GetPushJson(rule string) (*PushJson, error) {
	var pj PushJson
	err := json.Unmarshal([]byte(rule), &pj)
	if err != nil {
		global.Log.Error("规则：%s", rule)
		global.Log.Error("规则解析错误：%s", err.Error()) 
		return &PushJson{}, err
	} else {
		return &pj, nil
	}
}

func GetSwitchIframe(m map[string]string, key string) string {
	if m != nil {
		return m[key]
	}
	return ""
}

func SwitchIframe(iframeList string) (bool, int, error) {
	if iframeList != "" {
		var iframes []string
		if iframeList != "" {
			iframes = strings.Split(iframeList, "|")
		}
		if len(iframes) > 0 {
			for i := 0; i < len(iframes); i++ {
				iframeSelector := iframes[i]
				if iframeSelector != "" {
					time.Sleep(time.Second * 1)
					iframe, err := global.Page.Find(iframeSelector).Elements()
					if err != nil {
						global.Log.Error(err.Error())
						return false, 0, err
					}
					if len(iframe) <= 0 {
						return false, 0, errors.New("未查找到iframe，Selector：" + iframeSelector)
					}
					err = global.Page.SwitchToRootFrameByName(iframe[0])
					if err != nil {
						global.Log.Error(err.Error())
						return false, 0, err
					}
					return true, len(iframes), nil
				}

			}
		}
	}
	return false, 0, nil
}

func RunPush(rule string, paramMap map[string]string, platformAcc string, platformId int64, switchAccount bool) (bool, bool, string, map[string]interface{}, error) {
	resultMap := map[string]interface{}{}
	pj, err := GetPushJson(rule)
	if err != nil {
		return false, false, pj.Title + "解析配置规则失败", nil, err
	}
	if len(pj.Fill) <= 0 {
		return false, false, pj.Title + "无发布规则配置", nil, nil
	}
	var flog bool
	if pj.Login != nil && len(pj.Login) > 0 {
		global.Navigate(pj.LoginUrl)
		if switchAccount {
			DeleteCookie()
		}
		_, flog = CheckIdentity(pj.Identity)
		if !flog || switchAccount {
			flog, _ = SetCookieLogin(platformId, pj.Domain, platformAcc, time.Duration(pj.TimeOut))
			_, flog2 := CheckIdentity(pj.Identity)
			if !flog || !flog2 || switchAccount {
				if pj.LoginUrl != global.PageUrl() {
					global.Navigate(pj.LoginUrl)
					time.Sleep(time.Second)
				}
				iframe, ic, err := SwitchIframe(GetSwitchIframe(pj.SwitchIframe, "login"))
				if err != nil {
					global.Log.Error(err.Error())
					url, _ := global.Page.URL()
					return false, false, "切换Iframe失败，PageURL：" + url, nil, nil
				}
				global.Log.Info("自动登陆中...")
				for i := 0; i < len(pj.Login); i++ {
					if CheckStopAtSite(pj.Domain) {
						return false, false, "已经不在任务网站了，结束任务执行", nil, nil
					}
					rFlog, result, _, _, err := HandleSelection(&pj.Login[i], paramMap)

					if err != nil {
						global.Log.Error(err.Error())
					}
					rule, _ := json.Marshal(pj.Login[i])
					global.Log.Warning("《登陆》执行任务：%s", string(rule))
					global.Log.Warning("《登陆》执行结果：%v，返回之：result：%s", rFlog, result)
				}
				flog, _ = CheckLogin(platformId, pj.Domain, pj.Identity, platformAcc)
				if iframe {
					for i := 0; i < ic; i++ {
						global.Page.SwitchToParentFrame()
					}
				}
			}
		}
	} else {
		global.Navigate(pj.Domain)
		_, flog = CheckIdentity(pj.Identity)
		if !flog {
			flog, _ = SetCookieLogin(platformId, pj.Domain, platformAcc, time.Duration(pj.TimeOut))
			_, flog2 := CheckIdentity(pj.Identity)
			if !flog || !flog2 {
				//等待登陆方式
				global.Navigate(pj.LoginUrl)
			} else {
				//设置Cookie
				global.Log.Info("设置Cookie成功")
			}
			flog, _ = CheckLogin(platformId, pj.Domain, pj.Identity, platformAcc)
		} else {
			global.Log.Info("已经是登陆状态")
		}
	}
	if !flog {
		return false, false, pj.Title + "登陆失败", nil, nil
	}
	time.Sleep(time.Second)
	sendUrl := pj.WriterUrl[rand.Intn(len(pj.WriterUrl))]
	global.Log.Warning("正在进入发送URL地址：v%", sendUrl)
	global.Navigate(sendUrl)
	if pj.Sleep > 0 {
		time.Sleep(time.Duration(pj.Sleep) * time.Second)
	}

	completFlog := false
	for i := 0; i < len(pj.Fill); i++ {
		iframe, ic, err := SwitchIframe(pj.Fill[i].SwitchIframe)
		if err != nil {
			global.Log.Error(err.Error())
			return false, false, "切换Iframe失败" + global.PageUrl(), nil, nil
		}
		if CheckStopAtSite(pj.Domain) {
			return false, false, "已经不在任务网站了，结束任务执行", nil, nil
		}
		rFlog, result, m, status, err := HandleSelection(&pj.Fill[i], paramMap)
		if err != nil {
			global.Log.Error(err.Error())
		} else {
			if status == Complete_Status {
				completFlog = rFlog
			}
			if rFlog && m != nil {
				for k, v := range m {
					resultMap[k] = v
				}
			}
		}
		if iframe {
			for i := 0; i < ic; i++ {
				global.Page.SwitchToParentFrame()
			}
		}
		rule, _ := json.Marshal(pj.Fill[i])
		global.Log.Warning("执行任务：%s", string(rule))
		global.Log.Warning("执行结果：%v，返回之：result：%s", rFlog, result)
	}

	global.Log.Warning("直接最终结果：%v", completFlog)
	return true, completFlog, "全部执行完成", resultMap, nil
}

func HandleSelection(p *PushFillJson, paramMap map[string]string) (bool, string, map[string]string, string, error) {
	if paramMap != nil && p.SelectorName != "" {
		for k, v := range paramMap {
			if strings.Contains(p.SelectorName, k+"/v") {
				p.SelectorName = strings.Replace(p.SelectorName, k+"/v", v, -1)
			}
		}
	}
	resultMap := map[string]string{}
	var result string
	var err error
	switch p.Handle {
	case Click:
		err = findSelection(p.Selector, p.SelectorName).Click()
		break
	case DoubleClick:
		err = findSelection(p.Selector, p.SelectorName).DoubleClick()
		break
	case Check:
		err = findSelection(p.Selector, p.SelectorName).Check()
	case Uncheck:
		err = findSelection(p.Selector, p.SelectorName).Uncheck()
	case Select:
		err = findSelection(p.Selector, p.SelectorName).Select(p.SelectorVal)
	case Submit:
		err = findSelection(p.Selector, p.SelectorName).Submit()
	case Fill:
		var text string
		if p.SelectorVal == "" {
			text = paramMap[p.Param]
		} else {
			text = p.SelectorVal
		}
		err = findSelection(p.Selector, p.SelectorName).Fill(text)
		break
	case Text:
		result, err = findSelection(p.Selector, p.SelectorName).Text()
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
		global.Log.Info("执行JS：%s", p.Js)
		err = global.Page.RunScript(p.Js, p.JsParam, &result)
		break
	case PageURL:
		result, err = global.Page.URL()
		break
	case Attr:
		result, err = findSelection(p.Selector, p.SelectorName).Attribute(p.Attr)
		break
	case Value:
		result = p.Result
		break
	}
	if p.Sleep > 0 {
		time.Sleep(time.Duration(p.Sleep) * time.Second)
	}
	if err == nil {
		//直接产生的结果比较
		if p.Result != "" {
			oldResult := p.Result
			if strings.Contains(p.Result, "/v") {
				oldResult = paramMap[strings.Replace(p.Result, "/v", "", 1)]
			}
			if !strings.Contains(result, oldResult) {
				err = errors.New("结果值预期不符合，结果：" + result + ",预期：" + p.Result)
			}
		}
		//等待结果，如等待手动输入验证码

		if p.Field != "" {
			resultMap[p.Field] = result
		}
	}
	if err != nil {
		global.Log.Error(err.Error())
	}
	return err == nil, result, resultMap, p.Status, nil
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

func CheckLogin(platformId int64, domain, identity, platformAcc string) (bool, string) {
	flog := false
	result := ""
	i := 0
	for {
		global.Log.Info("检测登陆状态")
		if CheckStopAtSite(domain) {
			result = "已经不在本网站了，结束检测登陆状态"
			global.Log.Info(result)
			break
		}

		c, f := CheckIdentity(identity)
		flog = f

		if flog {

			cookieJson, _ := json.Marshal(c)

			global.Log.Info(string(cookieJson[:]))

			localStorage, err := global.PageLocalStorage()
			if err != nil {
				global.Log.Error("PageLocalStorage：" + err.Error())
			}
			session := "" //global.Page.Session()
			dao.SetUnionCookies(domain, string(cookieJson), localStorage, session, platformId, global.Acc.Id, platformAcc)
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

func CheckStopAtSite(domain string) bool {
	thisUrl, _ := global.Page.URL()
	u, _ := url.Parse(domain)
	domainName := utils.Substr(u.Host, strings.LastIndex(utils.Substr(u.Host, 0, strings.LastIndex(u.Host, ".")), "."), len(u.Host))
	global.Log.Info("原网站：%s,现网站：%s", thisUrl, domainName)
	return !strings.Contains(thisUrl, domainName)
}

func SetCookieLogin(platformId int64, doMan, platformAcc string, timeOut time.Duration) (bool, error) {
	cookieRd, localStorage, sessions, err := dao.GetUnionCookies(platformId, global.Acc.Id, platformAcc, timeOut) //global.RD.GetString
	// (doMan)
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

	global.PageAddLocalStorage(localStorage)

	global.PageAddSessionStorage(sessions)

	if timeOut == -3 {
		return false, nil
	} else {
		return true, nil
	}
}

func DeleteCookie() bool {
	err := global.Page.ClearCookies()
	if err != nil {
		return false
	}
	if err := global.Page.RunScript("localStorage.clear();", nil, nil); err != nil {
		return false
	}

	if err := global.Page.RunScript("sessionStorage.clear();", nil, nil); err != nil {
		return false
	}
	return true
}

func CheckIdentity(identity string) ([]*http.Cookie, bool) {
	c, err := global.Page.GetCookies()
	if err != nil {
		return nil, false
	}
	for _, apiCookie := range c {
		if strings.Contains(identity, "@") {
			idents := strings.Split(identity, "@")
			if len(idents) > 1 {
				if apiCookie.Name == idents[0] && apiCookie.Value == idents[1] {
					return c, true
				}
			}
		} else {
			if apiCookie.Name == identity && apiCookie.Value != "" {
				return c, true
			}
		}
	}
	return nil, false
}

func main() {
	pf := &PushFillJson{
		"ID",
		"share-modal",
		"",
		"",
		"Click",
		1000,
		"alert(1)",
		map[string]interface{}{"hiveHtml": "hiveHtml", "host": "host"},
		"title",
		"已发布",
		"",
		"",
		"",
	}

	var pfs []PushFillJson

	pfs = append(pfs, *pf)
	pf = &PushFillJson{
		"ID2",
		"share-modal2",
		"123456",
		"",
		"Text",
		1000,
		"alert(1)",
		map[string]interface{}{"hiveHtml": "hiveHtml", "host": "host"},
		"title",
		"已发布",
		"",
		"",
		"",
	}
	pfs = append(pfs, *pf)
	st := PushJson{
		"简书",
		"http://www.jianshu.com",
		"https://www.jianshu.com/sign_in",
		"remember_user_token",
		[]string{"http://www.jianshu.com/writer#/"},
		1000,
		pfs,
		nil,
		nil,
		0,
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
