package test

import (
	"github.com/sclevine/agouti"
	"testing"
)

func TestApp(t *testing.T) {

	driver := agouti.ChromeDriver(agouti.ChromeOptions("args", []string{
		"--start-maximized",
		"--disable-infobars",
		"--app=https://i.qq.com/?rd=1",
		"--webkit-text-size-adjust"}))
	driver.Start()
	page, err2 := driver.NewPage()
	if err2 != nil {
		println(err2.Error())
		return
	}
	page.Navigate("http://www.baidu.com")
}
