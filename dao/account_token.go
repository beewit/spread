package dao

import (
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
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
	return convert.ToString(m[0]["token"]), nil
}

func InsertToken(token string, acc *global.Account) (bool, error) {
	iw, _ := utils.NewIdWorker(1)
	id, _ := iw.NextId()
	sql := `DELETE FROM account_token WHERE account_id=?;INSERT INTO account_token(id,account_id,token,ut_time) values(?,?,?,?)`
	x, err := global.SLDB.Insert(sql, acc.Id, id, acc.Id, token, time.Now().Format("2006-01-02 15:04:05"))
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
