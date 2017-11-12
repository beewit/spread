package dao

import (
	"database/sql"
	"fmt"
	"github.com/beewit/beekit/sqlite"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread-update/update"
	"github.com/beewit/spread/global"
	"github.com/pkg/errors"
)

func GetVersion() (int, error) {
	sql := `SELECT version FROM version WHERE id=1 LIMIT 1`
	m, err := global.SLDB.Query(sql)
	if err != nil {
		return 0, err
	}
	if len(m) <= 0 {
		return 0, errors.New("数据库版本错误")
	}
	return convert.MustInt(m[0]["version"]), nil
}

func UpdateVersion(version int) error {
	sql := `UPDATE version SET version=? WHERE id=1`
	_, err := global.SLDB.Update(sql, version)
	if err != nil {
		return err
	}
	return err
}

//检查数据库更新
func CheckDatabase(newVersion int) error {
	thisVersion, err := GetVersion()
	if err != nil {
		return err
	}
	if thisVersion != newVersion {
		//版本维护
		switch newVersion {
		case 1:
		case 2:
		case 3:

		}
	}
	return nil
}

func GetOldTableName(table string) (old string, err error) {
	date := utils.CurrentDate()
	sql := fmt.Sprintf(`SELECT name FROM sqlite_master WHERE name LIKE '_%s_old_%s%%' ORDER BY name DESC LIMIT 1`, table, date)
	m, err := global.SLDB.Query(sql)
	if err != nil {
		return old, err
	}
	if len(m) <= 0 {
		return fmt.Sprintf("_%s_old_%s", table, date), nil
	}
	name := convert.ToString(m[0]["name"])
	return name, nil
}

//查询所有表
func SelectTables() {
	// select * from sqlite_master WHERE type = "table";
}

//查询所有表所有字段
func SelectFail(table string) {
	//PRAGMA table_info(version)
}

//查询表所有数据并导入新表数据
func ImportData() {
	//1、查询旧版数据库所有数据库表，进行循环操作获取表字段
	//SelectTables()
	//2、查询旧版表字段
	//SelectFail()
	//3、插入数据
	//INSERT INTO table(name) VALUES(value)
}

func CheckUpdate() (err error) {
	var version int
	version, err = GetVersion()
	if err != nil {
		return
	}
	_, err = update.DBUpdate(update.Version{Major: version, Minor: 0, Patch: 0}, func(fileNames []string) {
		if len(fileNames) > 0 {
			for _, name := range fileNames {
				if name == "spread.db" {
					InitDB()
				}
			}
		}
	})
	return
}

var (
	OldDB *sqlite.SqlConnPool
	NewDB *sqlite.SqlConnPool
)

/**
特别注意，新版本必须包含兼容老版本数据库结构
*/
func InitDB() (err error) {
	var flog, OldFlog bool
	flog, err = utils.PathExists("spread.db")
	if err != nil {
		return
	}
	OldFlog, err = utils.PathExists("spread.db.old")
	if err != nil {
		return
	}
	if !flog || !OldFlog {
		err = errors.New("更新数据库文件已损坏无法更新")
		return
	}

	OldDB = &sqlite.SqlConnPool{
		DriverName:     "sqlite3",
		DataSourceName: "spread.db.old",
	}
	OldDB.SqlDB, err = sql.Open(OldDB.DriverName, OldDB.DataSourceName)
	if err != nil {
		return
	}
	NewDB = &sqlite.SqlConnPool{
		DriverName:     "sqlite3",
		DataSourceName: "spread.db",
	}
	NewDB.SqlDB, err = sql.Open(NewDB.DriverName, NewDB.DataSourceName)
	if err != nil {
		return
	}
	ImportData()
	return
}
