package gats

import (
	//"fmt"
	"reflect"
)

func getBool(name string, data interface{}) bool {
	return reflect.ValueOf(data).Elem().FieldByName(name).Bool()
	//val := reflect.ValueOf(data).Elem()
	//valType := val.Type()
	//for i := 0; i < val.NumField(); i++ {
	//	if name == valType.Field(i).Name {
	//		return val.Field(i).Bool(), true
	//	}
	//}
	//return false
}

func getLength(name string, data interface{}) int {
	return reflect.ValueOf(data).Elem().FieldByName(name).Len()
	//val := reflect.ValueOf(data).Elem()
	////defer recover() // this is expected if length is not defined on the type of data
	//field := val.FieldByName(name)
	//return field.Len(), field != reflect.Zero(val.Type())
}
