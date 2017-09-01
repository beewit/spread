package main

import (
	"fmt"
	"encoding/json"
)

func push() {

}

type PushJson struct {
	Title    string `json:"title"`
	Domain   string `json:"domain"`
	LoginUrl string `json:"loginUrl"`
	Identity string `json:"identity"`
	Writer   string `json:"writer"`
	Fill     []PushFillJson
}

type PushFillJson struct {
	Selector     string `json:"selector"`
	SelectorName string `json:"selectorName"`
	Handle       string `json:"handle"`
	Sleep        string `json:"sleep"`
	Js           string `json:"js"`
	JsParam      string `json:"jsParam"`
	Param        string `json:"param"`
	Result       string `json:"result"`
}

func main() {
	pf := PushFillJson{
		"ID",
		"share-modal",
		"Click",
		"1000",
		"alert(1)",
		`{"hiveHtml": "hiveHtml", "host": "host"}`,
		"title",
		"已发布",
	}

	var pfs []PushFillJson

	pfs = append(pfs, pf)
	pf = PushFillJson{
		"ID2",
		"share-modal2",
		"Click",
		"1000",
		"alert(1)",
		`{"hiveHtml": "hiveHtml", "host": "host"}`,
		"title",
		"已发布",
	}
	pfs = append(pfs, pf)
	st := PushJson{
		"简书",
		"http://www.jianshu.com",
		"https://www.jianshu.com/sign_in",
		"remember_user_token",
		"http://www.jianshu.com/writer#/",
		pfs,
	}

	b, err := json.Marshal(st)

	if err != nil {
		fmt.Println("encoding faild")
	} else {
		j := string(b)
		var pj PushJson
		fmt.Println(j)
		err := json.Unmarshal(b, &pj)
		if err != nil {
			println(err.Error())
		} else {
			println("结果：" + pj.Title)
		}
	}

}




##
{
  "title": "简书",
  "domain": "http://www.jianshu.com",
  "loginUrl": "https://www.jianshu.com/sign_in",
  "identity": "remember_user_token",
  "writer": "http://www.jianshu.com/writer#/",
  "fill": [
    {
      "handle": "Click",
      "selector": "ID",
      "selectorName": "share-modal",
      "selectorVal": "share-modal",
      "sleep": 1000,
      "js": "$('#new-note').click()",
      "jsParam": {
        "hiveHtml": "hiveHtml",
        "host": "host"
      },
      "param": "title",
      "result": "已发布"
    }
  ]
}

###
{
  "title": "标题",
  "domain": "主域名",
  "loginUrl": "登陆Url",
  "identity": "登陆身份识别",
  "writerUrl": "新建文章url",
  "fill": [
    {
      "handle": "操作项 Click/DoubleClick/Text/Fill",
      "selector": "Selector/ID/Class/Name",
      "selectorName": "#id",
      "selectorVal": "填充值",
      "sleep": 秒,
      "js": "$('#new-note').click()",
      "jsParam": {
        "hiveHtml": "hiveHtml",
        "host": "host"
      },
      "param": "title",
      "result": "已发布"
    }
  ]
}

## 知乎

    {
      "handle": "Js",
      "sleep": 1,
      "js": "document.getElementsByTagName(\"textarea\")[0].value=title",
      "jsParam": {
        "title": "title/v"
      }
    },
    {
      "handle": "Js",
      "sleep": 1,
      "js": "document.getElementsByClassName(\"public-DraftEditor-content\")[0].innerHTML=content;",
      "jsParam": {
        "content": "content/v"
      }
    },
    {
      "handle": "Js",
      "sleep": 1,
      "js": "document.getElementsByClassName(\"PublishPanel-triggerButton\")[0].disabled=false;document.getElementsByClassName(\"PublishPanel-wrapper\")[0].click()"
    },
    {
      "handle": "Text",
      "sleep": 0,
      "result": "已发布"
    }


##新浪

    {
      "handle": "Js",
      "sleep": 1,
      "js": "document.getElementsByClassName(\"hive-load\")[0].className+=' hide';",
      "sleep": 1
    },
    {
      "handle": "Click",
      "Selector": "selector",
      "SelectorName": ".addopt a",
      "sleep": 1
    },
    {
      "handle": "Fill",
      "Selector": "selector",
      "SelectorName": "input[node-type=\"title\"]",
      "Param": "title",
      "sleep": 1
    },
    {
      "handle": "Fill",
      "Selector": "selector",
      "SelectorName": "#editor",
      "Param": "content",
      "sleep": 1
    },






### 图片上传
	global.Page.Navigate("http://www.baidu.com")
	global.Page.Find(".soutu-btn").Click()
	time.Sleep(10 * time.Second)
	println(global.Page.Find(".UploadFile"))
	global.Page.Find(".upload-pic").UploadFile("/home/zxblovelc/goProjects/src/github.com/beewit/spread/app/static/img/logo-80*80.png")
	global.Page.RunScript("alert(1)", nil, nil)


### sina.json

    {
      "handle": "Fill",
      "Selector": "selector",
      "SelectorName": "#editor",
      "Param": "content",
      "sleep": 5
    },