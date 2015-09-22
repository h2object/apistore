# apistore
baidu api store for golang, 百度APIStore接口Golang实现

## 快速指南(Quick Start)


````golang

	import "github.com/h2object/apistore/currency"

	client := currency.NewClient()
	// 获取汇率币种类型
	types, err := client.Catagories()

	// 获取汇率转换
	data := ExchangeData{
		FromCurrency: "JPY",
		ToCurrency: "CNY",
		Amount: 1000,
	}

	err = client.Exchange(&data)
	log.Println(data)

````	

## 已完成接入(持续接入中)

-	[汇率](http://apistore.baidu.com/apiworks/servicedetail/119.html)
-	[推送](http://push.baidu.com) 接入中...

