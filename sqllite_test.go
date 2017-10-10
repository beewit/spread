package main

import (
	"github.com/beewit/spread/dao"
	"github.com/beewit/spread/global"
	"testing"
)

func TestUnionCookies(t *testing.T) {
	acc := new(global.Account)
	acc.Id = 9
	global.Acc = acc

	f, e := dao.SetUnionCookies("sina", "12345647899999", "", "", 1, 1, "")
	if e != nil {
		t.Error(e)
	}
	t.Log(f)
}
