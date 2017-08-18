package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
)

func Identity(c echo.Context) error {
	return utils.Success(c, "操作成功", "")
}
