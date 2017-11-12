package handler

import (
	"fmt"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread/global"
	"github.com/labstack/echo"
)

func GetTask(c echo.Context) error {
	runTask := fmt.Sprintf(`
	function loadScript(url, callback) {
		var script = document.createElement("script");
		script.type = "text/javascript";
		if (typeof (callback) != "undefined") {
			if (script.readyState) {
				script.onreadystatechange = function () {
					if (script.readyState == "loaded" || script.readyState == "complete") {
						script.onreadystatechange = null;
						callback();
					}
				};
			} else {
				script.onload = function () {
					callback();
				};
			}
		}
		script.src = url;
		document.body.appendChild(script);
	};
	function stopTask(key) {
    	loadScript("%s/task/stop.js?key=" + key + "&rand=" + Date.parse(new Date()))
	};
	function taskCallback(result) {
		var trTemplete = "<tr><td>{name}</td><td>{content}</td><td>{handle}</td></tr>";
		var trStr = '';
		for (var s in result) {
			if (result[s].state) {
				trStr += trTemplete.replace('{name}', result[s].name).replace('{content}', result[s].content).replace('{handle}', "<a onclick='stopTask(\"" + s + "\")'>停止</a>")
			}
		}
		if (trStr == "") {
			trStr = '<td colspan="3">暂无任务</td>'
		}
		document.getElementById("taskList").innerHTML = trStr
	};document.getElementById("version").innerHTML=%s`, global.Host, global.VersionStr)
	return utils.SuccessRespone(c, runTask+";taskCallback("+convert.ToObjStr(global.TaskList)+");")
}

func StopTask(c echo.Context) error {
	key := c.FormValue("key")
	if key != "" {
		global.DelTask(key)
	}
	return utils.SuccessRespone(c, "任务停止成功！")
}
