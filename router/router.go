package router

import (
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/handler"

	"fmt"

	//	"time"

	"github.com/labstack/echo"
	"github.com/sclevine/agouti"
)

func Start() {
	ip := global.CFG.Get("server.ip")
	port := global.CFG.Get("server.port")
	host := fmt.Sprintf("http://%s:%s", ip, port)

	e := echo.New()
	e.Static("/app", "app")
	e.File("/", "app/page/index.html")
	handlerConfig(e)

	global.Driver = agouti.ChromeDriver(agouti.ChromeOptions("args", []string{
		"--start-maximized",
		"--disable-infobars",
		"--app=http://www.jq22.com/demo/svgloader-150105194218/",
		"--webkit-text-size-adjust"}))
	global.Driver.Start()
	var err error
	global.Page, err = global.Driver.NewPage()
	if err != nil {
		fmt.Println("Failed to open page.")
	}
	go run(host)
	go e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))

	println("ceshi")
}

func handlerConfig(e *echo.Echo) {
	e.POST("/auth/identity", handler.Identity)
}

func run(host string) {

	global.Page.Navigate("http://www.jianshu.com")

	arguments := map[string]interface{}{"hiveHtml": global.HiveHtml, "host": host}

	js := ";$(function () {$('body').append(hiveHtml)});" + global.HiveJs

	global.Page.RunScript(js, arguments, nil)

}
