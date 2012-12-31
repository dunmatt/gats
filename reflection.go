package gats

import (
	"fmt"
	"reflect"
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
	if cont == nil {
		return false, fmt.Errorf("No field named %v in the data (or no data provided)", name)
	}
	switch reflect.ValueOf(cont.data).Kind().String() {
	case "ptr":
		field := reflect.ValueOf(cont.data).Elem().FieldByName(name)
		if field.IsValid() {
			return field.Bool(), nil
		}
	}
	return getBool(name, cont.enclosure)
}

func getLength(name string, cont *context) (int, error) {
	if cont == nil {
		return -1, fmt.Errorf("No field named %v in the data (or no data provided)", name)
	}
	switch reflect.ValueOf(cont.data).Kind().String() {
	case "ptr":
		field := reflect.ValueOf(cont.data).Elem().FieldByName(name)
		if field.IsValid() {
			return field.Len(), nil
		}
	}
	return getLength(name, cont.enclosure)
}

func getItem(name string, index int, cont *context) (*context, error) {
	if cont == nil {
		return nil, fmt.Errorf("No field named %v in the data (or no data provided)", name)
	}
	switch reflect.ValueOf(cont.data).Kind().String() {
	case "ptr":
		field := reflect.ValueOf(cont.data).Elem().FieldByName(name)
		if field.IsValid() {
			if 0 <= index && index <= field.Len() {
				return makeContext(field.Index(index).Interface(), cont), nil
			}
			return nil, fmt.Errorf("Index %v out of bounds.", index)
		}
	}
	return getItem(name, index, cont.enclosure)
}
