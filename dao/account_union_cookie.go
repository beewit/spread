package dao

import (
	"time"
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
)

func GetUnionCookies(domain string, accId int64) (string, error) {
	m, err := GetUnionCookie(domain, accId)
	if err != nil {
		return "", err
	}
	return convert.ToString(m["cookies"]), nil
}

func SetUnionCookies(domain, cookies string, accId int64) (bool, error) {
	if global.Acc == nil {
		return false, nil
	}
	iw, _ := utils.NewIdWorker(1)
	id, _ := iw.NextId()
	m, err := GetUnionCookie(domain, accId)
	if err != nil {
		return false, err
	}
	var x int64
	var err2 error
	if m == nil {
		//修改原有Cookie
		sql := `INSERT INTO account_union_cookie(id,account_id,domain,cookies,ut_time,ct_time,status) values(?,?,?,?,?,?,1)`
		nt := time.Now().Format("2006-01-02 15:04:05")
		x, err2 = global.SLDB.Insert(sql, id, accId, domain, cookies, nt, nt)
	} else {
		sql := `UPDATE account_union_cookie SET  cookies=?,ut_time=? WHERE domain=? AND account_id=?`
		nt := time.Now().Format("2006-01-02 15:04:05")
		x, err2 = global.SLDB.Update(sql, cookies, nt, domain, accId)
	}
	if err2 != nil {
		return false, err
	}
	return x > 0, err
}

func GetUnionCookie(domain string, accId int64) (map[string]interface{}, error) {
	if global.Acc == nil {
		return nil, nil
	}
	sql := `SELECT cookies FROM account_union_cookie WHERE domain=? AND account_id=? ORDER BY ut_time DESC LIMIT 1`
	m, err := global.SLDB.Query(sql, domain, accId)
	if err != nil {
		return nil, err
	}
	if len(m) <= 0 {
		return nil, nil
	}
	return m[0], nil
}
