package main

import (
	"fmt"
	"runtime"

	"github.com/beewit/spread/global"
	"github.com/beewit/spread/handler"
	"github.com/beewit/spread/router"
	"github.com/lxn/walk"

	"log"

	"github.com/beewit/beekit/utils"
	"github.com/sclevine/agouti"
	"time"
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
	mw.ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		global.Page.Page.NextWindow()
	})
	mw.addAction(nil, "显示主界面").Triggered().Attach(func() {
		go func() {
			err := Start()
			if err != nil {
				walk.MsgBox(mw, "工蜂小智-系统提示", "工蜂小智启动失败,错误："+err.Error(), walk.MsgBoxIconInformation)
			}
		}()
	})
	title := "工蜂小智-系统提示"
	taskMenu := mw.addMenu("进行中的任务")
	mWechatGroupAdd := mw.addAction(taskMenu, "停止添加微信群")
	mWechatGroupAdd.SetEnabled(false)
	mWechatGroupAdd.Triggered().Attach(func() {
		global.DelTask(global.TASK_WECHAT_ADD_GROUP)
		walk.MsgBox(mw, title, "关闭批量添加微信群成功，请等待本次流程完毕", walk.MsgBoxOK)
	})
	mWechatMessageSend := mw.addAction(taskMenu, "停止发送微信消息")
	mWechatMessageSend.SetEnabled(false)
	mWechatMessageSend.Triggered().Attach(func() {
		global.DelTask(global.TASK_WECHAT_SEND_MESSAGE)
		walk.MsgBox(mw, title, "关闭批量发送微信消息成功", walk.MsgBoxOK)
	})

	mWechatUserAdd := mw.addAction(taskMenu, "停止添加微信群成员")
	mWechatUserAdd.SetEnabled(false)
	mWechatUserAdd.Triggered().Attach(func() {
		global.DelTask(global.TASK_WECHAT_ADD_GROUP_USER)
		walk.MsgBox(mw, title, "关闭批量添加微信群成员成功", walk.MsgBoxOK)
	})

	mPlatformPush := mw.addAction(taskMenu, "停止发送平台内容")
	mPlatformPush.SetEnabled(false)
	mPlatformPush.Triggered().Attach(func() {
		global.DelTask(global.TASK_PLATFORM_PUSH)
		walk.MsgBox(mw, title, "关闭批量添加微信群成员成功，请等待本次流程完毕", walk.MsgBoxOK)
	})

	go func() {
		for {
			/*
				"TASK_PLATFORM_PUSH":         "平台自动化营销内容群发",
				"TASK_WECHAT_ADD_GROUP":      "自动化添加微信群",
				"TASK_WECHAT_SEND_MESSAGE":   "批量发送微信群或人的消息",
				"TASK_WECHAT_ADD_GROUP_USER": "自动化发起添加微信群成员"
			*/
			taskFlog := false
			task := global.GetTask(global.TASK_PLATFORM_PUSH)
			if task == nil || !task.State {
				if mPlatformPush.Enabled() {
					mPlatformPush.SetEnabled(false)
				}
			} else {
				taskFlog = true
				if !mPlatformPush.Enabled() {
					mPlatformPush.SetEnabled(true)
				}
			}

			task = global.GetTask(global.TASK_WECHAT_ADD_GROUP)
			if task == nil || !task.State {
				if mWechatGroupAdd.Enabled() {
					mWechatGroupAdd.SetEnabled(false)
				}
			} else {
				taskFlog = true
				if !mWechatGroupAdd.Enabled() {
					mWechatGroupAdd.SetEnabled(true)
				}
			}

			task = global.GetTask(global.TASK_WECHAT_SEND_MESSAGE)
			if task == nil || !task.State {
				if mWechatMessageSend.Enabled() {
					mWechatMessageSend.SetEnabled(false)
				}
			} else {
				taskFlog = true
				if !mWechatMessageSend.Enabled() {
					mWechatMessageSend.SetEnabled(true)
				}
			}

			task = global.GetTask(global.TASK_WECHAT_ADD_GROUP_USER)
			if task == nil || !task.State {
				if mWechatUserAdd.Enabled() {
					mWechatUserAdd.SetEnabled(false)
				}
			} else {
				taskFlog = true
				if !mWechatUserAdd.Enabled() {
					mWechatUserAdd.SetEnabled(true)
				}
			}
			if taskFlog {

			}
			time.Sleep(time.Second)
		}
	}()

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
	utils.CloseChrome()
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
		"--app=" + global.LoadPage,
		"--webkit-text-size-adjust"}))
	global.Driver.Start()
	var err error
	global.Page.Page, err = global.Driver.NewPage()
	if err != nil {
		global.Log.Error("Failed to open page.")
		return err
	}
	go func() {
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
