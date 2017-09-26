package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"container/list"
	"fmt"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/beewit/beekit/utils/convert"
	"github.com/go-vgo/robotgo"
	"golang.org/x/net/html/charset"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
)

var (
	l          = list.New()
	todo       = list.New()
	gourpCount = 10
	m          = map[string]interface{}{}

	listUrls = list.New()
)

func TestCrawlerBBS(t *testing.T) {
	pageTodo("http://www.bbsbaba.com/")
	addTodo()
	groupSum()
}

func groupSum() {
	i := 0
	for {
		println("统计数据-------------")
		if l.Len() > 0 {
			c := gourpCount
			if l.Len() < gourpCount {
				c = l.Len()
			}
			for j := 0; j < c; j++ {
				t := l.Front()
				l.Remove(t)
				go groupDay(convert.ToString(t.Value))
			}
		}
		if i > 30 {
			str, _ := convert.ToMapStr(m)
			println(str)
			println("结束统计xxxxxxxxxxxxxxxxxxxxx")
			break
		}
		i++
		time.Sleep(time.Second * 1)
	}
}

func TestRegexp2(t *testing.T) {
	dat, err := ioutil.ReadFile("sitedetail.txt")
	if err != nil {
		println(dat)
	}
	str := string(dat)
	strs := strings.Split(str, "\n")
	println("所有：", len(strs))

	olds := getOlds()
	println("已抓", len(olds))
	var s string
	for i := 0; i < len(strs); i++ {
		s = strings.TrimSpace(strs[i])
		if s != "" && strings.Contains(s, "s.asp") != StrsContains(olds, s) {
			pageDetail(s)
			time.Sleep(time.Second * 2)
		}
	}
	//groupDay("http://www.tanzhou.com.cn/bbs/forum.php")
}

func StrsContains(strs []string, s string) bool {
	for i := 0; i < len(strs); i++ {
		if strs[i] == s {
			return true
		}
	}
	return false
}

func getOlds() []string {
	dat, err := ioutil.ReadFile("siteOld.txt")
	if err != nil {
		println(dat)
	}
	str := string(dat)
	strs := strings.Split(str, "\n")
	return strs
}

func groupDay(url string) {
	if url == "" {
		return
	}
	resp, e := http.Get(url)
	if e != nil {
		println(e.Error)
		return
	}
	var ir io.Reader
	if strings.Contains(resp.Header.Get("Content-Type"), "utf-8") {
		ir = resp.Body
	} else {
		var errIr error
		ir, errIr = charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
		if errIr != nil {
			println(errIr.Error)
			return
		}
	}

	doc, err := goquery.NewDocumentFromReader(ir)
	if err != nil {
		println(err.Error())
		return
	}
	html, err2 := doc.Html()
	if err2 != nil {
		println(err2.Error())
		return
	}
	if html != "" {
		num, _ := getDayNum(html)
		if num != "" {
			m[url] = num
			WriteAt("url:" + num)
			println("今日数据[", num, "]："+url)
		} else {
			println("解析今日数据失败：" + url)
		}
	} else {
		println("解析Html失败：" + url)
	}
}

func groupDayNew(url string) (string, string, string) {
	if url == "" {
		return "", "", ""
	}
	resp, e := http.Get(url)
	if e != nil {
		println(e.Error)
		return "", "", ""
	}
	var ir io.Reader
	if strings.Contains(resp.Header.Get("Content-Type"), "utf-8") {
		ir = resp.Body
	} else {
		var errIr error
		ir, errIr = charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
		if errIr != nil {
			println(errIr.Error)
			return "", "", ""
		}
	}

	doc, err := goquery.NewDocumentFromReader(ir)
	if err != nil {
		println(err.Error())
		return "", "", ""
	}
	html, err2 := doc.Html()
	if err2 != nil {
		println(err2.Error())
		return "", "", ""
	}
	if html != "" {
		t, y := getDayNum(html)
		return t, y, doc.Find("title").Text()

	}
	return "", "", ""
}

func TestRegexp3(t *testing.T) {
	str := `登录
	新用户注册
	用其他账号登录:
	`
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, "	", "", -1)
	println(str)
}

