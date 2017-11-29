package main

import (
	"fmt"
	"runtime"

	"github.com/beewit/spread/global"
	"github.com/beewit/spread/handler"
	"github.com/beewit/spread/router"
	"github.com/lxn/walk"

	"log"

	"time"

	"sync"

	"os"

	"errors"
	"io/ioutil"
	"strconv"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread-update/update"
	"github.com/beewit/wechat-ai/ai"
	"github.com/sclevine/agouti"
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
	go func() {
		err := Start()
		if err != nil {
			walk.MsgBox(mw, "工蜂小智-系统提示", "工蜂小智启动失败,错误："+err.Error(), walk.MsgBoxIconInformation)
			Stop()
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
		go mw.openSpread()
	})
	mw.addAction(nil, "显示主界面").Triggered().Attach(func() {
		go mw.openSpread()
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

	mQQGroupAdd := mw.addAction(taskMenu, "停止添加QQ群")
	mQQGroupAdd.SetEnabled(false)
	mQQGroupAdd.Triggered().Attach(func() {
		global.DelTask(global.TASK_QQ_ADD_GROUP)
		walk.MsgBox(mw, title, "关闭批量添加QQ群成功，请等待本次流程完毕", walk.MsgBoxOK)
	})

	mQQFriendAdd := mw.addAction(taskMenu, "停止添加QQ好友")
	mQQFriendAdd.SetEnabled(false)
	mQQFriendAdd.Triggered().Attach(func() {
		global.DelTask(global.TASK_QQ_ADD_FRIEND)
		walk.MsgBox(mw, title, "关闭批量添加QQ好友成功，请等待本次流程完毕", walk.MsgBoxOK)
	})

	mQQMessageSend := mw.addAction(taskMenu, "停止发送QQ消息")
	mQQMessageSend.SetEnabled(false)
	mQQMessageSend.Triggered().Attach(func() {
		global.DelTask(global.TASK_QQ_SEND_MESSAGE)
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
		taskFlog := false
		time.Sleep(time.Second * 5)
		for {
			/*
				"TASK_PLATFORM_PUSH":         "平台自动化营销内容群发2",
				"TASK_WECHAT_ADD_GROUP":      "自动化添加微信群",
				"TASK_WECHAT_SEND_MESSAGE":   "批量发送微信群或人的消息",
				"TASK_WECHAT_ADD_GROUP_USER": "自动化发起添加微信群成员"
			*/
			taskFlog = false
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

			task = global.GetTask(global.TASK_QQ_SEND_MESSAGE)
			if task == nil || !task.State {
				if mQQMessageSend.Enabled() {
					mQQMessageSend.SetEnabled(false)
				}
			} else {
				taskFlog = true
				if !mQQMessageSend.Enabled() {
					mQQMessageSend.SetEnabled(true)
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

			task = global.GetTask(global.TASK_QQ_ADD_GROUP)
			if task == nil || !task.State {
				if mQQGroupAdd.Enabled() {
					mQQGroupAdd.SetEnabled(false)
				}
			} else {
				taskFlog = true
				if !mQQGroupAdd.Enabled() {
					mQQGroupAdd.SetEnabled(true)
				}
			}

			task = global.GetTask(global.TASK_QQ_ADD_FRIEND)
			if task == nil || !task.State {
				if mQQFriendAdd.Enabled() {
					mQQFriendAdd.SetEnabled(false)
				}
			} else {
				taskFlog = true
				if !mQQFriendAdd.Enabled() {
					mQQFriendAdd.SetEnabled(true)
				}
			}

			if taskFlog && global.VoiceSwitch {
				openVoice()
			} else {
				closeVoice()
			}
			time.Sleep(time.Second)
		}
	}()
	mw.addAction(nil, "工蜂小智官网").Triggered().Attach(func() {
		global.Log.Info("打开工蜂小智官网")
		err := utils.Open(global.API_DOMAIN)
		if err != nil {
			global.Log.Error(err.Error())
		}
	})

	mw.addAction(nil, "联系我们").Triggered().Attach(func() {
		global.Log.Info("打开工蜂小智-联系我们")
		err := utils.Open(global.ContactPage)
		if err != nil {
			global.Log.Error(err.Error())
		}
	})

	mVoiceSwitch := mw.addAction(nil, "关闭声音")
	mVoiceSwitch.Triggered().Attach(func() {
		if global.VoiceSwitch {
			mVoiceSwitch.SetText("打开声音")
			global.VoiceSwitch = false
		} else {
			mVoiceSwitch.SetText("关闭声音")
			global.VoiceSwitch = true
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

var syncMutex *sync.Mutex

func (mw *MyWindow) openSpread() {
	syncMutex.Lock()
	global.Log.Info("正在打开工蜂小智主界面")
	defer func() {
		if err := recover(); err != nil {
			global.Log.Info("新打开工蜂小智主界面")
			global.Page.Page, err = global.Driver.NewPage()
			if err != nil {
				walk.MsgBox(mw, "工蜂小智-系统提示", "工蜂小智启动失败，请重新启动程序,错误："+convert.ToString(err), walk.MsgBoxIconInformation)
			} else {
				global.Page.Page.Navigate(global.Host)
			}
			syncMutex.Unlock()
		}
	}()
	global.Page.Page.NextWindow()
	title, err := global.Page.Title()
	if err != nil {
		global.Log.Error("global.Page.Title ERROR:" + err.Error())
		panic(err)
	} else {
		ai.ForegroundWindow("Chrome_WidgetWin_1", title)
	}
	global.Log.Info("已打开工蜂小智主界面")
	syncMutex.Unlock()
}

func openVoice() {
	if global.Page.Page != nil {
		js := fmt.Sprintf("var tipMp3=document.getElementById('hive_tip');if(tipMp3==null){"+
			"var tipMp3 = document.createElement('audio');"+
			"tipMp3.id='hive_tip';"+
			"tipMp3.loop='loop';"+
			"tipMp3.src='%s/app/static/media/hive_tip.mp3';"+
			"tipMp3.autoplay='autoplay';"+
			"tipMp3.volume = 0.2;"+
			"document.body.appendChild(tipMp3);}", global.Host)
		if err := global.Page.RunScript(js, nil, nil); err != nil {
			global.Log.Error("播放声音，错误：%s", err.Error())
		}
	}
}

func closeVoice() {
	if global.Page.Page != nil {
		js := "var audio=document.getElementsByTagName('audio');for(var i = 0;i<(audio.length) * 2;i++){audio[0].remove();}"
		if err := global.Page.RunScript(js, nil, nil); err != nil {
			global.Log.Error("播放声音，错误：%s", err.Error())
		}
	}
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

func main() {
	defer func() {
		if err := recover(); err != nil {
			errStr := fmt.Sprintf("《程序出现严重错误，终止运行！》，ERROR：%v", err)
			global.Logs(errStr)
		}
	}()
	mw := NewMyWindow()
	iManPid := fmt.Sprint(os.Getpid())
	if err := ProcExsit(); err == nil {
		pidFile, _ := os.Create("pid.pid")
		defer pidFile.Close()
		pidFile.WriteString(iManPid)

	} else {
		str := fmt.Sprintf("*************************【%s】*************************", err.Error())
		global.Logs(str)
		println(str)
		walk.MsgBox(mw, "工蜂小智-系统提示", "程序已打开，右下角的图标可打开主界面！", walk.MsgBoxIconInformation)
		os.Exit(1)
		return
	}
	utils.CloseChrome()
	global.InitGlobal()
	global.Log.Info("启动程序,当前版本：%s", global.VersionStr)
	_, err := update.CheckUpdate(global.Version, true)
	if err == nil {
		//启动更新程序
		flog, err := utils.PathExists("spread-update.exe")
		if err != nil || !flog {
			global.Log.Error("程序已破坏")
			//提示更新程序错误
			walk.MsgBox(mw, "工蜂小智-系统提示", "工蜂小智核心文件已损坏，请重新下载安装！", walk.MsgBoxIconError)
			utils.Open(global.API_DOMAIN)
			return
		} else {
			global.Log.Info("启动更新程序")
			_, err := os.StartProcess("spread-update.exe", nil, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
			if err != nil {
				global.Log.Error(err.Error())
			}
			time.Sleep(time.Second)
			Stop()
			return
		}
	}
	global.Log.Info("无版本更新")
	mw.init()
	mw.AddNotifyIcon()
	mw.Run()
}

// 判断进程是否启动
func ProcExsit() (err error) {
	pid, err := os.Open("pid.pid")
	defer pid.Close()
	if err == nil {
		filePid, err := ioutil.ReadAll(pid)
		if err == nil {
			pidStr := fmt.Sprintf("%s", filePid)
			pid, _ := strconv.Atoi(pidStr)
			_, err := os.FindProcess(pid)
			if err == nil {
				return errors.New("工蜂小智已启动")
			}
		}
	}
	return nil
}

func clearPid() {
	os.Truncate("pid.pid", 0)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Start() error {
	defer func() {
		if err := recover(); err != nil {
			global.Log.Error("《程序出现严重错误，终止运行！》，ERROR：%v", err)
		}
	}()
	syncMutex = new(sync.Mutex)
	go router.Router()
	runtime.GOMAXPROCS(runtime.NumCPU())
	load := global.Host
	acc := handler.CheckClientLogin()
	if acc == nil {
		load = global.API_SSO_DOMAIN + "?backUrl=" + global.Host + "/ReceiveToken"
	} else {
		global.Acc = acc
	}
	global.Driver = agouti.ChromeDriver(agouti.ChromeOptions("args", []string{
		"--user-data-dir=ChromeUserData",
		"--gpu-process",
		"--start-maximized",
		"--disable-infobars",
		"--app=" + global.LoadPage,
		"--webkit-text-size-adjust",
	}))
	global.Driver.Start()
	var err error
	global.Page.Page, err = global.Driver.NewPage()
	if err != nil {
		global.Log.Error("启动ChromeDriver失败，请重启. Error：%s", err.Error())
		return err
	}
	go func() {
		global.Page.Navigate(load)
	}()
	return nil
}

func Stop() {
	clearPid()
	global.Log.Info("退出桌面应用")
	if global.Page.Page != nil {
		global.Page.Page.CloseWindow()
	}
	if global.Driver != nil {
		global.Driver.Stop()
	}
	global.Log.Info("退出服务")
	router.Stop()
	utils.CloseSpread()
}
