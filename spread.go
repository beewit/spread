package main

import (
	"fmt"
	//"time"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/parser"
	"github.com/beewit/spread/router"
	"github.com/sclevine/agouti"
)

func main() {
	start()
	router.Router()

}

func start() {
	global.Driver = agouti.ChromeDriver(agouti.ChromeOptions("args", []string{"--start-maximized", "--disable-infobars",
		"--app=http://www.jq22.com/demo/svgloader-150105194218/", "--webkit-text-size-adjust"}))
	global.Driver.Start()
	var err error
	global.Page, err = global.Driver.NewPage()
	if err != nil {
		fmt.Println("Failed to open page.")
	}
	global.Navigate(global.Host)
	paramMap := map[string]string{"title": "ArchLinux 以服务的方式启动Redis", "content": `启动redis服务
sudo systemctl restart redis

修改redis配置文件
nano /etc/redis.conf

测试密码正确性
127.0.0.1:6379> auth 123456`}
	rule := utils.Read(utils.JsonPath("parser", "zhihu.json"))
	flog, result, err2 := parser.RunPush(rule, paramMap)
	println(flog, result, err2)
}
