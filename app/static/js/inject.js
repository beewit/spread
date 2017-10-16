var time = 300;
var host = "";
var local = false;
if (location.href.indexOf("127.0.0.1") > -1 || location.href.indexOf('localhost') > -1) {
    local = true
} else {
    host = "http://127.0.0.1:8080";
}


if (document.getElementsByTagName("title").length > 0) {
    document.getElementsByTagName("title")[0].innerHTML = "工蜂小智";
} else {
    var title = document.createElement("title");
    title.innerHTML = "工蜂小智";
    document.head.appendChild(title);
}
if (document.querySelectorAll("[href='" + host + "/app/page/admin/assets/i/favicon.png']").length <= 0) {
    var link = document.createElement("link");
    link.setAttribute("type", "image/png");
    link.setAttribute("rel", "icon");
    link.setAttribute("href", host + '/app/page/admin/assets/i/favicon.png');
    document.head.appendChild(link);
}

$hiveIframe = document.getElementById("hive-iframe");
$hiveIframe.style.width = '0px';


//each 待定
var $dataSrc = document.querySelectorAll('[data-src]');
for (var i = 0; i < $dataSrc.length; i++) {
    $dataSrc[i].src = host + $dataSrc[i].getAttribute("data-src");
}

if (!local) {
    var $dataHref = document.querySelectorAll('[data-href]');
    for (var i = 0; i < $dataHref.length; i++) {
        $dataHref[i].addEventListener('click', function (e) {
            if (confirm("正在执行任务，是否确定终止任务！"))
                location.href = host + "?lastUrl=" + this.getAttribute("data-href");
        })
    }
}

//$appitemHook
var $appitemHook = document.getElementsByClassName('appitem-hook');
for (var i = 0; i < $appitemHook.length; i++) {

    $appitemHook[i].addEventListener('click', function (e) {

        for (var k = 0; k < $appitemHook.length; k++) {
            var _className = $appitemHook[k].className;
            $appitemHook[k].className = _className.replace(new RegExp('active', 'gi'), '');
        }

        this.className = this.className + ' active';

        var $hiveIframeDiv = document.getElementsByClassName('hive-iframe-div');
        for (var k = 0; k < $hiveIframeDiv.length; k++) {
            var _className = $hiveIframeDiv[k].className;

            $hiveIframeDiv[k].className = _className + ' active';
        }

        $hiveIframe.src = host + this.getAttribute('data-href')

        var $hive = document.getElementsByClassName('hive');

        if (new RegExp('hive-site', 'g').test($hive[0].className)) {
            var $cutOut = document.getElementsByClassName('cut-out');

            for (var k = 0; k < $cutOut.length; k++) {
                $cutOut[k].style.display = 'block';
            }
        }

    }, false);

}

var $cotOut = document.getElementsByClassName('cut-out');

for (var i = 0; i < $cotOut.length; i++) {

    $cotOut[i].addEventListener('click', function (e) {

        var $hiveIframeDiv = document.getElementsByClassName('hive-iframe-div');
        for (var k = 0; k < $hiveIframeDiv.length; k++) {
            var _className = $hiveIframeDiv[k].className;

            $hiveIframeDiv[k].className = _className.replace(new RegExp('active', 'gi'), '');
        }

        for (var k = 0; k < $cotOut.length; k++) {
            $cotOut[k].style.display = 'none';
        }
    }, false);
}

if (local) {
    var $hiveIframeDiv = document.getElementsByClassName('hive-iframe-div')[0];
    $hiveIframeDiv.className += ' active';
    document.getElementById("hive-iframe").src = host + $hiveIframe.getAttribute('data-src');

} else {
    var $hive = document.getElementsByClassName('hive')[0];
    $hive.className += " hive-site";
}


var winWidth = 0;
var winHeight = 0;

function findDimensions() //函数：获取尺寸
{
    //获取窗口宽度
    if (window.innerWidth)
        winWidth = window.innerWidth;
    else if ((document.body) && (document.body.clientWidth))
        winWidth = document.body.clientWidth;
    //获取窗口高度
    if (window.innerHeight)
        winHeight = window.innerHeight;
    else if ((document.body) && (document.body.clientHeight))
        winHeight = document.body.clientHeight;
    //通过深入Document内部对body进行检测，获取窗口大小
    if (document.documentElement && document.documentElement.clientHeight && document.documentElement.clientWidth) {
        winHeight = document.documentElement.clientHeight;
        winWidth = document.documentElement.clientWidth;
    }
    var $hiveIframe = document.getElementsByClassName('hive-iframe');
    for (var i = 0; i < $hiveIframe.length; i++) {
        $hiveIframe[i].style.width = winWidth + 'px';
        $hiveIframe[i].style.height = winHeight + 'px';
    }
}

findDimensions();
//调用函数，获取数值
window.onresize = findDimensions;

function GetRequest() {
    var url = location.search; //获取url中"?"符后的字串
    var theRequest = new Object();
    if (url.indexOf("?") != -1) {
        var str = url.substr(1);
        strs = str.split("&");
        for (var i = 0; i < strs.length; i++) {
            theRequest[strs[i].split("=")[0]] = unescape(strs[i].split("=")[1]);
        }
    }
    return theRequest;
}

var Request = GetRequest();

if (Request["lastUrl"]) {
    document.getElementById("hive-iframe").src = Request["lastUrl"];
}

