package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/parser"
	"github.com/labstack/echo"
)

func Push(c echo.Context) error {
	title := c.FormValue("title")
	content := c.FormValue("content")

	go func() {
		println("开始执行简书分发")
		rule := utils.Read(utils.JsonPath("parser", "jianshu.json"))
		flog, result, err2 := PushComm(title, content, rule)
		println("简书分发", flog, result, err2)

		println("开始执行知乎分发")
		rule = utils.Read(utils.JsonPath("parser", "zhihu.json"))
		flog, result, err2 = PushComm(title, content, rule)
		println("知乎分发", flog, result, err2)

		println("开始执行新浪分发")
		rule = utils.Read(utils.JsonPath("parser", "sina.json"))
		flog, result, err2 = PushComm(title, content, rule)
		println("微博分发", flog, result, err2)
	}()

	return utils.Success(c, "正在发布中", "")
}

func PushComm(title string, content string, rule string) (bool, string, error) {
	//println("Title：", title, "，Content:", content, "，规则：", rule)
	paramMap := map[string]string{
		"loginName": "18223277005",
		"loginPwd":  "13696433488wb",
		"title":     title,
		"content":   content}
	return parser.RunPush(rule, paramMap)
}
