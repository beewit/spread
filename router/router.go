package router

import (
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/handler"

	"fmt"

	//	"time"

	"github.com/labstack/echo"
	"github.com/sclevine/agouti"
)

func Init() {
	ip := global.CFG.Get("server.ip")
	port := global.CFG.Get("server.port")
	host := fmt.Sprintf("http://%s:%s", ip, port)

	e := echo.New()
	e.Static("/app", "app")
	e.File("/", "app/page/index.html")
	handlerConfig(e)

	println("端口" + host)
	//host = "http://www.jianshu.com"
	global.Driver = agouti.ChromeDriver(agouti.ChromeOptions("args", []string{"--start-maximized", "--disable-infobars", "--app=" + host}))
	global.Driver.Start()
	var err error
	global.Page, err = global.Driver.NewPage()
	if err != nil {
		fmt.Println("Failed to open page.")
	}
	go run()
	go e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))

	println("ceshi")
}

func handlerConfig(e *echo.Echo) {
	e.POST("/auth/identity", handler.Identity)
}

func run() {
	//global.Page.Navigate("http://www.baidu.com")
}
