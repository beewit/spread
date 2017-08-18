package global

import (
	"github.com/beewit/beekit/conf"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/redis"
	"github.com/sclevine/agouti"
)

var (
	CFG    = conf.New("config.json")
	DB     = mysql.DB
	RD     = redis.Cache
	Driver *agouti.WebDriver
	Page   *agouti.Page
)
