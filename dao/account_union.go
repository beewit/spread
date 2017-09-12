package dao

import (
	"time"
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils"
)

func SetUnion(platform, acc, pwd string, accId int64) (bool, error) {
	if global.Acc == nil {
		return false, nil
	}
	iw, _ := utils.NewIdWorker(1)
	id, _ := iw.NextId()
	m, err := GetUnion(platform, accId)
	if err != nil {
		return false, err
	}
	var x int64
	var err2 error
	if m == nil {
		//修改原有Cookie
		sql := `INSERT INTO account_union(id,platform,status,ct_time,ut_time,platform_account,platform_password,account_id) values(?,?,1,?,?,?,?,?)`
		nt := time.Now().Format("2006-01-02 15:04:05")
		x, err2 = global.SLDB.Insert(sql, id, platform, nt, nt, acc, pwd, accId)
	} else {
		sql := `UPDATE account_union SET platform_account=?,platform_password=?,ut_time=? WHERE platform=? AND account_id=?`
		nt := time.Now().Format("2006-01-02 15:04:05")
		x, err2 = global.SLDB.Update(sql, acc, pwd, nt, platform, accId)
	}
	if err2 != nil {
		return false, err
	}
	return x > 0, err
}

func GetUnion(platform string, accId int64) (map[string]interface{}, error) {
	if global.Acc == nil {
		return nil, nil
	}
	sql := `SELECT * FROM account_union WHERE platform=? AND account_id=? ORDER BY ut_time DESC LIMIT 1`
	m, err := global.SLDB.Query(sql, platform, accId)
	if err != nil {
		return nil, err
	}
	if len(m) <= 0 {
		return nil, nil
	}
	return m[0], nil
}

func GetUnionList(accId int64) ([]map[string]interface{}, error) {
	if global.Acc == nil {
		return nil, nil
	}
	sql := `SELECT * FROM account_union WHERE account_id=? ORDER BY ut_time DESC`
	m, err := global.SLDB.Query(sql, accId)
	if err != nil {
		return nil, err
	}
	if len(m) <= 0 {
		return nil, nil
	}
	return m, nil
}

func UpdateUnionPhoto(nickname, photo, platform string, accId int64) (bool, error) {
	if global.Acc == nil {
		return false, nil
	}
	nt := time.Now().Format("2006-01-02 15:04:05")
	sql := `UPDATE account_union SET  nickname=?,photo=?,ut_time=? WHERE platform=? AND account_id=?`
	x, err := global.SLDB.Update(sql, nickname, photo, nt, platform, accId)
	if err != nil {
		return false, err
	}
	return x > 0, err
}