func TestRegexp(t *testing.T) {
	src := strings.TrimSpace(`<p class="chart z">今日: <em>3108</em><span class="pipe">|</span>昨日: <em>18903</em><span class="pipe">|</span>帖子: <em>373955820</em><span class="pipe">|</span>会员: <em rel="last">30316800</em><span class="pipe">|</span>欢迎新会员: <em><a href="space-uid-30316800.html" target="_blank" class="xi2">牛奶泡芒果喝</a></em></p>`)

	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("<script[^>]*?>.*?</script>")
	src = re.ReplaceAllString(src, "")

	//去除所有的标签
	re, _ = regexp.Compile("<[^>]*>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("<[^>]*>")
	src = strings.Replace(src, " ", "", -1)
	src = strings.Replace(src, "：", ":", -1)

	fmt.Println(src)

	fmt.Println("--------------------------------------------")

	re = regexp.MustCompile("今日:(\\d+)")
	fmt.Println("--------------------------------------------")
	fmt.Println(re.FindString(src))
	fmt.Println("--------------------------------------------")
	data := re.Find([]byte(src))

	fmt.Println(strings.Replace(string(data), "今日:", "", -1))
}

func getDayNum(html string) (string, string) {
	src := strings.TrimSpace(html)

	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("<script[^>]*?>.*?</script>")
	src = re.ReplaceAllString(src, "")

	//去除所有的标签
	re, _ = regexp.Compile("<[^>]*>")
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile("<[^>]*>")
	src = strings.Replace(src, " ", "", -1)
	src = strings.Replace(src, "：", ":", -1)
	src = strings.Replace(src, "\n", "", -1)
	src = strings.Replace(src, "\r", "", -1)
	src = strings.Replace(src, "	", "", -1)

	fmt.Println("--------------------------------------------")

	re = regexp.MustCompile("今日:(\\d+)")

	re2 := regexp.MustCompile("昨日:(\\d+)")
	fmt.Println("--------------------------------------------")
	fmt.Println(re.FindString(src))
	fmt.Println(re2.FindString(src))
	fmt.Println("--------------------------------------------")
	data := re.Find([]byte(src))
	relust := string(data)

	data2 := re2.Find([]byte(src))
	relust2 := string(data2)

	var today, yesterday string
	if relust != "" {
		today = strings.Replace(string(data), "今日:", "", -1)
	}
	if relust2 != "" {
		yesterday = strings.Replace(string(data2), "昨日:", "", -1)
	}
	return today, yesterday
}

func pageTodo(url string) {
	println("启动爬虫")
	doc, err := goquery.NewDocument(url)
	if err != nil {
		println(err.Error())
		return
	}

	// Find the review items
	doc.Find("#tab2 a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			println("加入ToDo：" + href)
			todo.PushBack(href)
		}
	})

	doc.Find("#ncontentbody .ncontentbdcen a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			println("加入List：" + href)
			l.PushBack(href)
		}
	})
	println("todo：", todo.Len(), "list：", l.Len())
}

func pageDetail(url string) {
	println("siteUrl抓取：", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		println(err.Error())
		return
	}

	doc.Find(".tjdetailrg a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			WriteAtPath(url, "./siteOld.txt")
			WriteAtPath(href, "./siteDoMain.txt")
			l.PushBack(href)
		}
	})
}

func TestPageDetail(t *testing.T) {
	doc, err := goquery.NewDocument("http://www.bbsbaba.com/s.asp?id=13239")
	if err != nil {
		println(err.Error())
		return
	}

	doc.Find(".tjdetailrg a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			l.PushBack(href)
		}
	})
}

func addTodo() {
	i := 0
	for {
		println("待抓队列追加-------------")
		if l.Len() > 0 {
			t := l.Front()
			l.Remove(t)
			go pageDetail(convert.ToString(t.Value))
		}
		if i > 30 {
			println("结束待抓队列追加xxxxxxxxxxxxxxxxxxxxx")
			break
		}
		i++
		time.Sleep(time.Second * 3)
	}
}

func TestRobotgo(t *testing.T) {
	fpid, err := robotgo.FindIds("chrome")
	if err == nil {
		fmt.Println("pids...", fpid)
	}

	isExist, err := robotgo.PidExists(100)
	if err == nil {
		fmt.Println("pid exists is", isExist)
	}

	abool := robotgo.ShowAlert("test", "robotgo")
	if abool == 0 {
		fmt.Println("ok@@@", "ok")
	}

	title := robotgo.GetTitle()
	fmt.Println("title@@@", title)
}

