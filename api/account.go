package api

import (
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils"
)

func CheckClientLogin(token string) *global.Account {
	b, err := uhttp.PostForm(global.API_SSO_DOMAN+"/pass/checkToken?token="+token, nil)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	rp := utils.ToResultParam(b)
	if rp.Ret != 200 {
		return nil
	}
	return global.ToInterfaceAccount(rp.Data)
}
