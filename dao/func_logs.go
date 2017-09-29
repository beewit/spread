package dao

import (
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils"
)

func SetFuncLogs(funcLogsMap map[string]interface{}) (bool, error) {
	var x int64
	var err error
	x, err = global.SLDB.InsertMap("func_logs", funcLogsMap)
	if err != nil {
		return false, err
	}
	return x > 0, err
}

func GetFuncLogsListPage(accId int64, pageIndex, pageSize int) (*utils.PageData, error) {
	if global.Acc == nil {
		return nil, nil
	}
	page, err := global.SLDB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "func_logs",
		Where:     "account_id=? ORDER BY ct_time DESC",
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, accId)
	if err != nil {
		return nil, err
	}
	return page, nil
}
