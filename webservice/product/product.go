package product

import (
	"encoding/json"
	"log"
)

type Product struct {
	ProductID      int    `json:"productid"`
	Manufacturer   string `json:"Manufacturer"`
	Sku            string `json:sku`
	Upc            string `json:upc`
	PricePerUnit   string `json:priceperunit`
	QuantityOnHand int    `json:quantityOnHand`
	ProductName    string `json:productName`
}

var Productlist []Product

func init() {
	productsJson := `[
		{
			"productid" :1 ,
			"Manufacturer": "M1",
			"sku":"Msku1",
			"upc":"upc1",
			"priceperunit": "11.001",
			"quantityOnHand":567,
			"productName": "Lamp shade"
		},
		{
			"productid" :2 ,
			"Manufacturer": "M2",
			"sku":"Msku2",
			"upc":"upc2",
			"priceperunit": "21.001",
			"quantityOnHand":267,
			"productName": " shade"
		},
		{
			"productid" : 3 ,
			"Manufacturer": "M3",
			"sku":"Msku3",
			"upc":"upc3",
			"priceperunit": "11.003",
			"quantityOnHand":36,
			"productName": "33 shade"
		}
	]`
	err := json.Unmarshal([]byte(productsJson), &Productlist)
	if err != nil {
		log.Fatal(err)
	}
}