func TestCrawlerBBSNew(t *testing.T) {
	hrefs := `http://www.bbsbaba.com/guangdong.html,http://www.bbsbaba.com/guangxi.html,http://www.bbsbaba.com/hunan.html,http://www.bbsbaba.com/hubei.html,http://www.bbsbaba.com/fujian.html,http://www.bbsbaba.com/jiangsu.html,http://www.bbsbaba.com/zhejiang.html,http://www.bbsbaba.com/anhui.html,http://www.bbsbaba.com/jiangxi.html,http://www.bbsbaba.com/henan.html,http://www.bbsbaba.com/hebei.html,http://www.bbsbaba.com/liaoning.html,http://www.bbsbaba.com/shandong.html,http://www.bbsbaba.com/shanxi.html,http://www.bbsbaba.com/shaanxi.html,http://www.bbsbaba.com/jilin.html,http://www.bbsbaba.com/sichuan.html,http://www.bbsbaba.com/guizhou.html,http://www.bbsbaba.com/yunnan.html,http://www.bbsbaba.com/hainan.html,http://www.bbsbaba.com/xinjiang.html,http://www.bbsbaba.com/xizang.html,http://www.bbsbaba.com/gansu.html,http://www.bbsbaba.com/qinghai.html,http://www.bbsbaba.com/xinwen.html,http://www.bbsbaba.com/yule.html,http://www.bbsbaba.com/game.html,http://www.bbsbaba.com/ruanjian.html,http://www.bbsbaba.com/junshi.html,http://www.bbsbaba.com/wenxue.html,http://www.bbsbaba.com/lvyou.html,http://www.bbsbaba.com/tiyu.html,http://www.bbsbaba.com/shouji.html,http://www.bbsbaba.com/diannao.html,http://www.bbsbaba.com/nanren.html,http://www.bbsbaba.com/meishi.html,http://www.bbsbaba.com/nvxing.html,http://www.bbsbaba.com/qinggan.html,http://www.bbsbaba.com/caijing.html,http://www.bbsbaba.com/zhiming.html,http://www.bbsbaba.com/boke.html,http://www.bbsbaba.com/tuiguang.html,http://www.bbsbaba.com/ningxia.html,http://www.bbsbaba.com/neimenku.html,http://www.bbsbaba.com/heilongjiang.html,http://www.bbsbaba.com/beijing.html,http://www.bbsbaba.com/tianjin.html,http://www.bbsbaba.com/shanghai.html,http://www.bbsbaba.com/chongqing.html,http://www.bbsbaba.com/hk.html,http://www.bbsbaba.com/jiaoyu.html,http://www.bbsbaba.com/yishu.html,http://www.bbsbaba.com/yinyue.html,http://www.bbsbaba.com/zhanzhang.html,http://www.bbsbaba.com/aihao.html,http://www.bbsbaba.com/liuxue.html,http://www.bbsbaba.com/jiankang.html,http://www.bbsbaba.com/yingshi.html,http://www.bbsbaba.com/gouwu.html,http://www.bbsbaba.com/jiaoyou.html,http://www.bbsbaba.com/qinzi.html,http://www.bbsbaba.com/gongye.html,http://www.bbsbaba.com/zongjiao.html,http://www.bbsbaba.com/gaoxiao.html,http://www.bbsbaba.com/fangchan.html,http://www.bbsbaba.com/qiche.html,http://www.bbsbaba.com/zhichang.html,http://www.bbsbaba.com/shuma.html,http://www.bbsbaba.com/sheji.html,http://www.bbsbaba.com/chongwu.html,http://www.bbsbaba.com/anquan.html,http://www.bbsbaba.com/chuangye.html,http://www.bbsbaba.com/sheying.html,http://www.bbsbaba.com/jiaju.html,http://www.bbsbaba.com/caipiao.html,http://www.bbsbaba.com/huwai.html,http://www.bbsbaba.com/nongye.html,http://www.bbsbaba.com/dongman.html,http://www.bbsbaba.com/jianshen.html,http://www.bbsbaba.com/mingxing.html,http://www.bbsbaba.com/falv.html,http://www.bbsbaba.com/fengshui.html,http://www.bbsbaba.com/xiuxian.html,http://www.bbsbaba.com/jiaotong.html,http://www.bbsbaba.com/gongyi.html,http://www.bbsbaba.com/fuzhuang.html,http://www.bbsbaba.com/fuwu.html,http://www.bbsbaba.com/aomen.html,http://www.bbsbaba.com/taiwan.html,http://www.bbsbaba.com/yazhou.html,http://www.bbsbaba.com/oumei.html,http://www.bbsbaba.com/wudao.html,http://www.bbsbaba.com/huahui.html,http://www.bbsbaba.com/zhubao.html,http://www.bbsbaba.com/diaoyu.html,http://www.bbsbaba.com/keji.html,http://www.bbsbaba.com/gupiao.html,http://www.bbsbaba.com/kuaiji.html,http://www.bbsbaba.com/biancheng.html,http://www.bbsbaba.com/guanggao.html`
	h := strings.Split(hrefs, ",")
	for i := 0; i < len(h); i++ {
		CrawlerBBSHrefs(h[i])
		time.Sleep(3 * time.Second)
	}

	println("todo：", todo.Len())
	getSiteHref()
	println("list：", l.Len())
	groupSum()

	println("map：", len(m))
	dayNums, _ := convert.ToMapStr(m)
	println(dayNums)
	WriteAt("," + dayNums)
}

