package main

import (
	"fmt"

	"runtime"

	"github.com/beewit/spread/global"
	"github.com/beewit/spread/handler"
	"github.com/beewit/spread/router"
	"github.com/sclevine/agouti"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	start()
	router.Router()
}
func start() {
	load := global.Host
	acc := handler.CheckClientLogin()
	if acc == nil {
		load = global.API_SSO_DOMAN + "?backUrl=" + global.Host + "/ReceiveToken"
	} else {
		global.Acc = acc
	}
	global.Driver = agouti.ChromeDriver(agouti.ChromeOptions("args", []string{
		"--gpu-process",
		"--in-process-gpu",
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
	go func() {
		global.Page.Navigate(load)
	}()
}
