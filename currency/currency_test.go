package currency

import (
	"log"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCurrency(t *testing.T) {
	client := NewCurrency("c49fe436cf12f81458f6304c1083b95a")
	catagories, err := client.Catagories()
	assert.Nil(t, err)
	log.Println(catagories)

	data := ExchangeData{
		FromCurrency: "JPY",
		ToCurrency: "CNY",
		Amount: 1000,
	}

	err = client.Exchange(&data)
	assert.Nil(t, err)
	log.Println(data)
}

