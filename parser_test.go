package main

import (
	"testing"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/handler"
	"github.com/sclevine/agouti"
	"fmt"
	"github.com/beewit/spread/global"
	"github.com/beewit/spread/api"
	"github.com/beewit/beekit/utils/uhttp"
	"time"
)

func TestParser(t *testing.T) {
	Parser("700农民不种田专画老虎 全村年收入过亿", `<div contenteditable="true" style="width:100%; height:100%;" class="w-e-text"><p>靠画老虎富甲一方的神奇村落</p><p>王公庄是一座位于豫鲁交界处，占地1400亩的小村落。和周围大多数靠种地打工度日的村庄不同，这里别有一番纸墨香。最为外界津津乐道的，是王公庄的“虎文化”生意。村里没有养老虎，却有700多村民都会画老虎，全村年售几万张画作，收入过亿。不少村民都坐拥多套房产，换了几部私家车，年收入近百万。</p><p>从民权县城驱车28公里，穿过一片玉米地，当看到一块文化广场石碑和一座恢宏的展览馆时，就到了久负盛名的“中国第一画虎村”王公庄。只见一排排二层民房整齐划一，一块块书画招牌琳琅满目，村子里鲜少有人闲逛，推开一间间画室却可以见到男女老少的村民都在静心作画，孩童也大都在一旁看书学习，不吵不闹。</p><div><img width="500px" src="https://ss1.baidu.com/6ONXsjip0QIZ8tyhnq/it/u=4167756370,3115892551&amp;fm=173&amp;s=D59B3FD756C14AEA6C9405730300A070&amp;w=500&amp;h=281&amp;img.JPG"></div><p>画室里多个村民齐画“百虎图”</p><p>王公庄有村民1600人左右，其中700多人都在从事着绘画行业，其中300多人堪称画师。村里有夫妻画家、父子画家、姐妹画家，三代同堂的绘画家庭不在少数，画的大都是老虎。上山虎、下山虎，虎头和丛林虎，百虎、千虎……风格不尽相同，却都可圈可点、栩栩如生。</p><div><img width="500px" src="https://ss2.baidu.com/6ONYsjip0QIZ8tyhnq/it/u=1415589314,824664614&amp;fm=173&amp;s=F5AB95570AA3A0D204A4A4A70300F043&amp;w=500&amp;h=281&amp;img.JPG"></div><p>给虎头画须</p><p>然而最为外界津津乐道的，却是王公庄的卖画生意。全国大中小型的书画市场上几乎都能觅得王公庄虎画的踪影。村里“虎王”的真迹动辄能卖出几十万、上百万的高价，可以跻身收藏界；普通画师的虎画也是被各地画商抢购一空；就连一般学生的临摹作业，也能卖出几百块到上千块不等的价格，经常被旅游景点用作装饰品。</p><p>于是乎，坊间形容画虎村是“一人画虎，鸡犬升天”。</p><p>画虎技艺竟是祖传，先辈“照猫画虎”？</p><p>王公庄卖出过最贵的一幅虎画，是“四大虎王”之首王建民的作品，售价高达九十多万元。王建民形容自己的画虎技艺，主要是在虎毛的创作上走出了自己的风格，“飘逸有厚度”，还带有一些动感。只见他一手执两笔交替描画虎毛，手法十分娴熟。</p><p>“我们并没有受过正规的培训，很多技法就是从书本上看到的，过去就是临摹画虎名家。”王建民说，整个王公庄的画虎技艺实质是从祖辈上传承下来的，像他从小就是跟随父亲和爷爷在村子里描画中堂和房梁，画的都是鱼、荷、虎之类有美好寓意的图像。</p><p>让他印象深刻的是，小时候每逢节日就跟着大人们在集市上挂卖老虎年画，“我问过我爷爷，你们见过老虎吗？他说没有，他们那代人就是照猫画虎。”</p><p>都说画龙画虎难画骨，然而王公庄里“照猫画虎”的技艺，却在传承了好几代之后成了炙手可热的民间艺术，惠泽至今。</p><div><img width="500px" src="https://ss2.baidu.com/6ONYsjip0QIZ8tyhnq/it/u=3644083727,2753691841&amp;fm=173&amp;s=B2115DCF065A39DA5E08813A0300D052&amp;w=500&amp;h=281&amp;img.JPG"></div><p>“虎王”王建民和他的儿子</p><p>到了王建民成年的时候，他和现今并称“四大虎王”的其他三个小伙伴，不再甘于只在村里的集市上卖画了，他们急切地想出去闯一闯。陆续去到过开封、洛阳等省内城市的书画市场，他们有遭受过轻蔑，也经常入不敷出。然而就是在一次次碰壁，一次次改进之后，王公庄的虎画终于打开了外地市场。</p><p>虎，谐音就是“福”。题材非常接地气，老百姓都很喜欢。渐渐地，他们便不再描画一些花鸟虫鱼了，而是主攻画老虎，为此还经常结伴到动物园里近距离观察老虎以提高画技。</p><div><img width="500px" src="https://ss0.baidu.com/6ONWsjip0QIZ8tyhnq/it/u=943469190,11861481&amp;fm=173&amp;s=2CB64F95C41959D25EA0F40C0300B0D3&amp;w=500&amp;h=281&amp;img.JPG"></div><p>他们常去距离王公庄最近的商丘动物园看老虎</p><p>在20多年前，王建民的一张六尺长的老虎画，就卖出了一百块的高价，而当时工人的工资也就只有30多块一个月。“画老虎能挣钱”的消息一下子在王公庄传播开来，亲戚、乡邻纷纷找到王建民要学习画老虎。他回忆说，从办学的第一天起，村里就呼一下冒出了五六十个学徒。</p><p>手画老虎，眼观市场</p><p>现年51岁的“虎王”王建民，以稼轩主人自居，在自己300平米的两层民宅里划出一个房间专门用来陈列虎画。这些巨幅老虎画作，伴随着他巡展过全国和海外十几个城市，每一幅他都舍不得出手。</p><div><img width="500px" src="https://ss0.baidu.com/6ONWsjip0QIZ8tyhnq/it/u=3918558033,1808569012&amp;fm=173&amp;s=E29B1DC7461049CA1C2674720300507B&amp;w=500&amp;h=281&amp;img.JPG"></div><p>王建民家里的虎画展示厅</p><p>当记者和他聊起眼下正热火朝天的画虎产业时，他颇有些得意地评价王公庄的村民，“不但是画家，还是市场揣摩家”。</p><p>早在“四大虎王”外出闯荡的时期，虎王们就开始在村子里培养售画经纪人了。“不是每个人都是艺术家，也不是每个人都有创作能力”，有绘画天分的村民安心作画，而活络一些的就带着他们的画作专门去跑市场，卖画的同时给村里带回市场的需求和讯息。</p><p>王建金是王公庄上第一个卖画的经纪人，早期他单枪匹马跑外地，睡过桥洞，风餐露宿。而今，他已经在王公庄上成立了一家经纪人公司，手里十几个员工分头在全国各地联络书画市场。</p><div><img width="500px" src="https://ss1.baidu.com/6ONXsjip0QIZ8tyhnq/it/u=1415003971,1192433827&amp;fm=173&amp;s=686001D16A124CDA162D584A0300A070&amp;w=500&amp;h=281&amp;img.JPG"></div><p>经纪人是村民画师与外界市场之间的桥梁</p><p>一次次接触外界之后，他引领着村里的画虎题材进行了几次变革，保证了村里的画作，始终走在市场的前沿。从单一的虎头、虎身，到带背景，有叙事的巨幅虎图；从单只上山虎、下山虎，到三口虎，寓意祝福的五虎、八虎、百虎；从国内流行的中西结合虎画，到畅销海外的老虎与古典美女……至今，他的画室里还珍藏着一幅摊开来足足有四百多米长的千虎图，当时就是为了打响“画虎第一村”的名气，他引领着村里几个重量级画师一同设计、组图。</p><p>“一人画虎，鸡犬升天”</p><p>“画虎村”的“虎”产业究竟有多虎虎生风？当今在王公庄年销量第一，收学徒第一的是“四小虎王”之一王建辉，他也是王建民最得意的学徒。</p><p>“小虎王”王建辉留着过耳长发，33岁的年纪已经画虎近20年。“当初家里穷，灯都没有，就点着煤火学画。冬天手都冻烂了，晚上要熬到三四点，夏天身上粘的都是纸。”虽然从小学习不好，但王建辉却凭借画虎改变了一生的命运，他是画虎村文化产业发展的一个缩影。</p><div><img width="500px" src="https://ss1.baidu.com/6ONXsjip0QIZ8tyhnq/it/u=1186621969,290163345&amp;fm=173&amp;s=F5AB95570AA3A0D204A4A4A70300F043&amp;w=500&amp;h=281&amp;img.JPG"></div><p>王建辉擅长画虎头</p><p>“这是要寄给山东一个老客户的”，王建辉卷起了十几幅画作，“一次几十万的订单现在是常态，他的四层酒楼里每个角落都用我的画装饰。”</p><p>“上山虎寓意步步高升，升官发财；下山虎镇宅辟邪保平安；五只老虎是五福临门；收藏百虎是纳百福……” 王建辉的妻子杨美菊是和他同期学习的画师，现在帮忙打理着画室和销售的生意，锻炼得精明能干。杨美菊回忆说，从当初结婚买不起一辆摩托车，到现在家里前后换了五部私家车，购置了房产，积蓄了存款，她很满意现在的生活。“我们现在一个月光发工资几十万，一年发工资也要几百万吧。”</p><p>王建辉其中一个画室旁边，是赵庆伟和王喜梅夫妻俩的装裱店。半夜十二点时分，夫妻俩仍然还在劳作，需要装裱的绘画订单实在太多。然而通常只要辛苦两天一夜，就能收入两千块钱。</p><div><img width="500px" src="https://ss2.baidu.com/6ONYsjip0QIZ8tyhnq/it/u=1810950227,465524744&amp;fm=173&amp;s=AFF2439583026D5B124DA15E030050F3&amp;w=500&amp;h=281&amp;img.JPG"></div><p>赵庆伟的装裱店是村里第一家</p><p>“刚结婚那会儿五块钱也是借，十块钱也是借，真的是穷怕了。”王喜梅从没下过田地，为此她禁不住夸赞赵庆伟，这辈子做对了裱画这一件事。她最难忘怀的是从王建辉夫妻俩手中接过第一笔裱画费用200块钱，“原来裱画是可以挣钱的！”当时的激动心情是至今一次结算几万块也无法比拟的。</p><p>赵庆伟带着记者看已经被废弃的老房子，这曾经是他亲手搬砖搭建起来供全家七口人一起居住的。夏天漏雨，冬天严寒，站在房前，赵庆伟忆苦思甜难掩激动之情。“以前是没地方住，现在我爸一套，我一套，还买了一间门面房，现在是住不过来，这是真的。” 赵庆伟笑得合不拢嘴。</p><div><img width="500px" src="https://ss1.baidu.com/6ONXsjip0QIZ8tyhnq/it/u=4013260450,4096512528&amp;fm=173&amp;s=4EF862D91C03E357001D811D03001056&amp;w=500&amp;h=281&amp;img.JPG"></div><p>赵庆伟曾亲手盖起的老房子</p><p>因为学不来画虎，读过大学的赵庆伟想到了学做裱画这门生意，他是王公庄第一个吃螃蟹的人，十几年裱画生意让赵庆伟的家里发生了翻天覆地的变化。“我这个二层楼一盖，村子里一下子多了好几家裱画店。”</p><p>“昔日困境贫如洗，立志苦干更天地。创业艰辛多风雨，彩虹映时知足惜。”这首出自赵庆伟之手的小诗正是他自身生活的真实写照。</p></div>`, 4)
}

