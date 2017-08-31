package main

import (
	"testing"

	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils/convert"
	"encoding/json"
	"strconv"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/beekit/utils"
)

func TestCreateTable(t *testing.T) {
	b, err := uhttp.PostForm("http://127.0.0.1:8081/pass/checkToken?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.jIrvk6panpeFNw-irK4XWotxmH3-YNQ11FTqPubnWx0", nil)
	println(string(b[:]))
	if err != nil {
		global.Log.Error(err.Error())
		t.Error(b)
	}
	rp := utils.ToResultParam(b)
	if rp.Ret != 200 {
		t.Log(b)
	}
	t.Log(b)
	//x, err := global.CreateLoginTokenTable()
	//checkErr(t, err)
	//checkInt(t, x)
	//flog, err2 := global.InsertToken("1234568978951231231231321231")
	//checkErr(t, err)
	//t.Log(flog)
	//token, err2 := global.QueryLoginToken()
	//checkErr(t, err2)
	//t.Log(token)
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		return
	}
}

func checkInt(t *testing.T, x int64) {
	if x <= 0 {
		t.Error("返回数据条数为" + strconv.FormatInt(x, 10))
		return
	}
}

func checkMap(t *testing.T, m map[string]interface{}) {
	if m != nil {
		t.Error("数据为Null" + convert.ToString(m))
		return
	}
	jsons, err := json.Marshal(m)
	checkErr(t, err)
	t.Log(jsons)
}

func checkMaps(t *testing.T, m []map[string]interface{}) {
	if m != nil {
		t.Error("数据为Null" + convert.ToString(m))
		return
	}
	jsons, err := convert.ToArrayMapStr(m)
	checkErr(t, err)
	t.Log(jsons)
}
