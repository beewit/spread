package main

import (
	"testing"
	"github.com/beewit/spread/dao"
	"github.com/beewit/spread/global"
)

func TestUnionCookies(t *testing.T) {
	acc := new(global.Account)
	acc.Id = 9
	global.Acc = acc
	f, e := dao.SetUnionCookies("sina", "12345647899999", 1)
	if e != nil {
		t.Error(e)
	}
	t.Log(f)
}

func TestGetUnionCookies(t *testing.T) {
	acc := new(global.Account)
	acc.Id = 9
	global.Acc = acc
	f, e := dao.GetUnionCookies("sina" , 1)
	if e != nil {
		t.Error(e)
	}
	t.Log(f)
}
