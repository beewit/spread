package global

import (
	"github.com/beewit/beekit/utils/convert"
)

func CreateLoginTokenTable() (int64, error) {
	sql := `CREATE TABLE IF NOT EXISTS login_token (
			  token VARCHAR(255)
		   )`
	return SLDB.Insert(sql)
}

func QueryLoginToken() (string, error) {
	sql := `SELECT token FROM login_token LIMIT 1`
	m, err := SLDB.Query(sql)
	if err != nil {
		return "", err
	}
	if len(m) <= 0 {
		return "", nil
	}
	return convert.ToString(m[0]["token"]), nil
}

func InsertToken(token string) (bool, error) {
	sql := `DELETE FROM login_token;INSERT INTO login_token(token) values(?)`
	x, err := SLDB.Insert(sql, token)
	if err != nil {
		return false, err
	}
	return x > 0, err
}

func QueryTableExists(table string) (bool, error) {
	sql := "`SELECT count(*) as num FROM sqlite_master WHERE type='table' AND name=?;`"
	m, err := SLDB.Query(sql, table)
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
