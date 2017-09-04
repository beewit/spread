package global

import (
	"fmt"
	"time"

	"github.com/beewit/beekit/conf"
	"github.com/beewit/beekit/log"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/sqlite"
	"github.com/beewit/beekit/utils"
	"github.com/sclevine/agouti"
	"strings"
)

var (
	CFG      = conf.New("config.json")
	DB       = mysql.DB
	SLDB     = sqlite.DB
	Driver   *agouti.WebDriver
	Page     *agouti.Page
	HiveHtml = utils.Read("app/page/index.html")
	HiveJs   = utils.Read("app/static/js/inject.js")
	Log      = log.Logger
	IP       = CFG.Get("server.ip")
	Port     = CFG.Get("server.port")
	Host     = fmt.Sprintf("http://%v:%v", IP, Port)
	Navigate = PageNavigate
	Acc      *Account
	Platform = map[string]int{"新浪微博": 1, "简书": 2, "知乎": 3}
)

const API_SERVICE_DOMAN = "http://hive.tbqbz.com"
const API_SSO_DOMAN = "http://sso.tbqbz.com"

func injection() {
	time.Sleep(300 * time.Millisecond)
	arguments := map[string]interface{}{"hiveHtml": HiveHtml, "host": Host}
	//jquery
	//js := ";$(function () {$('body').append(hiveHtml)});" + HiveJs
	js := "var hiveHtmlDiv = document.createElement('div');hiveHtmlDiv.innerHTML=hiveHtml;document.body.appendChild(hiveHtmlDiv);" + HiveJs
	Page.RunScript(js, arguments, nil)
}

func PageAlertMsg(tip, url string) {
	js := fmt.Sprintf("alert('%v');localhost.href='%v'", tip, url)
	Page.RunScript(js, nil, nil)
}

func PageSuccessMsg(tip, url string) {
	PageMsg("#19a010", tip, url)
}

func PageErrorMsg(tip, url string) {
	PageMsg("#f33a3a", tip, url)
}

func PageMsg(status, tip, url string) {
	tipDiv := fmt.Sprintf(`<div style="
    position: fixed;
    width: 100%%;
    height: 100%%;
    background-color: rgba(0, 0, 0, 0.36);
    z-index: 999999;
    text-align: center;top:0"><span style="
    background-color: %s;
    padding: 20px 50px;
    color: #fff;
    line-height: 50px;
    font-size: 16px;
    border-radius: 5px;
    margin-top: 20px;
    top: 20px;
    font-weight: 900;">%s</span></div>`, status, tip)
	js := fmt.Sprintf("var div = document.createElement('div');div.innerHTML=`%v`;document.body.appendChild(div);setTimeout(function () {     location.href='%v';    },1500)", tipDiv, url)
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