func Parser(title, content string, t int) {
	var rule, result string
	var flog bool
	var err2 error
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
	switch t {
	case 1:
		println("开始执行简书分发")
		//   utils.JsonPath("parser", "./jianshu.json")
		rule = utils.Read("./parser/jianshu.json")
		flog, result, err2 = handler.PushComm(title, content, rule)
		println("简书分发", flog, result, err2)
		break
	case 2:
		println("开始执行知乎分发")
		//utils.JsonPath("parser", "./zhihu.json")
		rule = utils.Read("./parser/zhihu.json")
		flog, result, err2 = handler.PushComm(title, content, rule)
		println("知乎分发", flog, result, err2)
		break
	case 3:
		println("开始执行新浪分发")
		//utils.JsonPath("parser", "./sina.json")
		rule = utils.Read("./parser/sina.json")
		flog, result, err2 = handler.PushComm(title, content, rule)
		println("微博分发", flog, result, err2)
		break
	case 4:
		global.Page.Navigate("http://www.baidu.com")
		global.Page.FindByID("kw").SendKeys("{ENTER}")
		break
	}
}

func TestPlatform(t *testing.T) {
	m, err := api.GetPlatformList()
	if err != nil {
		t.Error(err)
	}
	t.Log(len(m))
}

func TestPlatformOne(t *testing.T) {
	m, err := api.GetPlatformOne("新浪微博")
	if err != nil {
		t.Error(err)
	}
	t.Log(len(m))
}

func TestPlatformPost(t *testing.T) {
	body, err := uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    "http://127.0.0.1:8090/api/platform/one?type=新浪微博",
	})
	if err != nil {
		t.Log(err)
	}
	t.Log(string(body[:]))
}

func TestFindText(t *testing.T) {
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
	global.Navigate("https://www.weibo.com/")

	time.Sleep(5 * time.Second)
	els, elsErr := global.Page.Find("#pl_unlogin_home_hots > div:nth-child(1) > div.UG_contents > div:nth-child(1) > div > div.pic.W_piccut_v > a > img").Elements()
	if elsErr != nil {
		println(elsErr.Error())
	}
	println("数量：", len(els))
	val, _ := els[0].GetAttribute("src")
	println("值：", val)
}

func TestMsg(t *testing.T) {
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
	global.Navigate("https://www.baidu.com/")
}
