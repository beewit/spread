package api

import (
	"github.com/beewit/spread/dao"
	"github.com/beewit/beekit/utils/uhttp"
	"encoding/json"
	"github.com/beewit/spread/global"
	"github.com/beewit/beekit/utils"
)

func ApiPost(url string, m map[string]string) (utils.ResultParam, error) {
	token, _ := dao.QueryLoginToken()
	nm := map[string]string{"token": token}
	if m != nil {
		for k, v := range m {
			nm[k] = v
		}
	}
	b, _ := json.Marshal(nm)
	body, err := uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    global.API_SERVICE_DOMAN + url,
		Body:   b,
	})
	if err != nil {
		return utils.ResultParam{}, err
	}
	return utils.ToResultParam(body), nil
}
