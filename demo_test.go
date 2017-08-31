package main

import (
	"testing"

	"github.com/beewit/beekit/utils/convert"
	"encoding/json"
	"strconv"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils/uhttp"
	"time"
	"fmt"
)

type Article struct {
	WebSite string
	Title   string
	Created string
}
type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf(`"%s"`, time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}
func TestTimeJson(t *testing.T) {
	str := `{"Created":"2016-03-20T20:44:25.371Z","Title":"测试标题5","WebSite":"5-wow.com"}`
	var a Article
	err := json.Unmarshal([]byte(str), &a)
	if err != nil {
		println(err.Error())
	}
	println("78978")
}

func TestCreateTable(t *testing.T) {

	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.HezrOWC6gyT06oTOoBDMs0_NYLNA59Fk2UhI2bZ25cU"

	b, err := uhttp.PostForm("http://127.0.0.1:8081/pass/checkToken?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Ul4N-Z_SOAQp01NaY-5Me2qOgJdGluHGzO1c_dCeT2s", nil)
	//rp := utils.ToResultParam(b)
	//if rp.Ret == utils.SUCCESS_CODE {
	//	acc := global.ToInterfaceAccount(rp.Data)
	//	println(acc.Nickname)
	//} else {
	//	t.Error(rp.Msg)
	//}
	//println(string(b[:]))
	//if err != nil {
	//	global.Log.Error(err.Error())
	//	t.Error(b)
	//}

	//str := `{"data":{"gender":null,"id":122068319091036160,"member_expir_time":null,"member_type_id":0,"member_type_name":null,"mobile":"18223277005","nickname":null,"photo":null},"msg":"有效token","ret":200}`
	println(string(b[:]))
	var rp utils.ResultParam
	err = json.Unmarshal(b, &rp)
	if err != nil {
		println(err.Error())
	} else {
		acc := global.ToInterfaceAccount(rp.Data)
		println(acc.Nickname)
	}

	//rp := utils.ToResultParam(b)
	//if rp.Ret != 200 {
	//	t.Log(b)
	//}
	//t.Log(b)
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
