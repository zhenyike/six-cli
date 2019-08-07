package six_cloud

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"strconv"
)

func (user *SixUser) GetFilesByPath(path string) ([]*SixFile, error) {
	var (
		page = 2
		body = `{"path":"` + path + `","pageSize":50,"page": 1}`
		info = gjson.Parse(user.Client.PostJson("https://api.6pan.cn/v2/files/page", body))
	)
	if !info.Get("success").Bool() {
		return nil, errors.New(info.Get("message").Str)
	}
	if info.Get("result.parent").Type == gjson.Null {
		return nil, errors.New("path not exists")
	}
	res := parseFiles(info.Get("result.list").Array())
	for int64(page) <= info.Get("result.totalPage").Int() {
		body = `{"path":"` + path + `","pageSize":50,"page": ` + strconv.FormatInt(int64(page), 10) + `}`
		info = gjson.Parse(user.Client.PostJson("https://api.6pan.cn/v2/files/page", body))
		if !info.Get("success").Bool() {
			return res, nil
		}
		res = append(res, parseFiles(info.Get("result.list").Array())...)
		page++
	}
	return res, nil
}

func parseFiles(list []gjson.Result) []*SixFile {
	var res []*SixFile
	for _, r := range list {
		var file *SixFile
		err := json.Unmarshal([]byte(r.Raw), &file)
		if err == nil {
			res = append(res, file)
		}
	}
	return res
}