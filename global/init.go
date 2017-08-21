package global

import (
	"github.com/beewit/beekit/conf"
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
)
