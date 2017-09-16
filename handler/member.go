package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils"
)

func GetMemberInfo(c echo.Context) error {
	return utils.Success(c, "", map[string]interface{}{"account": global.Acc, "host": global.Host})
}
