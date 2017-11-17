package dao

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/encrypt"
	"github.com/beewit/spread/global"
	"strings"
	"time"
)

func QueryLoginToken() (string, error) {
	sql := `SELECT token FROM account_token ORDER BY ut_time DESC LIMIT 1`
	m, err := global.SLDB.Query(sql)
	if err != nil {
		return "", err
	}
	if len(m) <= 0 {
		return "", nil
	}
	return strings.Replace(convert.ToString(m[0]["token"]), encrypt.NewRsae().Md532(utils.GetMac()), "", 1), nil
}

func InsertToken(token string, acc *global.Account) (bool, error) {
	println("InsertToken 添加Token ")
	iw, _ := utils.NewIdWorker(1)
	id, _ := iw.NextId()
	sql := `DELETE FROM account_token WHERE account_id=?;INSERT INTO account_token(id,account_id,token,ut_time) values(?,?,?,?)`
	x, err := global.SLDB.Insert(sql, acc.Id, id, acc.Id, token+encrypt.NewRsae().Md532(utils.GetMac()), time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return false, err
	}
	return x > 0, err
}
func DeleteToken(acc *global.Account) (bool, error) {
	println("DeleteToken 删除Token ")
	sql := `DELETE FROM account_token WHERE account_id=?;`
	x, err := global.SLDB.Delete(sql, acc.Id)
	if err != nil {
		return false, err
	}
	return x > 0, err
}

func QueryWechatLogin(accId int64) (string, error) {
	sql := `SELECT login_wechat_user FROM account_token WHERE account_id=? ORDER BY ut_time DESC LIMIT 1`
	m, err := global.SLDB.Query(sql, accId)
	if err != nil {
		return "", err
	}
	if len(m) <= 0 {
		return "", nil
	}
	return convert.ToString(m[0]["login_wechat_user"]), nil
}

func InsertWechatLogin(wxLoginInfo string, acc *global.Account) (bool, error) {
	sql := `UPDATE account_token SET login_wechat_user=? WHERE account_id=?`
	x, err := global.SLDB.Insert(sql, wxLoginInfo, acc.Id)
	if err != nil {
		return false, err
	}
	return x > 0, err
}

func DeleteWechatLogin(acc *global.Account) (bool, error) {
	sql := `UPDATE account_token SET login_wechat_user=null WHERE account_id=?`
	x, err := global.SLDB.Delete(sql, acc.Id)
	if err != nil {
		return false, err
	}
	return x > 0, err
}

func QueryTableExists(table string) (bool, error) {
	sql := "`SELECT count(*) as num FROM sqlite_master WHERE type='table' AND name=?;`"
	m, err := global.SLDB.Query(sql, table)
	if err != nil {
		return false, err
	}
	if len(m) <= 0 {
		return false, nil
	}
	num, err := convert.ToInt64(m[0]["num"])
	if err != nil {
		return false, err
	}
	return num > 0, nil
}
