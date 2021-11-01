package models

import (
	"encoding/json"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/buger/jsonparser"
	"net/http"
	"net/url"
	"strings"
)

var ua2 = `okhttp/3.12.1;jdmall;android;version/10.1.2;build/89743;screen/1440x3007;os/11;network/wifi;`

type AutoGenerated struct {
	ClientVersion string `json:"clientVersion"`
	Client        string `json:"client"`
	Sv            string `json:"sv"`
	St            string `json:"st"`
	UUID          string `json:"uuid"`
	Sign          string `json:"sign"`
	FunctionID    string `json:"functionId"`
}

func getSign() *AutoGenerated {
	data, _ := httplib.Get("https://pan.smxy.xyz/sign").SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36").Bytes()
	t := &AutoGenerated{}
	json.Unmarshal(data, t)
	logs.Info(t.Sign)
	if t != nil {
		t.FunctionID = "genToken"
	}
	return t
}

func getKey(WSCK string) (string, error) {
	v := url.Values{}
	s := getSign()
	v.Add("functionId", s.FunctionID)
	v.Add("clientVersion", s.ClientVersion)
	v.Add("client", s.Client)
	v.Add("uuid", s.UUID)
	v.Add("st", s.St)
	v.Add("sign", s.Sign)
	v.Add("sv", s.Sv)
	req := httplib.Post(`https://api.m.jd.com/client.action?` + v.Encode())
	req.Header("cookie", WSCK)
	req.Header("User-Agent", ua2)
	req.Header("content-type", `application/x-www-form-urlencoded; charset=UTF-8`)
	req.Header("charset", `UTF-8`)
	req.Header("accept-encoding", `br,gzip,deflate`)
	req.Body(`body=%7B%22action%22%3A%22to%22%2C%22to%22%3A%22https%253A%252F%252Fplogin.m.jd.com%252Fcgi-bin%252Fm%252Fthirdapp_auth_page%253Ftoken%253DAAEAIEijIw6wxF2s3bNKF0bmGsI8xfw6hkQT6Ui2QVP7z1Xg%2526client_type%253Dandroid%2526appid%253D879%2526appup_type%253D1%22%7D&`)
	data, err := req.Bytes()
	if err != nil {
		return "", err
	}
	logs.Info(string(data))
	logs.Info("获取token正常")
	tokenKey, _ := jsonparser.GetString(data, "tokenKey")
	ptKey, err := appjmp(tokenKey)
	logs.Info(ptKey)
	if err != nil {
		return "", err
	}
	return ptKey, nil
}

func appjmp(tokenKey string) (string, error) {
	v := url.Values{}
	v.Add("tokenKey", tokenKey)
	v.Add("to", ``)
	v.Add("client_type", "android")
	v.Add("appid", "879")
	v.Add("appup_type", "1")
	req := httplib.Get(`https://un.m.jd.com/cgi-bin/app/appjmp?` + v.Encode())
	req.Header("User-Agent", ua2)
	req.Header("accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3`)
	req.SetCheckRedirect(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	})
	rsp, err := req.Response()
	if err != nil {
		return "", err
	}
	cookies := strings.Join(rsp.Header.Values("Set-Cookie"), " ")
	//ptKey := FetchJdCookieValue("pt_key", cookies)
	return cookies, nil
}
