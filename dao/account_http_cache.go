package dao

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
)

const (
	WECHAT = 0
	QQ     = 1
)

func InsertHTTPCache(iden, nickName, cache string, t int, acc *global.Account) (bool, error) {
	iw, _ := utils.NewIdWorker(1)
	id, _ := iw.NextId()
	m := map[string]interface{}{}
	m["id"] = id
	m["iden"] = iden
	m["nick_name"] = nickName
	m["type"] = t
	m["cache"] = cache
	m["ct_time"] = utils.CurrentTime()
	m["account_id"] = acc.Id
	_, err := global.SLDB.InsertMap("account_http_cache", m)
	return err == nil, err
}

func UpdateHTTPCache(iden, cache string, t int, acc *global.Account) (bool, error) {
	_, err := global.SLDB.Update("UPDATE account_http_cache SET cache=? WHERE iden=? AND type=? AND account_id=?",
		cache, iden, t, acc.Id)
	return err == nil, err
}

func QueryHTTPCacheList(t int, acc *global.Account) ([]map[string]interface{}, error) {
	m, err := global.SLDB.Query("SELECT * FROM account_http_cache WHERE type=? AND account_id=?", t, acc.Id)
	return m, err
}

func QueryHTTPCache(iden string, t int, acc *global.Account) (map[string]interface{}, error) {
	m, err := global.SLDB.Query("SELECT * FROM account_http_cache WHERE iden=? AND type=? AND account_id=? LIMIT 1", iden, t, acc.Id)
	if err != nil {
		return nil, err
	}
	if len(m) <= 0 {
		return nil, nil
	}
	return m[0], nil
}
