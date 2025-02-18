package utils

import (
	"fmt"
	"reflect"

	"github.com/zeromicro/go-zero/core/logx"
)

func CopyStruct(src, dst interface{}, isCopyEmptyString bool) error {
	sval := reflect.ValueOf(src).Elem()
	dval := reflect.ValueOf(dst).Elem()
	for i := 0; i < sval.NumField(); i++ {
		val := sval.Field(i)
		name := sval.Type().Field(i).Name
		kind := sval.Type().Field(i).Type.Kind()
		// fmt.Println(name, kind, val)
		if kind == reflect.Slice {
			continue
		}
		if kind == reflect.Struct {
			continue
		}

		if kind == reflect.String || kind == reflect.Int64 || kind == reflect.Int || kind == reflect.Uint64 {
			dvalue := dval.FieldByName(name)

			if dvalue.IsValid() {
				dkind := dvalue.Kind()
				if dkind == kind {
					// if !isCopyEmptyString {
					// 	if kind == reflect.String && val.String() == "" {
					// 	} else if kind == reflect.Int64 && val.Int() == 0 {
					// 	} else if kind == reflect.Uint64 && val.Uint() == 0 {
					// 	} else if kind == reflect.Int && val.Int() == 0 {
					// 	} else {
					// 		fmt.Println(name, kind, val)
					// 		dvalue.Set(val)

					// 	}

					// } else {
					fmt.Println(name, kind, val)
					dvalue.Set(val)
					// }

				} else {
					logx.Error("Err while copy ", dkind, " ", kind, " ", name)
				}
			} else {

			}
		}

	}
	return nil
}
