package main

import (
	"fmt"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/spread/global"
	"encoding/json"
)

func main() {
	body, err := uhttp.Cmd(uhttp.Request{
		Method: "GET",
		URL:    global.API_SERVICE_DOMAN + "/api/template",
	})
	if err != nil {
	}
	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err == nil {
		for k, v := range res {
			if k == "data" {
				data, err2 := json.Marshal(v)
				if err2 != nil {
					println(err2.Error())
				}
				res[k] = string(data[:])
			}
		}
		fmt.Sprintf("%+v", res)
	} else {
	}
}
