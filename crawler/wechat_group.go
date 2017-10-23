package crawler

import (
	"encoding/json"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/parser"
	"container/list"
	"github.com/beewit/beekit/utils/convert"
)

type Spider struct {
	StartUrl  []string              `json:"start_url"`
	Domain    string                `json:"domain"`
	Fill      []parser.PushFillJson `json:"fill"`
	Page      *utils.AgoutiPage
	Queue     *list.List
	DoneQueue *list.List
	FailQueue *list.List
}

func Start(page *utils.AgoutiPage, rule string) (spider *Spider, err error) {
	spider, err = GetPushJson(rule)
	if err != nil {
		return
	}
	if len(spider.StartUrl) > 0 {
		for i := 0; i < len(spider.StartUrl); i++ {
			spider.Queue.PushBack(spider.StartUrl[i])
		}
	}
	spider.Page = page
	return
}

func (spider *Spider) Run() {
	ele := spider.Queue.Front()
	if ele != nil {
		url := convert.ToString(ele.Value)
		spider.Page.Navigate(url)


		if spider.Queue.Len() > 0 {
			spider.Run()
		}
	}
	parser.HandleSelection(&spider.Fill[0], nil)

}

func GetPushJson(rule string) (*Spider, error) {
	var spider Spider
	err := json.Unmarshal([]byte(rule), &spider)
	if err != nil {
		global.Log.Error("规则解析错误：%s", err.Error())
		return &Spider{}, err
	} else {
		return &spider, nil
	}
}
