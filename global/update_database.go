package global

import (
	"database/sql"
	"fmt"

	"github.com/beewit/beekit/sqlite"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread-update/update"
	"github.com/pkg/errors"
)

func GetVersion() (int, error) {
	err := InitNewDB()
	if err != nil {
		return 0, err
	}
	defer NewDB.Close()
	sql := `SELECT version FROM local_version WHERE id=1 LIMIT 1`
	m, err := NewDB.Query(sql)
	if err != nil {
		return 0, err
	}
	if len(m) <= 0 {
		return 0, nil
	}
	return convert.MustInt(m[0]["version"]), nil
}

func UpdateVersion(version int) error {
	sql := `DELETE FROM local_version WHERE id=1;INSERT INTO local_version(id,version) VALUES(1,?)`
	_, err := NewDB.Update(sql, version)
	if err != nil {
		return err
	}
	return err
}

//查询所有表
func SelectTables() ([]map[string]interface{}, error) {
	sql := "select * from sqlite_master WHERE type = 'table'"
	return OldDB.Query(sql)
}

func QueryOldDBData(pageIndex int, table string) (*utils.PageData, error) {
	return OldDB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     table,
		PageIndex: pageIndex,
		PageSize:  100,
	})
}

//插入旧版数据
func InsertOldDBData(pageIndex int, table string) error {
	page, err := QueryOldDBData(pageIndex, table)
	if err != nil {
		return err
	}
	if page.Count <= 0 {
		return nil
	}
	Log.Info(" InsertOldDBData 正在转移本地数据库数据->新数据库中")
	//执行数据转移
	for _, v := range page.Data {
		_, err = NewDB.InsertMap(table, v)
		if err != nil {
			return err
		}
	}
	if page.PageIndex < page.PageSize {
		pageIndex++
		err = InsertOldDBData(pageIndex, table)
		if err != nil {
			return err
		}
	}
	return nil
}

//查询表所有数据并导入新表数据
func ImportData() error {
	Log.Info(" ImportData 准备本地数据库升级")
	//1、查询旧版数据库所有数据库表，进行循环操作获取表字段
	m, err := SelectTables()
	if err != nil {
		return err
	}
	if len(m) < 0 {
		return errors.New("无表结构")
	}
	for _, v := range m {
		tableName := convert.ToObjStr(v["name"])
		err = InsertOldDBData(1, tableName)
		if err != nil {
			println(err)
			continue
		}
	}
	return nil
}

func CheckUpdateDB() (err error) {
	var version int
	Log.Info("正在检查数据库版本")
	version, err = GetVersion()
	if err != nil {
		return
	}
	Log.Info("当前数据库版本：v%v", version)
	_, err = update.DBUpdate("app", update.Version{Major: 1, Minor: 0, Patch: version}, func(fileNames []string, rel update.Release) {
		if len(fileNames) > 0 {
			for _, name := range fileNames {
				if name == "spread.db" {
					Log.Info("初始化迁移数据")
					err2 := InitNewDB()
					if err2 != nil {
						Log.Info("InitNewDB ERROR", err.Error())
						return
					}
					defer NewDB.Close()
					err2 = InitOldDB()
					if err2 != nil {
						Log.Info("InitOldDB ERROR", err.Error())
						return
					}
					defer OldDB.Close()
					Log.Info("准备导入数据")
					err2 = ImportData()
					if err2 != nil {
						Log.Info("InitOldDB ERROR", err.Error())
						return
					}
					Log.Info("导入数据完成")
					UpdateVersion(rel.Patch)
					Log.Info("数据转移完成")
				}
			}
		}
	})
	if err != nil {
		Log.Info(err.Error())
	}
	return
}

var (
	OldDB *sqlite.SqlConnPool
	NewDB *sqlite.SqlConnPool
)

func InitNewDB() (err error) {
	var flog bool
	flog, err = utils.PathExists(SQLITE_DATABASE)
	if err != nil {
		return
	}
	if !flog {
		err = errors.New("更新数据库文件已损坏无法更新")
		return
	}
	NewDB = &sqlite.SqlConnPool{
		DriverName:     "sqlite3",
		DataSourceName: SQLITE_DATABASE,
	}
	NewDB.SqlDB, err = sql.Open(NewDB.DriverName, NewDB.DataSourceName)
	if err != nil {
		return
	}
	return
}

/**
特别注意，新版本必须包含兼容老版本数据库结构
*/
func InitOldDB() (err error) {
	var OldFlog bool
	oldDB := fmt.Sprintf("%s.old", SQLITE_DATABASE)
	OldFlog, err = utils.PathExists(oldDB)
	if err != nil {
		return
	}
	if !OldFlog {
		err = errors.New("更新数据库文件已损坏无法更新")
		return
	}

	OldDB = &sqlite.SqlConnPool{
		DriverName:     "sqlite3",
		DataSourceName: oldDB,
	}
	OldDB.SqlDB, err = sql.Open(OldDB.DriverName, OldDB.DataSourceName)
	if err != nil {
		return
	}
	return
}
