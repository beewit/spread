package global

import (
	"encoding/json"
	"github.com/beewit/beekit/utils/convert"
	"time"
	"fmt"
)

type JSONTime time.Time

//实现它的json序列化方法
func (this JSONTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

type Account struct {
	Id              int64 `json:"id"`
	Gender          string `json:"gender"`
	Mobile          string `json:"mobile"`
	Photo           string `json:"photo"`
	Nickname        string `json:"nickname"`
	MemberTypeName  string `json:"member_type_name"`
	MemberExpirTime string `json:"member_expir_time"`
	MemberTypeId    int64 `json:"member_type_id"`
}

func ToByteAccount(b []byte) *Account {
	var rp = new(Account)
	err := json.Unmarshal(b[:], &rp)
	if err != nil {
		Log.Error(err.Error())
		return nil
	}
	return rp
}

func ToMapAccount(m map[string]interface{}) *Account {
	b := convert.ToMapByte(m)
	if b == nil {
		return nil
	}
	return ToByteAccount(b)
}

func ToInterfaceAccount(m interface{}) *Account {
	b := convert.ToInterfaceByte(m)
	if b == nil {
		return nil
	}
	return ToByteAccount(b)
}
