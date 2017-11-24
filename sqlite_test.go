package main

import (
	"testing"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/spread/global"
)

func Test(t *testing.T) {
	CheckClientLogin("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.NmfAoV1q5u1qc3xCGO8nC3qTo2Eu90sV6xrpcm9Tqto")
}

func CheckClientLogin(token string) *global.Account {
	b, err := uhttp.PostForm(global.API_SSO_DOMAIN+"/pass/checkToken?token="+token, nil)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	println(string(b))
	rp := utils.ToResultParam(b)
	if rp.Ret != 200 {
		return nil
	}
	acc := global.ToInterfaceAccount(rp.Data)
	acc.Token = token
	return acc
}
