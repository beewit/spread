package global

import (
	"fmt"
	"time"

	"strings"

	"encoding/json"
	"github.com/beewit/beekit/conf"
	"github.com/beewit/beekit/log"
	"github.com/beewit/beekit/sqlite"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/wechat-ai/smartQQ"
	"github.com/beewit/wechat-ai/smartWechat"
	"github.com/sclevine/agouti"
)

const (
	API_DOMAIN         = "http://www.tbqbz.com:8080"
	API_SERVICE_DOMAIN = "http://hive.tbqbz.com:8082"
	API_SSO_DOMAIN     = "http://sso.tbqbz.com:8081"
	PAGE_SIZE          = 10
	FUNC_WECHAT        = 6
	FUNC_QQ            = 7
)

var (
	CFG          = conf.New("config.json")
	SLDB         = sqlite.DB
	Driver       *agouti.WebDriver
	HiveHtml     = utils.Read("app/page/index.html")
	HiveJs       = utils.Read("app/static/js/inject.js")
	Log          = log.Logger
	IP           = CFG.Get("server.ip")
	Port         = CFG.Get("server.port")
	Host         = fmt.Sprintf("http://%v:%v", IP, Port)
	Navigate     = PageNavigate
	Acc          *Account
	Page         = *new(utils.AgoutiPage)
	WechatClient *smartWechat.WechatClient
	QQClient     = smartQQ.NewQQClient(&smartQQ.QQClient{})
	TaskList     = map[string]*Task{}
	VoiceSwitch  = true
	LOAD_PAGE    = API_DOMAIN + "/page/load.html"
	CONTACT_PAGE = API_DOMAIN + "/page/about/contact.html"
)

func injection() {
	time.Sleep(300 * time.Millisecond)
	arguments := map[string]interface{}{"hiveHtml": HiveHtml, "host": Host}
	js := "var hiveHtmlDiv = document.createElement('div');hiveHtmlDiv.innerHTML=hiveHtml;document.body.appendChild(hiveHtmlDiv);" + HiveJs
	Page.RunScript(js, arguments, nil)
}

func PageAlertMsg(tip, url string) {
	js := fmt.Sprintf("alert('%v');localhost.href='%v'", tip, url)
	Page.RunScript(js, nil, nil)
}

func PageSuccessMsg(tip, url string) {
	PageJumpMsg("#19a010", tip, url)
}

func PageErrorMsg(tip, url string) {
	PageJumpMsg("#f33a3a", tip, url)
}

func PageMsg(tip string) {
	PageJumpMsg("#ffb12c", tip, "")
}

func PageJumpMsg(status, tip, url string) {
	tipDiv := fmt.Sprintf(`<div id="pageMsg" style="
    position: fixed;
    width: 100%%;
    height: 100%%;
    background-color: rgba(0, 0, 0, 0.36);
    z-index: 999999998;
    text-align: center;top:0;">
	<span style="
    background-color: %s;
    padding: 20px 50px;
    color: #fff;
    line-height: 50px;
    font-size: 16px;
    border-radius: 5px;
    margin-top: 20px;
    top: 20px;
    font-weight: 900;position: relative;"
	onclick="var pageMsg= document.getElementById('pageMsg');pageMsg.parentNode.removeChild(pageMsg);">%s
	<a style="position: absolute;
    right: 4px;
    border-radius: 50%%;
    background-color: #fff;
    color: #464545;
    font-size: 12px;
    height: 40px;
    width: 40px;
    line-height: 40px;
    top: 8px;
    cursor: pointer;">关闭</a></span></div>`, status, tip)
	urls := ""
	if url != "" {
		if strings.Index(url, "http") == -1 {
			url = Host + "?lastUrl=" + url
		}
		urls = fmt.Sprintf("setTimeout(function () {     location.href='%v';    },1500)", url)
	}
	js := fmt.Sprintf("var pageMsg = document.getElementById('pageMsg'); if(pageMsg!=null) pageMsg.parentNode.removeChild(pageMsg);var div = document.createElement('div');div.innerHTML=`%v`;document.body.appendChild(div);%s", tipDiv, urls)
	Page.RunScript(js, nil, nil)
}

func PageNavigate(url string) {
	Page.Navigate(url)
	go injection()
}

func PageFindValue(selector string) string {
	if strings.Contains(selector, "@") {
		str := strings.Split(selector, "@")
		return PageFindAttr(str[0], str[1])
	}
	txt, elsErr := Page.Find(selector).Text()
	if elsErr != nil {
		Log.Error(elsErr.Error())
		return ""
	}
	return txt
}

func PageFindAttr(selector, attr string) string {
	els, elsErr := Page.Find(selector).Elements()
	if elsErr != nil {
		Log.Error(elsErr.Error())
		return ""
	}
	if len(els) > 0 {
		val, _ := els[0].GetAttribute(attr)
		return val
	}
	return ""
}

func PageUrl() string {
	url, _ := Page.URL()
	return url
}

func PageLocalStorage() (string, error) {
	var result string
	err := Page.RunScript("return JSON.stringify(localStorage);", nil, &result)
	return result, err
}

func PageAddLocalStorage(ls string) bool {
	if ls == "" {
		return false
	}
	m := map[string]string{}
	err := json.Unmarshal([]byte(ls), &m)
	if err != nil {
		Log.Error("json转换失败：" + ls)
		return false
	}
	for k, v := range m {
		arguments := map[string]interface{}{"key": k, "value": v}
		err = Page.RunScript("localStorage.setItem(key,value)", arguments, nil)
		if err != nil {
			Log.Error(fmt.Sprintf("localStorage.setItem('%s','%s')失败", k, v))
		} else {
			Log.Info("localStorage.setItem('%s','%s')成功", k, v)
		}
	}
	return true
}

func PageSessionStorageByKey(key string) (string, error) {
	var result string
	arguments := map[string]interface{}{"key": key}
	err := Page.RunScript("return sessionStorage.getItem(key);", arguments, &result)
	return result, err
}

func PageSessionStorage() (string, error) {
	var result string
	err := Page.RunScript("return JSON.stringify(sessionStorage);", nil, &result)
	return result, err
}

func PageAddSessionStorage(ss string) bool {
	if ss == "" {
		return false
	}
	m := map[string]string{}
	err := json.Unmarshal([]byte(ss), &m)
	if err != nil {
		Log.Error("json转换失败：" + ss)
		return false
	}
	for k, v := range m {
		arguments := map[string]interface{}{"key": k, "value": v}
		err = Page.RunScript("sessionStorage.setItem(key,value)", arguments, nil)
		if err != nil {
			Log.Error(fmt.Sprintf("sessionStorage.setItem('%s','%s')失败", k, v))
		} else {
			Log.Info("sessionStorage.setItem('%s','%s')成功", k, v)
		}
	}
	return true
}
