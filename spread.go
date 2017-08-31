package main

import (
	"fmt"

	"github.com/beewit/spread/global"
	"github.com/beewit/spread/router"
	"github.com/sclevine/agouti"
	"github.com/beewit/spread/handler"
)

func main() {
	start()
	router.Router()
}
func start() {
	load := global.Host
	if !CheckClientLogin() {
		load = global.API_SSO_DOMAN + "?backUrl=" + global.Host + "/ReceiveToken"
	}
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
	go func() {
		global.Page.Navigate(load)
	}()
}

func CheckClientLogin() bool {
	token, err := global.QueryLoginToken()
	if err != nil {
		global.Log.Error(err.Error())
		panic(err)
	}
	if token == "" {
		return false
	}
	return handler.CheckClientLogin(token)
}
