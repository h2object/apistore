package apistore

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"errors"
	"strings"
)

var UserAgent = "Golang h2object/apistore package"

// Baidu Request Preparer Interface Definition
type RequestPreparer interface {
	Prepare(*http.Request) *http.Request
}

// Baidu Response Parser Interface Definition
type ResposeParser interface{
	Parse(ret interface{}, resp *http.Response) error 
}

// Baidu Request URL Builder
func BuildHttpURL(addr string, uri string, vals url.Values) *url.URL {
	u := &url.URL{
		Scheme: "http",
		Host: addr,
		Path: uri,
	}
	if vals != nil {
		u.RawQuery = vals.Encode()	
	}
	return u
}

func BuildHttpsURL(addr string, uri string, vals url.Values) *url.URL {
	u := &url.URL{
		Scheme: "https",
		Host: addr,
		Path: uri,
	}
	if vals != nil {
		u.RawQuery = vals.Encode()	
	}
	return u	
}

type Client struct {
	conn     	*http.Client
	preparer   	RequestPreparer
	parser 		ResposeParser
}

func NewClient(prepare RequestPreparer, parse ResposeParser, clt *http.Client) *Client{
	if clt == nil {
		clt = http.DefaultClient
	}
	return &Client{
		conn: clt,
		preparer: prepare,
		parser: parse,
	}
}

func (c *Client) Prepare(prepare RequestPreparer) {
	c.preparer = prepare
}

func (c *Client) Parser(parse ResposeParser) {
	c.parser = parse
}



func (c *Client) sent(method string, u *url.URL, bodyType string, body io.Reader, bodyLength int) (resp *http.Response, err error) {
	var req *http.Request

	upperMethod := strings.ToUpper(method)
	switch upperMethod {
	case "GET":
		fallthrough
	case "POST":
		fallthrough
	case "PATCH":
		fallthrough
	case "PUT":
		fallthrough
	case "DELETE":
		req, err = http.NewRequest(upperMethod, u.String(), body)
		req.Header.Set("Content-Type", bodyType)
		if bodyLength > 0 {
			req.ContentLength = int64(bodyLength)	
		}		
		if err != nil {
			return
		}	
	default:
		err = errors.New("unsupport method: " + method)
		return
	}
	
	return c.do(req)
}

func (c *Client) do(req *http.Request) (resp *http.Response, err error) {
	var real *http.Request
	if c.preparer != nil {
		real = c.preparer.Prepare(req)
	} else {
		real = req
	}
	real.Header.Set("User-Agent", UserAgent)

	resp, err = c.conn.Do(real)
	if err != nil {
		return
	}
	return
}

func (c *Client) Get(u *url.URL, ret interface{}) error {
	resp, err := c.sent("GET", u, "", nil, 0)
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) GetResponse(u *url.URL) (*http.Response, error) {
	return c.sent("GET", u, "", nil, 0)
}

func (c *Client) Post(u *url.URL, bodyType string, body io.Reader, length int64, ret interface{}) error {
	resp, err := c.sent("POST", u, bodyType, body, int(length))
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) Put(u *url.URL, bodyType string, body io.Reader, length int64, ret interface{}) error {
	resp, err := c.sent("PUT", u, bodyType, body, int(length))
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) Patch(u *url.URL, bodyType string, body io.Reader, length int64, ret interface{}) error {
	resp, err := c.sent("PATCH", u, bodyType, body, int(length))
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) Delete(u *url.URL, ret interface{}) error {
	resp, err := c.sent("DELETE", u, "", nil, 0)
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) PostJson(u *url.URL, data interface{}, ret interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := c.sent("POST", u, "application/json", bytes.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) PutJson(u *url.URL, data interface{}, ret interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := c.sent("PUT", u, "application/json", bytes.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) PatchJson(u *url.URL, data interface{}, ret interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := c.sent("PATCH", u, "application/json", bytes.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) PostForm(u *url.URL, form map[string][]string, ret interface{}) error {
	msg := url.Values(form).Encode()
	resp, err := c.sent("POST", u, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) PutForm(u *url.URL, form map[string][]string, ret interface{}) error {
	msg := url.Values(form).Encode()
	resp, err := c.sent("PUT", u, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) PatchForm(u *url.URL, form map[string][]string, ret interface{}) error {
	msg := url.Values(form).Encode()
	resp, err := c.sent("PATCH", u, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) PostMultiPartForm(u *url.URL, multipart *MultipartForm, ret interface{}) error {
	ct, err := multipart.ContentType()
	if err != nil {
		return err
	}

	rd, err := multipart.Reader()
	if err != nil {
		return err
	}

	sz, err := multipart.Size()
	if err != nil {
		return err
	}

	resp, err := c.sent("POST", u, ct, rd, sz)
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) PutMultiPartForm(u *url.URL, multipart *MultipartForm, ret interface{}) error {
		ct, err := multipart.ContentType()
	if err != nil {
		return err
	}

	rd, err := multipart.Reader()
	if err != nil {
		return err
	}

	sz, err := multipart.Size()
	if err != nil {
		return err
	}

	resp, err := c.sent("PUT", u, ct, rd, sz)
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}

func (c *Client) PatchMultiPartForm(u *url.URL, multipart *MultipartForm, ret interface{}) error {
	ct, err := multipart.ContentType()
	if err != nil {
		return err
	}

	rd, err := multipart.Reader()
	if err != nil {
		return err
	}

	sz, err := multipart.Size()
	if err != nil {
		return err
	}

	resp, err := c.sent("PATCH", u, ct, rd, sz)
	if err != nil {
		return err
	}
	return c.parser.Parse(ret, resp)
}


