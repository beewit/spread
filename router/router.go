package router

import (
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/handler"

	"fmt"

	"time"

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

	global.Driver = agouti.ChromeDriver(agouti.ChromeOptions("args", []string{"--start-maximized", "--disable-infobars",
		"--app=http://www.jq22.com/demo/svgloader-150105194218/", "--webkit-text-size-adjust"}))
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

	//global.Page.Navigate("https://www.jianshu.com/sign_in#/")
	//handler.CheckLogin("jianshu", "www.jianshu.com", "remember_user_token")
	flog, result := handler.SetLoginCookie()
	println(flog, result)
	arguments := map[string]interface{}{"hiveHtml": global.HiveHtml, "host": host}

	js := ";$(function () {$('body').append(hiveHtml)});" + global.HiveJs

	global.Page.RunScript(js, arguments, nil)

	time.Sleep(1 * time.Second)

	global.Page.Navigate("http://www.jianshu.com/writer#/")

	global.Page.FindByID("note_title").Fill("每年的七夕，相爱的人儿，总在银河下面，相依相偎，诉说情深。")
	global.Page.FindByClass("kalamu-area").Fill(`<div>
    <section class="lanrenmb.com" id="xmy1476449703967" onmouseup="ajh(this.id);">
        <section class=" cur" id="xmy1476449565972" onmouseup="ajh(this.id);" style="">
            <section class="">
                <section style="text-align:center;margin:1em auto;">
                    <section style="width:1.5em;height:1.5em;border-radius:50%;background-color: rgb(2240,88,145);display:inline-block;vertical-align:top;" data-width="1.5em"></section>
                    <section style="width:0.8em;height:0.8em;border-radius:50%;background-color: rgb(135,81,175);display:inline-block;vertical-align:top;margin-right:-1.4em;margin-top:-1em;" data-width="0.8em"></section>
                    <section id="qqq0" class="lanrenmb.com" style="width:3em;height:3em;border-radius:50%;background-color: rgba(62,194,173,0.85);font-size:26px;text-align:center;line-height:3em;display:inline-block;margin-left:-0em;margin-top:0.4em;transform:rotate(0);-webkit-transform:rotate(0)" data-width="3em">
                        <span class="shuzi"><span style="color: rgb(255, 255, 255); font-size: 36px;"><!-- -->1<!-- --></span></span>
                    </section>
                    <section style="width:1.2em;height:1.2em;border-radius:50%;background-color: rgb(245,212,71);display:inline-block;vertical-align:top;margin-left:-2.2em;" data-width="1.2em"></section>
                    <section style="text-align:center;margin-top:-2em;margin-right:3.5em;">
                        <section style="width:2.1em;height:2.1em;border-radius:50%;background-color: rgb(54,153,184);display:inline-block;" data-width="2.1em"></section>
                    </section>
                    <section style="margin-top:-0.4em;text-align:center;margin-left:10%;">
                        <section style="width:0.8em;height:0.8em;border-radius:50%;background-color: rgb(224,100,88);display:inline-block;" data-width="0.8em"></section>
                        <section style="width:1em;height:1em;border-radius:50%;background-color: rgb(134,186,101);display:inline-block;margin-left:-0.2em;" data-width="1em"></section>
                    </section>
                </section>
            </section>
        </section>
    </section><br/>
    <div>
        <p>
            <br/>
        </p>
        <section class="lanrenmb.com" style="position: static; box-sizing: border-box; border: 0px none; padding: 0px;" data-id="85695">
            <section style="width:100%;box-sizing: border-box;" data-width="100%">
                <section style="width:25%;float:left;" data-width="25%">
                    <section class="wxqq-borderBottomColor" style="opacity: 0.8; margin-top: 5px; border-bottom-width: 10px; border-bottom-style: solid; border-bottom-color: rgb(209, 100, 27); border-top-color: rgb(235, 103, 148); box-sizing: border-box; color: inherit; float: left; border-left-width: 6px !important; border-left-style: solid !important; border-left-color: transparent !important; border-right-width: 6px !important; border-right-style: solid !important; border-right-color: transparent !important;"></section>
                    <section class="wxqq-borderRightColor" style="opacity: 0.4; border-right-width: 10px; border-left-width: 0px; border-right-style: solid; border-right-color: rgb(209, 100, 27); border-left-color: rgb(235, 103, 148); display: inline-block; float: left; color: inherit; margin-top: 5px; margin-left: 10px; margin-right: 5px; transform: rotate(10deg); border-bottom-width: 6px !important; border-top-width: 6px !important; border-top-style: solid !important; border-bottom-style: solid !important; border-top-color: transparent !important; border-bottom-color: transparent !important;"></section>
                    <section class="wxqq-borderLeftColor" style="border-left-width: 20px; border-right-width: 0px; border-left-style: solid; border-left-color: rgb(209, 100, 27); border-right-color: rgb(235, 103, 148); display: inline-block; float: left; color: inherit; transform: rotate(10deg); border-bottom-width: 10px !important; border-top-width: 15px !important; border-top-style: solid !important; border-bottom-style: solid !important; border-top-color: transparent !important; border-bottom-color: transparent !important;"></section>
                </section>
                <section style="width:50%;text-align:center;float: left;padding: 0px 5px;" data-width="50%">
                    <section style="display:inline-block;">
                        <span style="font-size:18px"><strong><span class="135brush" data-brushtype="text" style="font-size:18px">一、懒人图文排版</span></strong></span>
                    </section>
                </section>
                <section style="float:right;width:25%;text-align: right;" data-width="25%">
                    <section style="display:inline-block;">
                        <section class="wxqq-borderBottomColor" style="opacity: 0.8; margin-top: 5px; border-bottom-width: 10px; border-bottom-style: solid; border-bottom-color: rgb(209, 100, 27); border-top-color: rgb(235, 103, 148); box-sizing: border-box; color: inherit; float: left; border-left-width: 6px !important; border-left-style: solid !important; border-left-color: transparent !important; border-right-width: 6px !important; border-right-style: solid !important; border-right-color: transparent !important;"></section>
                        <section class="wxqq-borderRightColor" style="opacity: 0.4; border-right-width: 10px; border-left-width: 0px; border-right-style: solid; border-right-color: rgb(209, 100, 27); border-left-color: rgb(235, 103, 148); display: inline-block; float: left; color: inherit; margin-top: 5px; margin-left: 10px; margin-right: 5px; transform: rotate(10deg); border-bottom-width: 6px !important; border-top-width: 6px !important; border-top-style: solid !important; border-bottom-style: solid !important; border-top-color: transparent !important; border-bottom-color: transparent !important;"></section>
                        <section class="wxqq-borderLeftColor" style="border-left-width: 20px; border-right-width: 0px; border-left-style: solid; border-left-color: rgb(209, 100, 27); border-right-color: rgb(235, 103, 148); display: inline-block; float: left; color: inherit; transform: rotate(10deg); border-bottom-width: 10px !important; border-top-width: 15px !important; border-top-style: solid !important; border-bottom-style: solid !important; border-top-color: transparent !important; border-bottom-color: transparent !important;"></section>
                    </section>
                </section>
            </section>
            <section style="clear:both;"></section>
        </section>
        <p>
            <br/>
        </p>
    </div>
</div>
<p>
    <br/>
</p>
<div>
    <section id="see7264" style="margin: 0px;padding:8px;width:380px; min-width:50px;border: 0px solid #DDD; background-color:#FFF;border-radius: 5px;zoom:0.93">
        <section label="Powered by " class="_editor" data-tools="懒人微信编辑器" powered-by="">
            <section style="width: 100%; box-sizing: border-box; background-image: url(http://buluo.lanrenmb.com/bj/images/jieri/17082112.jpg);background-size:100%" powered-by="">
                <section style="background:rgba(196,4,0,0.6);line-height: 24px;box-sizing:border-box;" powered-by="">
                    <section style="color:#fff;letter-spacing: 2px;line-height: 26px;padding:10px 10px;text-align:justify;box-sizing:border-box;" powered-by="">
                        <p style="margin:0; padding:0; font-size:14px;" powered-by="">
                            &nbsp;&nbsp;&nbsp;&nbsp;我们常说，时光易老，真情难觅。虽然世间凉薄，但我更愿意相信传说中的故事，那份坚贞不移的深情，那一寸寸的相思，用无数年的时光交错，流转成为一段浪漫爱情的象征。
                        </p>
                        <p style="margin:0; padding:0; font-size:14px;text-align:center;border-radius:50%;" powered-by="">
                            <img src="http://buluo.lanrenmb.com/bj/images/jieri/17082113.jpg" style="width:100%;border-radius:50%;"/>
                        </p>
                        <p style="margin:0; padding:0; font-size:14px;" powered-by="">
                            &nbsp;&nbsp;&nbsp;&nbsp;如今，每年的七夕，相爱的人儿，总在银河下面，相依相偎，诉说情深。而你，在爱的季节里，是不是也和我一样，隔着一片海，放飞心中的无限执念？看一片片花瓣飘摇而落，我在七夕的天空下，任相思的长河漫过心海，泛滥成伤。
                        </p>
                        <p style="margin:0; padding:0; font-size:14px;" powered-by="">
                            <br/>
                        </p>
                        <p style="margin:0; padding:0; font-size:14px;" powered-by="">
                            “背景可扩展，长文效果更佳哦！”
                        </p>
                    </section>
                </section>
            </section>
        </section>
    </section>
</div>
<p>
    <br/>
</p>
<p>
    <br/>
</p>`)

}
