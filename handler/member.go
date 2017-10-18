package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
	"github.com/labstack/echo"
)

func GetMemberInfo(c echo.Context) error {
	global.Acc = CheckClientLogin()
	return utils.Success(c, "", map[string]interface{}{"account": global.Acc, "host": global.Host})
}

func CreateWechatQrCode(c echo.Context) error {
	qrCode, err := utils.CreateQrCode("http://m.tbqbz.com/account/wechatBind?token=" + global.Acc.Token)
	if err != nil {
		return utils.ErrorNull(c, "生成绑定微信二维码错误！")
	}
	return utils.SuccessNullMsg(c, qrCode)
}
