package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

type Discount struct {
	GreaterOrEqual int `json:"greater_or_equal"`
	BaseDiscount   int `json:"base_discount"`
	MemberDiscount int `json:"member_discount"`
}

type PriceDetail struct {
	BasePrice   int        `json:"base_price"`
	MemberPrice int        `json:"member_price"`
	Discounts   []Discount `json:"discounts"`
}

type PricingData struct {
	Pricing map[string]PriceDetail `json:"pricing"`
}

func CalculatePrice(data *PricingData, pricingType string, count int) (int, int, int, int, int, int) {
	detail := data.Pricing[pricingType]
	// sort discounts by GreaterOrEqual
	sort.Slice(detail.Discounts, func(i, j int) bool {
		return detail.Discounts[i].GreaterOrEqual < detail.Discounts[j].GreaterOrEqual
	})

	fmt.Println(detail.Discounts)

	basePrice := detail.BasePrice * (count)
	memberPrice := detail.MemberPrice * (count)

	baseDiscount := 100
	memberDiscount := 100

	for _, discount := range detail.Discounts {
		if count >= discount.GreaterOrEqual {
			baseDiscount = 100 + discount.BaseDiscount
			memberDiscount = 100 + discount.MemberDiscount

		} else {
			break
		}
	}

	fmt.Println(baseDiscount, memberDiscount)

	basePrice = basePrice * baseDiscount / 100
	memberPrice = memberPrice * memberDiscount / 100

	return detail.BasePrice * (count), detail.MemberPrice * (count), basePrice, memberPrice, baseDiscount, memberDiscount
}

func validateDiscounts(discounts []Discount) bool {

	sort.Slice(discounts, func(i, j int) bool {
		return discounts[i].GreaterOrEqual < discounts[j].GreaterOrEqual
	})

	for i := 0; i < len(discounts)-1; i++ {
		// 检查后一个条件
		for j := i + 1; j < len(discounts); j++ {
			// 当大于等于值更大时，折扣应该更小或相等
			if discounts[j].GreaterOrEqual > discounts[i].GreaterOrEqual {
				if discounts[j].BaseDiscount > discounts[i].BaseDiscount ||
					discounts[j].MemberDiscount > discounts[i].MemberDiscount {
					return false
				}
			}
		}
	}
	return true
}

func validateDiscountJson(data *PricingData) (bool, error) {

	detail := data.Pricing["daily"]

	validate := validateDiscounts(detail.Discounts)
	if !validate {
		return false, errors.New("daily discount is invalid")
	}

	detail = data.Pricing["hourly"]

	validate = validateDiscounts(detail.Discounts)
	if !validate {
		return false, errors.New("hourly discount is invalid")
	}

	detail = data.Pricing["monthly"]

	validate = validateDiscounts(detail.Discounts)
	if !validate {
		return false, errors.New("monthly discount is invalid")
	}

	detail = data.Pricing["yearly"]
	validate = validateDiscounts(detail.Discounts)
	if !validate {
		return false, errors.New("yearly discount is invalid")
	}

	return true, nil
}

func test() {

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
                    "greater_or_equal": 2,
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

}

func Str2PricingData(jsonData string) PricingData {

	var data PricingData
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println(err)
		// panic(err)
	}

	return data

}
