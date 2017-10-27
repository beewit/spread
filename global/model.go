package global

import (
	"encoding/json"
	"fmt"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"time"
)

const (
	TASK_PLATFORM_PUSH         = "TASK_PLATFORM_PUSH"         //平台自动化群发文章内容
	TASK_WECHAT_ADD_GROUP      = "TASK_WECHAT_ADD_GROUP"      //自动化添加微信群
	TASK_WECHAT_SEND_MESSAGE   = "TASK_WECHAT_SEND_MESSAGE"   //批量发送微信群或人的消息
	TASK_WECHAT_ADD_GROUP_USER = "TASK_WECHAT_ADD_GROUP_USER" //自动化发起添加微信群成员
)

var (
	TaskNameMap = map[string]string{
		"TASK_PLATFORM_PUSH":         "平台自动化营销内容群发",
		"TASK_WECHAT_ADD_GROUP":      "自动化添加微信群",
		"TASK_WECHAT_SEND_MESSAGE":   "批量发送微信群或人的消息",
		"TASK_WECHAT_ADD_GROUP_USER": "自动化发起添加微信群成员"}
)

type JSONTime time.Time

//实现它的json序列化方法
func (this JSONTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

type Task struct {
	Name     string `json:"name"`
	Content  string `json:"content"`
	LastTime string `json:"last_time"`
	State    bool   `json:"state"`
}

type Account struct {
	Id       int64  `json:"id"`
	Gender   string `json:"gender"`
	Mobile   string `json:"mobile"`
	Photo    string `json:"photo"`
	Nickname string `json:"nickname"`
	Token    string
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

func UpdateTask(key, content string) {
	task := TaskList[key]
	if task == nil {
		task = new(Task)
	}
	task.Name = TaskNameMap[key]
	task.Content = content
	task.State = true
	task.LastTime = utils.CurrentTime()
	TaskList[key] = task
}

func DelTask(key string) {
	task := TaskList[key]
	if task != nil {
		task.State = false
	}
}

func GetTask(key string) *Task {
	return TaskList[key]
}
