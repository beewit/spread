package main

import (
	"fmt"

	"github.com/beewit/spread/global"
	"github.com/beewit/spread/router"
	"github.com/sclevine/agouti"
)

func main() {
	start()
	router.Router()
}
func start() {
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
		global.Page.Navigate(global.Host)
	}()

}
