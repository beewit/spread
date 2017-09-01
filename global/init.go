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
	Platform    = map[string]int{"新浪微博": 1, "简书": 2, "知乎": 3}
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

func PageNavigate(url string) {
	Page.Navigate(url)
	go injection()
}
