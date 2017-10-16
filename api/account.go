package api

import (
	"fmt"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/spread/global"
	"github.com/labstack/echo"
)

func CheckClientLogin(token string) *global.Account {
	b, err := uhttp.PostForm(global.API_SSO_DOMAN+"/pass/checkToken?token="+token, nil)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	global.Log.Info("SSOToken：%s", string(b))
	rp := utils.ToResultParam(b)
	if rp.Ret != 200 {
		return nil
	}
	acc := global.ToInterfaceAccount(rp.Data)
	acc.Token = token
	return acc
}

func DeleteToken(token string) bool {
	url := fmt.Sprintf("/pass/deleteToken?token=%s", token)
	b, err := uhttp.PostForm(global.API_SSO_DOMAN+url, nil)
	if err != nil {
		global.Log.Error(err.Error())
		return false
	}
	rp := utils.ToResultParam(b)
	if rp.Ret != 200 {
		return false
	}
	return true
}

func UpdatePwd(c echo.Context) error {
	pwd := c.FormValue("pwd")
	pwdNew := c.FormValue("pwdNew")
	url := fmt.Sprintf("/api/account/updatePwd?pwd=%s&pwdNew=%s", pwd, pwdNew)

	rp, err := ApiPost(global.API_SERVICE_DOMAN+url, nil)
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "修改密码失败")
	}
	if rp.Ret != 200 {
		return utils.ErrorNull(c, "修改密码失败")
	}
	return utils.SuccessNull(c, "修改密码成功")
}
