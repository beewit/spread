package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
	"github.com/labstack/echo"
)

func GetMemberInfo(c echo.Context) error {
	return utils.Success(c, "", map[string]interface{}{"account": global.Acc, "host": global.Host})
}
