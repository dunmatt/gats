package gats

import (
	"exp/html"
	"fmt"
	"reflect"
	"unsafe"
)

type context struct {
	enclosure *context
	data      interface{}
}

func makeContext(data interface{}, parent *context) *context {
	return &context{
		enclosure: parent,
		data:      data,
	}
}

func getBool(name string, cont *context) (bool, error) {
	field := getField(name, cont)
	if field.IsValid() {
		return field.Bool(), nil
	}
	return false, fmt.Errorf("No field named %v in the data (or no data provided)", name)
}

func getString(name string, cont *context) (string, error) {
	field := getField(name, cont)
	if field.IsValid() {
		return field.String(), nil
	}
	return "", fmt.Errorf("No field named %v in the data (or no data provided)", name)
}

func isString(name string, cont *context) bool {
	return getField(name, cont).Type().String() == "string"
}

func getStringMap(name string, cont *context) (map[string]string, error) {
	field := getField(name, cont)
	if field.IsValid() {
		results := make(map[string]string)
		for _, k := range field.MapKeys() {
			results[k.String()] = field.MapIndex(k).String()
		}
		return results, nil
	}
	return nil, fmt.Errorf("No field named %v in the data (or no data provided)", name)
}

func getHtmlNode(name string, cont *context) (*html.Node, error) {
	field := getField(name, cont)
	if field.IsValid() {
		return (*html.Node)(unsafe.Pointer(field.Pointer())), nil
	}
	return nil, fmt.Errorf("No field named %v in the data (or no data provided)", name)
}

func getLength(name string, cont *context) (int, error) {
	field := getField(name, cont)
	if field.IsValid() {
		return field.Len(), nil
	}
	return -1, fmt.Errorf("No field named %v in the data (or no data provided)", name)
}

func getItem(name string, index int, cont *context) (*context, error) {
	field := getField(name, cont)
	if field.IsValid() {
		if 0 <= index && index <= field.Len() {
			return makeContext(field.Index(index).Interface(), cont), nil
		}
		return nil, fmt.Errorf("Index %v out of bounds.", index)
	}
	return nil, fmt.Errorf("No field named %v in the data (or no data provided)", name)
}

func getField(name string, cont *context) (result reflect.Value) {
	if cont == nil {
		return reflect.Zero(reflect.TypeOf(name)) // any value of zero is fine here, all are !IsValid
	}
	switch reflect.ValueOf(cont.data).Kind().String() {
	case "ptr":
		result = reflect.ValueOf(cont.data).Elem().FieldByName(name)
	case "struct":
		result = reflect.ValueOf(cont.data).FieldByName(name)
	}
	if result.IsValid() {
		return result
	}
	return getField(name, cont.enclosure)
}
