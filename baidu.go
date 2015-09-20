package apistore

import (
	"fmt"
	"net/http"
	"encoding/json"
)

const BaiduApistore = "apis.baidu.com"

type BaiduPreparer struct{
	apikey string
}

func NewBaiduPreparer(apikey string) *BaiduPreparer{
	return &BaiduPreparer{
		apikey: apikey,
	}
}

func (preparer *BaiduPreparer) Prepare(req *http.Request) *http.Request {
	req.Header.Del("apikey")
	req.Header.Add("apikey", preparer.apikey)
	return req
}

// interface for baidu apistore's service error
type BaiduErrorParser interface{
	parse(*http.Response) error
}

type BaiduParser struct{
	error_parser BaiduErrorParser
}

func NewBaiduParser(ep BaiduErrorParser) *BaiduParser{
	return &BaiduParser{
		error_parser: ep,
	}
}

func (parser *BaiduParser) Parse(ret interface{}, resp *http.Response) error {
	defer resp.Body.Close()

	if resp.StatusCode/100 == 2 {
		if ret != nil && resp.ContentLength != 0 {
			err := json.NewDecoder(resp.Body).Decode(ret)
			if err != nil {
				return err
			}
		}
		if resp.StatusCode == 200 {
			return nil
		}
	}

	if parser.error_parser == nil {
		return fmt.Errorf("code:%d, reason:%s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	return parser.error_parser.parse(resp)
}

