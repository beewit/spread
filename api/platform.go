package api

import (
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"net/url"
)

func GetPlatformList() ([]map[string]string, error) {
	body, err := uhttp.PostForm(global.API_SERVICE_DOMAN+"/api/platform", nil)
	if err != nil {
		return nil, nil
	}
	rp := utils.ToResultParam(body)
	if rp.Ret == utils.SUCCESS_CODE {
		m, err2 := convert.Obj2ListMapString(rp.Data)
		return m, err2
	} else {
		return nil, nil
	}
}

func GetPlatformOne(t string) (map[string]string, error) {
	v := url.Values{}
	v.Add("type", t)
	body, err := uhttp.PostForm(global.API_SERVICE_DOMAN+"/api/platform/one", v)
	if err != nil {
		return nil, nil
	}
	rp := utils.ToResultParam(body)
	if rp.Ret == utils.SUCCESS_CODE {
		m, err2 := convert.Obj2MapString(rp.Data)
		return m, err2
	} else {
		return nil, nil
	}
}
