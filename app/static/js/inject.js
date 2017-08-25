var time = 0;
var host = "";
if (location.href.indexOf("127.0.0.1") <= -1) {
    host = "http://127.0.0.1:8080"
}
if (!window.jQuery) {
    var jqueryJs = "<script src=\"" + host + "/app/page/admin/assets/js/jquery.min.js\">" + "</scr" + "ipt>";
    if (host != "") {
        document.getElementsByTagName("body")[0].innerHTML += jqueryJs;
    } else {
        document.write(jqueryJs)
    }
    time = 300;
}
setTimeout(function () {
    $(function () {
        $(".hive-iframe").css("width", "")
        $.each($("[data-src]"), function (i, item) {
            $(item).attr("src", host + $(item).attr("data-src"));
        });
        $(".appitem-hook").click(function () {
            $(".appitem-hook").removeClass("active");
            $(this).addClass("active");
            $(".hive-iframe-div").addClass("active");
            $(".hive-iframe").attr("src", host + $(this).attr("data-href"));
            if ($(".hive").hasClass("hive-site")) {
                $(".cut-out").show()
            }
        });
        $(".cut-out").click(function () {
            $(".hive-iframe-div").removeClass("active");
            $(".cut-out").hide()
        });
        if (location.href.indexOf('127.0.0.1') > -1) {
            $(".hive-iframe-div").addClass("active");
            $(".hive-iframe").attr("src", host + $(".hive-iframe").attr("data-src"))
        } else {
            $(".hive").addClass("hive-site")
        }
    });
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
        //结果输出至两个文本框
        $(".hive-iframe").css("height", winHeight + "px").css("width", winWidth + "px");
    }

    findDimensions();
    //调用函数，获取数值
    window.onresize = findDimensions;
}, time)


