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

		//		time.Sleep(500 * time.Millisecond)
		//		println("开始执行简书分发")
		//		rule := utils.Read(utils.JsonPath("parser", "jianshu.json"))
		//		flog, result, err2 := handler.PushComm("十年前，大家都在哼《认真的雪》。十年后，大家都在唱薛之谦", "你创业卖衣服开火锅店年入百万，你写段子爆红大江南北，你上综艺节目频频曝光是为了什么呢？—做音乐。 前几日，亲戚们小聚。 酒足饭饱之后，众人提议去KTV亮亮嗓子，可把我们这些...", rule)
		//		println("简书分发", flog, result, err2)

		//global.Navigate("https://www.zhihu.com")

		global.Page.Navigate(global.Host)

		//		paramMap := map[string]string{
		//			"loginName": "18223277005",
		//			"loginPwd":  "13696433488wb",
		//			"title":     "路上被人认出是一名长期健身者是什么体验？",
		//			"content": `我一直没想到
		//		健身能改变我的命运
		//		来公司上班面试的时候.......
		//		应聘的是文案策划和编辑

		//		结果来这了之后发现这里妹纸特！别！多！`}
		//		rule := utils.Read(utils.JsonPath("parser", "sina.json"))
		//		flog, result, err2 := parser.RunPush(rule, paramMap)
		//		println(flog, result, err2)

	}()

}