func TestWriteFile(t *testing.T) {
	WriteAt("张三丰12")
}

func WriteAt(content string) {
	WriteAtPath(content, "./site.txt")
}

func WriteAtPath(content, path string) {
	//以读写方式打开文件，如果不存在，则创建
	file2, error := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0766)
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println(file2)
	defer file2.Close()
	n, _ := file2.Seek(0, os.SEEK_END)
	file2.WriteAt([]byte("\n"+content), n)
}

func CrawlerBBSHrefs(url string) {
	println("----------------------------------")
	println("启动爬虫", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		println(err.Error())
		return
	}

	// Find the review items
	doc.Find(".colList_clss a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			println("加入ToDo：" + href)
			todo.PushBack(href)
		}
	})
}

func getSiteHref() {
	println("待抓队列追加-------------")
	if todo.Len() > 0 {
		for i := 0; i < todo.Len(); i++ {
			t := todo.Front()
			todo.Remove(t)
			if t.Value != nil {
				pageDetail(convert.ToString(t.Value))
			}
			time.Sleep(time.Second * 3)
		}
	}
}

func TestBBS(t *testing.T) {
	dat, err := ioutil.ReadFile("siteDoMain.txt")
	if err != nil {
		println(dat)
	}
	str := string(dat)
	strs := strings.Split(str, "\n")
	for i := 0; i < len(strs); i++ {
		listUrls.PushBack(strs[i])
	}

	for {
		if listUrls.Len() > 0 {
			c := gourpCount
			if listUrls.Len() < gourpCount {
				c = listUrls.Len()
			}
			for j := 0; j < c; j++ {
				t := listUrls.Front()
				listUrls.Remove(t)
				if t.Value != nil {
					//go func() {
						today, yesterday, title := groupDayNew(convert.ToString(t.Value))
						if today == "" && yesterday == "" && title == "" {

						}else {
							iw, _ := utils.NewIdWorker(1)
							id, _ := iw.NextId()
							m := make(map[string]interface{})
							m["id"] = id
							m["title"] = title
							m["url"] = convert.ToString(t.Value)
							if today != "" {
								m["today"] = today
							}
							if yesterday != "" {
								m["yesterday"] = yesterday
							}
							if today == "" && yesterday == "" {
								m["state"] = 0
							} else {
								m["state"] = 1
							}
							x, err := global.SLDB.InsertMap("bbs", m)
							if err != nil {
								global.Log.Error(err.Error())
							} else {
								println("保存数据成功！", convert.ToString(x))
							}
						}
					//}()
				}
			}
		} else {
			break
		}
		time.Sleep(time.Second * 3)
	}
	println("====================================={结束完成}=====================================")
}
