package gats

import (
	//"fmt"
	"reflect"
)

func getBool(name string, data interface{}) (bool, bool) {
	val := reflect.ValueOf(data).Elem()
	valType := val.Type()
	for i := 0; i < val.NumField(); i++ {
		if name == valType.Field(i).Name {
			return val.Field(i).Bool(), true
		}
	}
	return false, false
}
