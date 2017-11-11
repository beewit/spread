package dao

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
	"time"
)

func SetUnion(platform, platformAcc, platformPwd, remark string, platformId, accId int64) (bool, error) {
	if global.Acc == nil {
		return false, nil
	}
	iw, _ := utils.NewIdWorker(1)
	id, _ := iw.NextId()
	m, err := GetUnion(platformId, accId, platformAcc)
	if err != nil {
		return false, err
	}
	nt := utils.CurrentTime()
	var x int64
	if m == nil {
		//修改原有Cookie
		sql := `INSERT INTO account_union(id,platform,status,ct_time,ut_time,platform_account,platform_password,account_id,platform_id,remark) values(?,?,1,?,?,?,?,?,?,?)`
		x, err = global.SLDB.Insert(sql, id, platform, nt, nt, platformAcc, platformPwd, accId, platformId, remark)
	} else {
		sql := `UPDATE account_union SET platform_password=?,ut_time=? WHERE platform_id=? AND account_id=? AND platform_account=? AND remark=?`
		x, err = global.SLDB.Update(sql, platformPwd, nt, platformId, accId, platformAcc, remark)
	}
	if err != nil {
		return false, err
	}
	return x > 0, err
}

func GetUnion(platformId, accId int64, platformAcc string) (map[string]interface{}, error) {
	if global.Acc == nil {
		return nil, nil
	}
	sql := `SELECT * FROM account_union WHERE platform_id=? AND account_id=? AND
		platform_account=? AND status=1 ORDER BY ut_time DESC LIMIT 1`
	m, err := global.SLDB.Query(sql, platformId, accId, platformAcc)
	if err != nil {
		return nil, err
	}
	if len(m) <= 0 {
		return nil, nil
	}
	return m[0], nil
}

func GetUnionList(platformId, accId int64) ([]map[string]interface{}, error) {
	if global.Acc == nil {
		return nil, nil
	}
	sql := `SELECT * FROM account_union WHERE account_id=? AND platform_id=? AND status=1 ORDER BY ut_time DESC`
	m, err := global.SLDB.Query(sql, accId, platformId)
	if err != nil {
		return nil, err
	}
	if len(m) <= 0 {
		return nil, nil
	}
	return m, nil
}

func GetUnionListByPlatformAcc(platformId, accId int64, platformAcc string) (map[string]interface{}, error) {
	if global.Acc == nil {
		return nil, nil
	}
	sql := `SELECT * FROM account_union WHERE account_id=? AND platform_id=? AND platform_account=? AND status=1 ORDER BY ut_time DESC
	LIMIT 1`
	m, err := global.SLDB.Query(sql, accId, platformId, platformAcc)
	if err != nil {
		return nil, err
	}
	if len(m) <= 0 {
		return nil, nil
	}
	return m[0], nil
}

func GetUnionListPage(accId int64, pageIndex, pageSize int) (*utils.PageData, error) {
	if global.Acc == nil {
		return nil, nil
	}
	page, err := global.SLDB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "account_union",
		Where:     "account_id=? AND status=1 ORDER BY ut_time DESC",
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, accId)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func UpdateUnionPhoto(nickname, photo, platformAcc string, platformId, accId int64) (bool, error) {
	if global.Acc == nil {
		return false, nil
	}
	nt := time.Now().Format("2006-01-02 15:04:05")
	sql := `UPDATE account_union SET  nickname=?,photo=?,ut_time=?,status=1 WHERE platform_id=? AND account_id=? AND platform_account=?`
	x, err := global.SLDB.Update(sql, nickname, photo, nt, platformId, accId, platformAcc)
	if err != nil {
		return false, err
	}
	return x > 0, err
}

func DeleteUnionById(id, accId int64) (bool, error) {
	sql := `DELETE FROM account_union_cookie WHERE platform_account IN
	(SELECT platform_account  FROM account_union  WHERE  account_id = ? AND id = ?) AND account_id = 5882608802350080;
	DELETE FROM account_union WHERE id=? AND account_id=?`
	x, err := global.SLDB.Delete(sql, accId, id, id, accId)
	if err != nil {
		return false, err
	}
	return x > 0, err
}
