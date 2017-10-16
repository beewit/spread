package main

import (
	"fmt"
	"runtime"

	"github.com/beewit/spread/global"
	"github.com/beewit/spread/handler"
	"github.com/beewit/spread/router"
	"github.com/lxn/walk"
	"github.com/sclevine/agouti"

	"log"

	"github.com/beewit/beekit/utils"
)

type MyWindow struct {
	*walk.MainWindow
	ni *walk.NotifyIcon
}

func NewMyWindow() *MyWindow {
	mw := new(MyWindow)
	var err error
	mw.MainWindow, err = walk.NewMainWindow()
	checkError(err)
	return mw
}

func (mw *MyWindow) init() {
	utils.CloseChrome()
	go func() {
		err := Start()
		if err != nil {
			walk.MsgBox(mw, "工蜂小智-系统提示", "工蜂小智启动失败,错误："+err.Error(), walk.MsgBoxIconInformation)
		}
	}()
}

func (mw *MyWindow) RunHttpServer() error {
	return nil
}

func (mw *MyWindow) AddNotifyIcon() {
	var err error
	mw.ni, err = walk.NewNotifyIcon()
	checkError(err)
	mw.ni.SetVisible(true)

	icon, err := walk.NewIconFromResourceId(3)
	checkError(err)
	mw.SetIcon(icon)
	mw.ni.SetIcon(icon)

	mw.addAction(nil, "显示主界面").Triggered().Attach(func() {
		go func() {
			err := Start()
			if err != nil {
				walk.MsgBox(mw, "工蜂小智-系统提示", "工蜂小智启动失败,错误："+err.Error(), walk.MsgBoxIconInformation)
			}
		}()
	})

	mw.addAction(nil, "工蜂小智官网").Triggered().Attach(func() {
		global.Log.Info("打开工蜂小智官网")
		err := utils.Open("http://www.tbqbz.com/")
		if err != nil {
			global.Log.Error(err.Error())
		}
	})

	mw.addAction(nil, "联系我们").Triggered().Attach(func() {
		global.Log.Info("打开工蜂小智-联系我们")
		err := utils.Open("http://www.tbqbz.com/page/about/contact.html")
		if err != nil {
			global.Log.Error(err.Error())
		}
	})

	mw.addAction(nil, "安全退出").Triggered().Attach(func() {
		defer func() {
			mw.ni.Dispose()
			mw.Dispose()
			walk.App().Exit(0)
		}()
		Stop()
	})

}

func (mw *MyWindow) addMenu(name string) *walk.Menu {
	helpMenu, err := walk.NewMenu()
	checkError(err)
	help, err := mw.ni.ContextMenu().Actions().AddMenu(helpMenu)
	checkError(err)
	help.SetText(name)

	return helpMenu
}

func (mw *MyWindow) addAction(menu *walk.Menu, name string) *walk.Action {
	action := walk.NewAction()
	action.SetText(name)
	if menu != nil {
		menu.Actions().Add(action)
	} else {
		mw.ni.ContextMenu().Actions().Add(action)
	}

	return action
}

func (mw *MyWindow) msgbox(title, message string, style walk.MsgBoxStyle) {
	mw.ni.ShowInfo(title, message)
	walk.MsgBox(mw, title, message, style)
}

func main() {
	mw := NewMyWindow()
	mw.init()
	mw.AddNotifyIcon()
	mw.Run()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Start() error {
	go router.Router()
	runtime.GOMAXPROCS(runtime.NumCPU())
	load := global.Host
	acc := handler.CheckClientLogin()
	if acc == nil {
		load = global.API_SSO_DOMAN + "?backUrl=" + global.Host + "/ReceiveToken"
	} else {
		global.Acc = acc
	}
	global.Driver = agouti.ChromeDriver(agouti.ChromeOptions("args", []string{
		"--user-data-dir=ChromeUserData",
		"--gpu-process",
		"--start-maximized",
		"--disable-infobars",
		"--app=http://sso.tbqbz.com/",
		"--webkit-text-size-adjust"}))
	global.Driver.Start()
	var err error
	global.Page.Page, err = global.Driver.NewPage()
	if err != nil {
		global.Log.Error("Failed to open page.")
		return err
	}
	go func() {
		global.Log.Info(load)
		global.Page.Navigate(load)
	}()
	return nil
}

func Stop() {
	global.Log.Info("退出桌面应用")
	global.Page.Page.CloseWindow()
	global.Driver.Stop()
	global.Log.Info("退出服务")
	router.Stop()
}

func Show() {
	err := global.Page.Page.NextWindow()
	if err != nil {
		global.Page.Page, err = global.Driver.NewPage()
		if err != nil {
			fmt.Println("Failed to open page.")
		}
		go func() {
			global.Page.Navigate(global.Host)
		}()
	}
}
