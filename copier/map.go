package copier

import (
	"errors"
	"github.com/jinzhu/copier"
	"reflect"
)

func Copy(from, to interface{}) error {
	err := copier.Copy(to, from)
	if err != nil {
		return err
	}

	tt := reflect.TypeOf(to).Elem()
	vv := reflect.ValueOf(to).Elem()

	for i := 0; i < tt.NumField(); i++ {
		tField := tt.Field(i).Type
		vField := vv.Field(i)
		if tField.Kind() == reflect.Ptr {
			if vField.Elem().Interface() == reflect.Zero(tField.Elem()).Interface() {
				vField.Set(reflect.Zero(tField))
			}
		}
	}

	return nil
}

func MapTo(from, to interface{}) error {
	fromType := reflect.TypeOf(from)
	fromValue := reflect.ValueOf(from)
	toType := reflect.TypeOf(to)
	toValue := reflect.ValueOf(to)

	if fromType.Kind() != reflect.Ptr || toType.Kind() != reflect.Ptr {
		return errors.New("from and to must be ptr")
	}

	if fromType.Kind() == reflect.Ptr {
		fromType = fromType.Elem()
		fromValue = fromValue.Elem()
	}

	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
		toValue = toValue.Elem()
	}

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return errors.New("from and to's elem must be struct")
	}

	nums := fromType.NumField()

	for i := 0; i < nums; i++ {
		name := fromType.Field(i).Name
		_type := fromType.Field(i).Type

		toField, ok := toType.FieldByName(name)
		if ok {
			zero := reflect.Zero(_type)
			fieldValue := fromValue.Field(i)
			if fieldValue.Interface() != zero.Interface() {
				if toValue.FieldByName(name).Kind() == reflect.Ptr {

					if _type.Kind() == reflect.Ptr {
						toValue.FieldByName(name).Set(fromValue.Field(i))
					} else {
						if toField.Type.Elem() != _type {
							toValue.FieldByName(name).Set(fromValue.Field(i).Addr().Convert(toField.Type))
						} else {
							toValue.FieldByName(name).Set(fromValue.Field(i).Addr())
						}
					}
				} else {
					if _type.Kind() == reflect.Ptr {
						toValue.FieldByName(name).Set(fromValue.Field(i).Elem())
					} else {
						if toField.Type != _type {
							toValue.FieldByName(name).Set(fromValue.Field(i).Convert(toField.Type))
						} else {
							toValue.FieldByName(name).Set(fromValue.Field(i))
						}
					}
				}
			}
		}
	}

	return nil
}
