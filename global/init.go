package global

import (
	"fmt"
	"time"

	"strings"

	"database/sql"
	"encoding/json"

	"os"

	"github.com/astaxie/beego/logs"
	"github.com/beewit/beekit/sqlite"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread-update/update"
	"github.com/beewit/spread/static"
	"github.com/beewit/wechat-ai/smartQQ"
	"github.com/beewit/wechat-ai/smartWechat"
	"github.com/sclevine/agouti"
)

const (
	API_DOMAIN         = "http://www.9ee3.com"
	API_SERVICE_DOMAIN = "http://hive.9ee3.com"
	API_SSO_DOMAIN     = "http://sso.9ee3.com"
	SQLITE_DATABASE    = "app/spread.db"

	PAGE_SIZE   = 10
	FUNC_WECHAT = 6
	FUNC_QQ     = 7
	VERSION_DB  = 1
)

var (
	//先改版本，在编译后上传到gitee.com做版本维护
	//请注意，此版本不能大于https://gitee.com/beewit/spread/releases/new  上的版本
	Version          = update.Version{Major: 1, Minor: 0, Patch: 10}
	VersionStr       = fmt.Sprintf("V%d.%d.%d", Version.Major, Version.Minor, Version.Patch)
	SLDB             *sqlite.SqlConnPool
	Driver           *agouti.WebDriver
	Log              = logs.GetBeeLogger()
	IP               = "127.0.0.1"
	Port             = "8080"
	Host             = fmt.Sprintf("http://%v:%v", IP, Port)
	Navigate         = PageNavigate
	Acc              *Account
	Page             = *new(utils.AgoutiPage)
	WechatClientList = map[string]*smartWechat.WechatClient{}
	WechatUUid       = map[string]*smartWechat.WechatLoginStatus{}
	WechatClient     *smartWechat.WechatClient
	QQClientList     = map[int64]*smartQQ.QQClient{}
	QQClient         = smartQQ.NewQQClient(&smartQQ.QQClient{})
	TaskList         = map[string]*Task{}
	VoiceSwitch      = true
	LoadPage         = API_DOMAIN + "/page/load.html"
	ContactPage      = API_DOMAIN + "/page/about/contact.html"
	HiveHtml         string
	HiveJs           string
)

func InitGlobal() {
	initLog()
	CheckSqliteDB()
	err := CheckUpdateDB()
	if err != nil {
		Log.Error(err.Error())
	}
	iniSqliteDB()
}

func CheckSqliteDB() {
	var flog bool
	var err error
	flog, err = utils.PathExists(SQLITE_DATABASE)
	if !flog {
		//创建数据库
		var file *os.File
		file, err = utils.CreateFile(SQLITE_DATABASE)
		if err != nil {
			Log.Error(err.Error())
			panic(err)
		}
		file.Write(static.FileAppSpreadDb)
		file.Close()
	}
}

func iniSqliteDB() {
	var err error
	SLDB = &sqlite.SqlConnPool{
		DriverName:     "sqlite3",
		DataSourceName: SQLITE_DATABASE,
	}
	SLDB.SqlDB, err = sql.Open(SLDB.DriverName, SLDB.DataSourceName)
	if err != nil {
		Log.Error(err.Error())
		panic(err)
		return
	}
}

func initLog() {
	conf := fmt.Sprintf(
		`{
			"filename": "%s",
			"maxdays": %s,
			"daily": %s,
			"rotate": %s,
			"level": %s,
			"separate": "[%s]"
		}`,
		"spread.log",
		"10",
		"true",
		"true",
		"7",
		"error",
	)
	logs.SetLogger(logs.AdapterMultiFile, conf)
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
}

func injection() {
	if HiveHtml == "" {
		HiveHtml = string(static.FileAppPageIndexHTML) //utils.Read("app/page/index.html")
	}
	if HiveJs == "" {
		HiveJs = string(static.FileAppStaticJsInjectJs) //utils.Read("app/static/js/inject.js")
	}
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

func Logs(errStr string) {
	errStr = time.Now().Format("2006-01-02 15:04:05") + "   " + errStr
	file, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND, 0x644)
	defer file.Close()
	if err != nil {
		println(errStr)
	} else {
		file.Write([]byte(errStr))
	}
}
