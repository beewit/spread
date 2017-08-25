package global

import (
	"fmt"
	"time"

	"github.com/beewit/beekit/conf"
	"github.com/beewit/beekit/log"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/redis"
	"github.com/beewit/beekit/utils"
	"github.com/sclevine/agouti"
)

var (
	CFG      = conf.New("config.json")
	DB       = mysql.DB
	RD       = redis.Cache
	Driver   *agouti.WebDriver
	Page     *agouti.Page
	HiveHtml = utils.Read("app/page/index.html")
	HiveJs   = utils.Read("app/static/js/inject.js")
	Log      = log.Logger
	IP       = CFG.Get("server.ip")
	Port     = CFG.Get("server.port")
	Host     = fmt.Sprintf("http://%s:%s", IP, Port)
	Navigate = PageNavigate
)

func injection() {
	time.Sleep(300 * time.Millisecond)
	arguments := map[string]interface{}{"hiveHtml": HiveHtml, "host": Host}
	//js := ";$(function () {$('body').append(hiveHtml)});" + HiveJs
	//js := ";document.getElementsByTagName(\"body\")[0].innerHTML += hiveHtml;" + HiveJs
	js := "var hiveHtmlDiv = document.createElement('div');hiveHtmlDiv.innerHTML=hiveHtml;document.body.appendChild(hiveHtmlDiv);" + HiveJs
	Page.RunScript(js, arguments, nil)
}

func PageNavigate(url string) {
	Page.Navigate(url)
	go injection()
}
