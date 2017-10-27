package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/spread/global"
	"github.com/labstack/echo"
)

func GetTask(c echo.Context) error {
	return utils.SuccessNullMsg(c, global.TaskList)
}
