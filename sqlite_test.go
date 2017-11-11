package main

import (
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/spread/dao"
	"testing"
)

func Test(t *testing.T) {
	dao.SetUnion(enum.QQ, "123456", "123456", "123456", enum.QQ_ID, 34)
}
