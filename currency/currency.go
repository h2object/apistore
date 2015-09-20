package currency

import (
	"fmt"
	"net/url"
	"encoding/json"
	"github.com/h2object/apistore"
)

//!------------------------------
type JSONResponse struct{
	ErrNum int 	`json:"errNum"`
	ErrMsg string	`json:"errMsg"`
	RetData json.RawMessage `json:"retData"`
}

//!------------------------------
type Currency struct{
	client *apistore.Client
}

func NewCurrency(apikey string) *Currency {
	prepare := apistore.NewBaiduPreparer(apikey)
	parser := apistore.NewBaiduParser(nil)
	c := apistore.NewClient(prepare, parser, nil)
	return &Currency{c}
}

func (c *Currency) Catagories() ([]string, error) {
	var js JSONResponse
	url := apistore.BuildHttpURL(apistore.BaiduApistore, "/apistore/currencyservice/type", nil)

	if err := c.client.Get(url, &js); err != nil {
		return nil, err
	}

	var types []string
	if err := json.Unmarshal(js.RetData, &types); err != nil {
		return nil, err
	}

	return types, nil
}

type ExchangeData struct{
	Date string `json:"date"`
	Time string `json:"time"`
	FromCurrency string `json:"fromCurrency"`
	ToCurrency string `json:"toCurrency"`
	Amount float64 `json:"amount"`
	Currency float64 `json:"currency"`
	ConvertAmount float64 `json:"convertedamount"`
}

func (c *Currency) Exchange(data *ExchangeData) error {	
	params := url.Values{}
	params.Set("fromCurrency", data.FromCurrency)
	params.Set("toCurrency", data.ToCurrency)
	params.Set("amount", fmt.Sprintf("%.4f", data.Amount))
	url := apistore.BuildHttpURL(apistore.BaiduApistore, "/apistore/currencyservice/currency", params)

	var js JSONResponse
	if err := c.client.Get(url, &js); err != nil {
		return err
	}

	if err := json.Unmarshal(js.RetData, data); err != nil {
		return err
	}

	return nil
}