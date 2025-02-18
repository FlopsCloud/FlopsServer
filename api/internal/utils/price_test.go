package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

/*
go test -v -run ^TestCalculatePrice$ fca/api/internal/utils
*/

func TestCalculatePrice(t *testing.T) {

	var jsonData = `{
		"pricing": {
			"hourly": {
				"base_price": 10 ,
				"member_price": 8 ,
				"discounts": []
			},
			"daily": {
				"base_price": 10 ,
				"member_price": 8 ,
				"discounts": [
					{   
						"greater_or_equal": 10,
						"base_discount": -5,
						"member_discount": -5
					},
					{
						"greater_or_equal": 20,
						"base_discount": -20,
						"member_discount": -10
					}
				]
			},
			"monthly": {
				"base_price": 10 ,
				"member_price": 8 ,
				"discounts": [
					{   
						"greater_or_equal": 11,
						"base_discount": -5,
						"member_discount": -5
					},
					{
						"greater_or_equal": 12,
						"base_discount": -20,
						"member_discount": -10
					}
				]
			},
			"yearly": {
				"base_price": 1000 ,
				"member_price": 800 ,
				"discounts": [
					{   
						"greater_or_equal": 1,
						"base_discount": -5,
						"member_discount": -5
					},
					{
						"greater_or_equal": 3,
						"base_discount": -20,
						"member_discount": -10
					}
				]
			}
		}
	}`

	var data PricingData
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	validate, err := validateDiscountJson(&data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(validate)

	orgbasePrice, orgmemberPrice, basePrice, memberPrice, baseDiscount, memberDiscount := CalculatePrice(&data, "hourly", 1)
	fmt.Println(orgbasePrice, orgmemberPrice, basePrice, memberPrice, baseDiscount, memberDiscount)

	orgbasePrice, orgmemberPrice, basePrice, memberPrice, baseDiscount, memberDiscount = CalculatePrice(&data, "yearly", 1)
	fmt.Println(orgbasePrice, orgmemberPrice, basePrice, memberPrice, baseDiscount, memberDiscount)
	orgbasePrice, orgmemberPrice, basePrice, memberPrice, baseDiscount, memberDiscount = CalculatePrice(&data, "yearly", 2)
	fmt.Println(orgbasePrice, orgmemberPrice, basePrice, memberPrice, baseDiscount, memberDiscount)
	orgbasePrice, orgmemberPrice, basePrice, memberPrice, baseDiscount, memberDiscount = CalculatePrice(&data, "yearly", 3)
	fmt.Println(orgbasePrice, orgmemberPrice, basePrice, memberPrice, baseDiscount, memberDiscount)

	orgbasePrice, orgmemberPrice, basePrice, memberPrice, baseDiscount, memberDiscount = CalculatePrice(&data, "yearly", 21)
	fmt.Println(orgbasePrice, orgmemberPrice, basePrice, memberPrice, baseDiscount, memberDiscount)

}
